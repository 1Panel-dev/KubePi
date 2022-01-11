import Layout from "@/business/app-layout/horizontal-layout"

const Clusters = {
    path: "/cluster",
    parent: true,
    global: true,
    sort: 1,
    component: Layout,
    name: "Cluster",
    meta: {
        title: "business.cluster.cluster",
        icon: "iconfont iconkubernets",
    },
    children: [
        {
            path: "/nodes",
            requirePermission: {
                apiGroup: "",
                resource: "nodes",
                verb: "list",
                scope: "cluster"
            },
            component: () => import("@/business/cluster/nodes"),
            name: "Nodes",
            meta: {
                title: "Nodes",
            }
        },
        {
            path: "/nodes/detail/:name",
            component: () => import("@/business/cluster/nodes/detail"),
            name: "NodeDetail",
            props: true,
            hidden: true,
            meta: {
                activeMenu: "/nodes",
            }
        },
      {
        path: "/nodes/edit/:name",
        component: () => import("@/business/cluster/nodes/edit"),
        name: "NodeEdit",
        props: true,
        hidden: true,
        meta: {
          activeMenu: "/nodes",
        }
      },
        {
            path: "/namespaces",
            component: () => import("@/business/cluster/namespaces"),
            name: "Namespaces",
            requirePermission: {
                apiGroup: "",
                resource: "namespaces",
                verb: "list",
                scope: "cluster"
            },
            meta: {
                title: "Namespaces"
            }
        },
        {
            path: "/namespaces/create",
            component: () => import("@/business/cluster/namespaces/create"),
            name: "NamespaceCreate",
            hidden: true,
            meta: {
                activeMenu: "/namespaces"
            }
        },
        {
            path: "/namespaces/detail/:name",
            props: true,
            component: () => import("@/business/cluster/namespaces/detail"),
            name: "NamespaceDetail",
            hidden: true,
            meta: {
                activeMenu: "/namespaces"
            }
        },
        {
            path: "/namespaces/edit/:name",
            props: true,
            component: () => import("@/business/cluster/namespaces/edit"),
            name: "NamespaceEdit",
            hidden: true,
            meta: {
                activeMenu: "/namespaces"
            }
        },
        {
            path: "/events",
            component: () => import("@/business/cluster/events"),
            name: "events",
            requirePermission: {
                apiGroup: "",
                resource: "events",
                verb: "list",
                scope: "cluster"
            },
            meta: {
                title: "Events",
            }
        },

        {
            path: "/customResources",
            component: () => import("@/business/custom-resource"),
            name: "CustomResourceDefinitions",
            requirePermission: {
                apiGroup: "apiextensions.k8s.io",
                resource: "customresourcedefinitions",
                verb: "list",
                scope:"cluster"
            },
            meta: {
                title: "CRD",
            },
        },
        {
            path: "/customResources/:name/detail",
            component: () => import("@/business/custom-resource/detail"),
            hidden: true,
            props: true,
            name: "CustomResourceDefinitionDetail",
            requirePermission: {
                apiGroup: "apiextensions.k8s.io",
                resource: "customresourcedefinitions",
                verb: "get",
                scope:"namespace"
            },
            meta: {
                activeMenu: "/customResources",
            },
        },
        {
            path: "/customResources/:name/edit",
            component: () => import("@/business/custom-resource/edit"),
            hidden: true,
            props: true,
            name: "CustomResourceDefinitionEdit",
            requirePermission: {
                apiGroup: "apiextensions.k8s.io",
                resource: "customresourcedefinitions",
                verb: "update",
                scope:"namespace"
            },
            meta: {
                activeMenu: "/customResources",
            },
        },
        {
          path: "/resource/:group/:names/:version/:name/edit",
          component: () => import("@/business/custom-resource/cr/edit"),
          hidden: true,
          props: true,
          name: "CustomResourceEdit",
          requirePermission: {
            apiGroup: "apiextensions.k8s.io",
            resource: "customresourcedefinitions",
            verb: "update",
            scope:"namespace"
          },
          meta: {
            activeMenu: "/customResources",
          },
        },
    ]
}

export default Clusters
