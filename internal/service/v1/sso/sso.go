package sso

import (
	"context"
	"errors"
	"fmt"
	v1Session "github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	v1Sso "github.com/KubeOperator/kubepi/internal/model/v1/sso"
	v1User "github.com/KubeOperator/kubepi/internal/model/v1/user"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/rolebinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/user"
	ssoClient "github.com/KubeOperator/kubepi/pkg/util/sso"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/coreos/go-oidc"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"time"
)

type Service interface {
	common.DBService
	TestConnect(sso *v1Sso.Sso) error
	List(options common.DBOptions) ([]v1Sso.Sso, error)
	Create(sso *v1Sso.Sso, options common.DBOptions) error
	Update(id string, sso *v1Sso.Sso, options common.DBOptions) error
	Status(options common.DBOptions) bool
	OpenID(openid *v1Sso.OpenID, options common.DBOptions) (v1Session.UserProfile, error)
	OpenIDConfig(clientId, clientSecret, issuerURL, redirectURL string) (*v1Sso.OpenID, error)
}

func NewService() Service {
	return &service{
		userService:        user.NewService(),
		roleBindingService: rolebinding.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	userService        user.Service
	roleBindingService rolebinding.Service
}

func (s *service) TestConnect(sso *v1Sso.Sso) error {
	// 测试连接，理论上不应强制用户开启SSO认证
	//if !sso.Enable {
	//	return errors.New("请先启用SSO")
	//}

	sc := ssoClient.NewSsoClient(sso.Protocol, sso.InterfaceAddress, sso.ClientId, sso.ClientSecret, sso.Enable)
	if err := sc.TestConnect(sso.InterfaceAddress); err != nil {
		return err
	}

	return nil
}

func (s *service) Create(sso *v1Sso.Sso, options common.DBOptions) error {
	sc := ssoClient.NewSsoClient(sso.Protocol, sso.InterfaceAddress, sso.ClientId, sso.ClientSecret, sso.Enable)
	// 当用户进行SSO配置时，应该为用户检测目标是否可连接
	if err := sc.TestConnect(sso.InterfaceAddress); err != nil {
		return err
	}

	db := s.GetDB(options)
	sso.UUID = uuid.New().String()
	sso.CreateAt = time.Now()
	sso.UpdateAt = time.Now()
	return db.Save(sso)
}

func (s *service) Update(id string, sso *v1Sso.Sso, options common.DBOptions) error {
	sc := ssoClient.NewSsoClient(sso.Protocol, sso.InterfaceAddress, sso.ClientId, sso.ClientSecret, sso.Enable)
	// 当用户进行SSO配置时，应该为用户检测目标是否可连接
	if err := sc.TestConnect(sso.InterfaceAddress); err != nil {
		return err
	}

	old, err := s.GetById(id, options)
	if err != nil {
		return err
	}
	sso.UUID = old.UUID
	sso.CreateAt = old.CreateAt
	sso.UpdateAt = time.Now()
	db := s.GetDB(options)
	if sso.Enable != old.Enable {
		err = db.UpdateField(sso, "Enable", sso.Enable)
		if err != nil {
			return err
		}
	}
	return db.Update(sso)
}

func (s *service) List(options common.DBOptions) ([]v1Sso.Sso, error) {
	db := s.GetDB(options)
	sso := make([]v1Sso.Sso, 0)
	if err := db.All(&sso); err != nil {
		return nil, err
	}
	return sso, nil
}

func (s *service) Status(options common.DBOptions) bool {
	db := s.GetDB(options)
	sso := make([]v1Sso.Sso, 0)
	if err := db.All(&sso); err != nil {
		return false
	}

	if len(sso) == 0 {
		return false
	}

	return sso[0].Enable
}

func (s *service) GetById(id string, options common.DBOptions) (*v1Sso.Sso, error) {
	db := s.GetDB(options)
	var sso v1Sso.Sso
	query := db.Select(q.Eq("UUID", id))
	if err := query.First(&sso); err != nil {
		return nil, err
	}
	return &sso, nil
}

func (s *service) OpenID(openid *v1Sso.OpenID, options common.DBOptions) (v1Session.UserProfile, error) {
	token, err := openid.Oauth2Config.Exchange(openid.Ctx, openid.Code)
	if err != nil {
		return v1Session.UserProfile{}, errors.New("交换Token失败: " + err.Error())
	}

	userInfo, err := openid.OidcProvider.UserInfo(openid.Ctx, oauth2.StaticTokenSource(token))
	if err != nil {
		return v1Session.UserProfile{}, errors.New("获取用户信息失败: " + err.Error())
	}
	// 获取用户名
	var claims struct {
		PreferredUsername string `json:"preferred_username"`
	}
	if err = userInfo.Claims(&claims); err != nil {
		return v1Session.UserProfile{}, err
	}

	// 初始化用户
	_, err = s.userService.GetByNameOrEmail(userInfo.Email, options)
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			// 创建本地账号，密码默认设置为`@=7kvi-$l*Pj+,s`，默认不开启MFA
			userProfile := &v1User.User{
				BaseModel: v1.BaseModel{
					ApiVersion: "v1",
					Kind:       "User",
				},
				Metadata: v1.Metadata{
					Name: claims.PreferredUsername,
				},
				NickName: claims.PreferredUsername,
				Email:    userInfo.Email,
				Language: openid.Language,
				IsAdmin:  false,
				Authenticate: v1User.Authenticate{
					Password: `@=7kvi-$l*Pj+,s`,
				},
				Type: v1User.LOCAL,
				Mfa: v1User.Mfa{
					Enable: false,
				},
			}
			tx, err := server.DB().Begin(true)
			if err != nil {
				return v1Session.UserProfile{}, err
			}
			if err = s.userService.Create(userProfile, common.DBOptions{DB: tx}); err != nil {
				_ = tx.Rollback()
				return v1Session.UserProfile{}, err
			}

			// 用户角色默认为ReadOnly
			binding := v1Role.Binding{
				BaseModel: v1.BaseModel{
					Kind:       "RoleBind",
					ApiVersion: "v1",
					CreatedBy:  "admin",
				},
				Metadata: v1.Metadata{
					Name: fmt.Sprintf("role-binding-%s-%s", "ReadOnly", claims.PreferredUsername),
				},
				Subject: v1Role.Subject{
					Kind: "User",
					Name: claims.PreferredUsername,
				},
				RoleRef: "ReadOnly",
			}
			if err = s.roleBindingService.CreateRoleBinding(&binding, common.DBOptions{DB: tx}); err != nil {
				_ = tx.Rollback()
				return v1Session.UserProfile{}, err
			}
			_ = tx.Commit()
			fmt.Println("SSO用户" + claims.PreferredUsername + "不存在，已自动创建本地账号")
		} else {
			return v1Session.UserProfile{}, errors.New(fmt.Sprintf("query user %s failed ,: %s", claims.PreferredUsername, err.Error()))
		}
	}

	// 设置profile
	return s.localProfile(claims.PreferredUsername, userInfo.Email)
}

