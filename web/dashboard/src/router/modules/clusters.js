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
        }
    ]
}

export default Clusters
