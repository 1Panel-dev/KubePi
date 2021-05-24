import Layout from "@/business/app-layout/horizontal-layout"

const Cluster = {
  path: "/cluster",
  sort: 1,
  component: Layout,
  name: "Cluster",
  meta: {
    title: "business.cluster.cluster",
    icon: "iconfont iconkubernets"
  },
  children: [
    // {
    //   path: "/cluster",
    //   component: () => import("@/business/cluster"),
    //   name: "Cluster",
    //   meta: {
    //     title: "business.cluster.cluster",
    //     icon: "iconfont iconkubernets"
    //   },
    // },
    {
      path: "/nodes",
      component: () => import("@/business/cluster/nodes"),
      name: "Nodes",
      meta: {
        title: "Nodes",
        icon: "iconfont iconnodes"
      }
    },
    {
      path: "/namespaces",
      component: () => import("@/business/cluster/namespaces"),
      name: "namespaces",
      meta: {
        title: "Namespaces",
        icon: "iconfont iconnamespace"
      }
    }
  ]
}

export default Cluster
