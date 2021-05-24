import Layout from "@/business/app-layout/horizontal-layout"

const Clusters = {
  path: "/cluster",
  sort: 1,
  component: Layout,
  name: "Cluster",
  meta: {
    title: "business.cluster.cluster",
    icon: "iconfont iconkubernets"
  },
  children: [
    {
      path: "/nodes",
      component: () => import("@/business/cluster/nodes"),
      name: "Nodes",
      meta: {
        title: "Nodes",
      }
    },
    {
      path: "/namespaces",
      component: () => import("@/business/cluster/namespaces"),
      name: "namespaces",
      meta: {
        title: "Namespaces",
      }
    }
  ]
}

export default Clusters
