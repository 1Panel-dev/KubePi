import Layout from "@/business/app-layout/horizontal-layout"

const Clusters = {
    path: "/clusters",
    sort: 1,
    component: Layout,
    name: "ClusterManagement",
    requirePermission: {
        resource: "clusters",
        verb: "list"
    },
    meta: {
        title: "business.cluster.list",
        icon: "iconfont iconkubernets",
    },
    children: [
        {
            path: "",
            component: () => import("@/business/cluster-management"),
            name: "Clusters",
            meta: {
                title: "business.cluster.list",
                activeMenu: "/clusters",
            }
        },
        {
            path: "create",
            component: () => import("@/business/cluster-management/create"),
            hidden: true,
            name: "ClusterCreate",
            meta: {
                activeMenu: "/clusters",
            }
        },
        {
            path: "detail/:name",
            props: true,
            component: () => import("@/business/cluster-management/detail"),
            hidden: true,
            name: "ClusterDetail",
            meta: {
                activeMenu: "/clusters",
            },
            children: [
                {
                    path: "members",
                    props: true,
                    component: () => import("@/business/cluster-management/detail/members"),
                    hidden: true,
                    name: "ClusterMembers",
                    meta: {
                        activeMenu: "/clusters",
                    },
                },
                {
                    path: "roles",
                    props: true,
                    component: () => import("@/business/cluster-management/detail/roles"),
                    hidden: true,
                    name: "ClusterRoles",
                    meta: {
                        activeMenu: "/clusters",
                    },
                    children: [
                        {
                            path: "",
                            redirect: "cluster"
                        },
                        {
                            path: "cluster",
                            props: true,
                            component: () => import("@/business/cluster-management/detail/roles/cluster"),
                            hidden: true,
                            name: "ClusterClusterRoles",
                            meta: {
                                activeMenu: "/clusters",
                            }
                        }, {
                            path: "namespace",
                            props: true,
                            component: () => import("@/business/cluster-management/detail/roles/namespace"),
                            hidden: true,
                            name: "NamespaceClusterRoles",
                            meta: {
                                activeMenu: "/clusters",
                            }
                        }
                    ],
                },
                {
                    path: "repos",
                    props: true,
                    component: () => import("@/business/cluster-management/detail/repos"),
                    hidden: true,
                    name: "ClusterRepos",
                    meta: {
                      activeMenu: "/clusters",
                    },
                }
            ]
        }
    ]
}

export default Clusters
