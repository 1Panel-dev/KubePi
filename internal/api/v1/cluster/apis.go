package cluster

import (
	"fmt"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (h *Handler) ListApiGroups() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		groups, err := client.ServerGroups()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", groups.Groups)
	}
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
		k := kubernetes.NewKubernetes(*c)
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
