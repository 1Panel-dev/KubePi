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
  ]
}

export default Clusters
