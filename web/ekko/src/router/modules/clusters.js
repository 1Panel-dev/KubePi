import Layout from "@/business/app-layout/horizontal-layout"

const Clusters = {
    path: "/clusters",
    sort: 1,
    component: Layout,
    name: "ClusterManagement",
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
                    path: "overview",
                    props: true,
                    component: () => import("@/business/cluster-management/detail/overview"),
                    hidden: true,
                    name: "ClusterOverview",
                    meta: {
                        activeMenu: "/clusters",
                    },
                },
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
                    path: "clusterroles",
                    props: true,
                    component: () => import("@/business/cluster-management/detail/clusterroles"),
                    hidden: true,
                    name: "ClusterRoles",
                    meta: {
                        activeMenu: "/clusters",
                    },
                }
            ]

        }
    ]
}

export default Clusters