func (s *service) OpenIDConfig(clientId, clientSecret, issuerURL, redirectURL string) (*v1Sso.OpenID, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		return nil, err
	}

	return &v1Sso.OpenID{
		Oauth2Config: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  redirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		OidcProvider: provider,
		Ctx:          ctx,
	}, nil
}

func (s *service) localProfile(username, email string) (v1Session.UserProfile, error) {
	u, err := s.userService.GetByNameOrEmail(email, common.DBOptions{})
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return v1Session.UserProfile{}, errors.New("username or password error")
		}
		return v1Session.UserProfile{}, errors.New(fmt.Sprintf("query user %s failed ,: %s", username, err.Error()))
	}

	handler := v1Session.NewHandler()
	permissions, err := handler.AggregateResourcePermissions(username)
	if err != nil {
		return v1Session.UserProfile{}, errors.New(err.Error())
	}
	return v1Session.UserProfile{
		Name:                u.Name,
		NickName:            u.NickName,
		Email:               u.Email,
		Language:            u.Language,
		ResourcePermissions: permissions,
		IsAdministrator:     u.IsAdmin,
		Mfa: v1Session.Mfa{
			Secret:   u.Mfa.Secret,
			Enable:   u.Mfa.Enable,
			Approved: false,
		},
	}, nil
}
