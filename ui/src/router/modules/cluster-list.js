import Layout from "@/business/app-layout/horizontal-layout"

const ClusterList = {
  path: "/clusterlist",
  sort: 1,
  component: Layout,
  name: "ClusterList",
  meta: {
    title: "business.cluster.list",
    icon: "iconfont iconkubernets",
    global: true
  },
  children: [
    {
      path: "/clusterlist",
      component: () => import("@/business/cluster-list"),
      name: "ClusterList",
      meta: {
        title: "business.cluster.list",
        global: true
      }
    },
    {
      path: "/import",
      component: () => import("@/business/cluster-list/import"),
      hidden: true,
      name: "ClusterImport",
      meta: {
        title: "business.cluster.import",
        global: true
      }
    },
  ]
}

export default ClusterList
