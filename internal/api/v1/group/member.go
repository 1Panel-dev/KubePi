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
	"reflect"
	"sort"
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

		if err := h.groupBindingService.Create(&binding, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		go h.modifyClusterUserCert(req.Name)
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
		if err := h.groupBindingService.Delete(gb.Name, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		go h.modifyClusterUserCert(memberName)
	}
}

func (h *Handler) modifyClusterUserCert(userName string) {
	log := server.Logger()

	log.Infof("modify user %s certificate", userName)
	cbs, err := h.clusterBindingService.GetBindingsBySubject(v1Cluster.Subject{Kind: "User", Name: userName}, common.DBOptions{})
	if err != nil && !errors.Is(err, storm.ErrNotFound) {
		log.Error(err)
		return
	}
	if !(len(cbs) > 0) {
		log.Infof("user %s not in any cluster,skip it", userName)
		return
	}
	// 查询现在用户在哪些用户组中

	gs, err := h.groupBindingService.ListByUserName(userName, common.DBOptions{})
	if err != nil && !errors.Is(err, storm.ErrNotFound) {
		log.Error(err)
		return
	}

	var orgNames []string
	for i := range gs {
		orgNames = append(orgNames, gs[i].GroupRef)
	}
	sort.Strings(orgNames)

	for i := range cbs {
		log.Infof("redeploy user %s certificate at cluster %s started", userName, cbs[i].ClusterRef)
		cert, err := certificate.ParseX509Certificate(cbs[i].Certificate)
		if err != nil {
			server.Logger().Error(err)
			return
		}
		currentGroups := cert.Subject.Organization
		sort.Strings(currentGroups)
		if !reflect.DeepEqual(orgNames, currentGroups) {
			c, err := h.clusterService.Get(cbs[i].ClusterRef, common.DBOptions{})
			if err != nil {
				log.Error(err)
				return
			}
			client := kubernetes.NewKubernetes(*c)
			userCert, err := client.CreateCommonUser(userName, orgNames...)
			if err != nil {
				log.Error(err)
				return
			}
			cbs[i].Certificate = userCert
			if err := h.clusterBindingService.UpdateClusterBinding(cbs[i].Name, &cbs[i], common.DBOptions{}); err != nil {
				log.Error(err)
				return
			}
		}
		log.Infof(" redeploy user %s certificate at cluster %s finished", userName, cbs[i].ClusterRef)
	}

}
