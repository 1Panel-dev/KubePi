package group

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/certificate"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func (h *Handler) GetGroupMembers() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		bindings, err := h.groupBindingService.ListByGroupName(groupName, common.DBOptions{})
		if err != nil {
			if !errors.As(err, &storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		var members []Member

		for i := range bindings {
			members = append(members, Member{
				Name:     bindings[i].UserRef,
				Group:    bindings[i].GroupRef,
				CreateAt: bindings[i].CreateAt,
				CreateBy: bindings[i].CreatedBy,
			})
		}

		ctx.Values().Set("data", members)
	}
}

func (h *Handler) CreateGroupMember() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		var req Member
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		binding := v1Group.Binding{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "GroupBinding",
				CreatedBy:  profile.Name,
			},
			Metadata: v1.Metadata{
				Name: fmt.Sprintf("%s-%s-binding", req.Name, groupName),
			},
			UserRef:  req.Name,
			GroupRef: groupName,
		}
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}

		if err := h.groupBindingService.Create(&binding, txOptions); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		cbs, err := h.clusterBindingService.GetBindingsBySubject(v1Cluster.Subject{Kind: "User", Name: req.Name}, txOptions)
		if err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		for i := range cbs {
			cert, err := certificate.ParseX509Certificate(cbs[i].Certificate)
			if err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			var currentGroupNames []string
			currentGroupNames = append(cert.Subject.Organization, groupName)

			c, err := h.clusterService.Get(cbs[i].ClusterRef, txOptions)
			if err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			client := kubernetes.NewKubernetes(*c)
			userCert, err := client.CreateCommonUser(req.Name, currentGroupNames...)
			if err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			cbs[i].Certificate = userCert
			if err := h.clusterBindingService.UpdateClusterBinding(cbs[i].Name, &cbs[i], txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		_ = tx.Commit()
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteGroupMember() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		memberName := ctx.Params().GetString("member")
		gb, err := h.groupBindingService.GetGroupBindingByGroupNameAndUserName(memberName, groupName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}

		//重新签发证书

		cbs, err := h.clusterBindingService.GetBindingsBySubject(v1Cluster.Subject{Kind: "User", Name: memberName}, txOptions)
		if err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		for i := range cbs {
			cert, err := certificate.ParseX509Certificate(cbs[i].Certificate)
			if err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			var currentGroupNames []string
			exists := false
			for j := range cert.Subject.Organization {
				if cert.Subject.Organization[j] != groupName {
					currentGroupNames = append(currentGroupNames, cert.Subject.Organization[j])
				} else {
					exists = true
				}
			}
			if exists {
				c, err := h.clusterService.Get(cbs[i].ClusterRef, txOptions)
				if err != nil {
					_ = tx.Rollback()
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
				client := kubernetes.NewKubernetes(*c)
				userCert, err := client.CreateCommonUser(memberName, currentGroupNames...)
				if err != nil {
					_ = tx.Rollback()
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
				cbs[i].Certificate = userCert
				if err := h.clusterBindingService.UpdateClusterBinding(cbs[i].Name, &cbs[i], txOptions); err != nil {
					_ = tx.Rollback()
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
			}
		}
		if err := h.groupBindingService.Delete(gb.Name, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
	}
}
