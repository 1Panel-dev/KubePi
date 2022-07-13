package cluster

import (
	"fmt"
	"strings"

	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (h *Handler) ListApiGroups() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		scope := ctx.Params().GetString("scope")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		_, rss, err := client.ServerGroupsAndResources()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		var (
			groups   []ApiGroupOption
			groupMap = make(map[string]struct{})
		)
		itemResource := ApiResourceOption{
			Resource: "*",
			Verbs:    []string{"create", "delete", "deletecollection", "get", "list", "patch", "update", "watch"},
		}
		if scope == "namespace" {
			groups = append(groups, ApiGroupOption{Group: "*", Resources: []ApiResourceOption{itemResource}})
		} else {
			groups = append(groups, ApiGroupOption{Group: "*", Resources: []ApiResourceOption{itemResource}})
		}
		for _, group := range rss {
			if strings.Contains(group.GroupVersion, "/") {
				group.GroupVersion = group.GroupVersion[0:strings.Index(group.GroupVersion, "/")]
			}
			name := group.GroupVersion
			if name == "v1" {
				name = ""
			}
			if _, ok := groupMap[name]; ok {
				continue
			}
			groupMap[name] = struct{}{}
			itemGroup := ApiGroupOption{Group: name, Resources: []ApiResourceOption{itemResource}}
			for _, resource := range group.APIResources {
				if !(scope == "namespace" && resource.Namespaced) && !(scope == "cluster" && !resource.Namespaced) {
					continue
				}
				itemGroup.Resources = append(itemGroup.Resources, ApiResourceOption{Resource: resource.Name, Verbs: resource.Verbs})
			}
			if len(itemGroup.Resources) > 1 {
				groups = append(groups, itemGroup)
			}
		}
		ctx.Values().Set("data", groups)
	}
}

type ApiGroupOption struct {
	Group     string              `json:"group"`
	Resources []ApiResourceOption `json:"resources"`
}
type ApiResourceOption struct {
	Resource string   `json:"resource"`
	Verbs    []string `json:"verbs"`
}

func (h *Handler) ListApiGroupResources() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		groupName := ctx.Params().GetString("group")
		scope := ctx.URLParam("scope")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		gss, rss, err := client.ServerGroupsAndResources()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		var groupVersions []v1.GroupVersionForDiscovery
		resourceSet := map[string]struct{}{}

		if groupName == "core" {
			groupName = ""
		}
		for i := range gss {
			if gss[i].Name == groupName {
				groupVersions = append(groupVersions, gss[i].Versions...)
			}
		}
		for i := range rss {
			for j := range groupVersions {
				if groupVersions[j].GroupVersion == rss[i].GroupVersion {
					for k := range rss[i].APIResources {
						if !strings.Contains(rss[i].APIResources[k].Name, "/") {
							switch scope {
							case "cluster":
								if !rss[i].APIResources[k].Namespaced {
									resourceSet[rss[i].APIResources[k].Name] = struct{}{}
								}
							case "namespace":
								if rss[i].APIResources[k].Namespaced {
									resourceSet[rss[i].APIResources[k].Name] = struct{}{}
								}
							default:
								resourceSet[rss[i].APIResources[k].Name] = struct{}{}
							}
						}
					}
				}
			}
		}
		var resources []string
		for k := range resourceSet {
			resources = append(resources, k)
		}
		ctx.Values().Set("data", resources)

	}
}
