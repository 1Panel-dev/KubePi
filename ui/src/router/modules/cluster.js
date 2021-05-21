import Layout from "@/business/app-layout/horizontal-layout"

const Cluster = {
  path: "/cluster",
  sort: 1,
  component: Layout,
  name: "Cluster",
  children: [
    {
      path: "/cluster",
      component: () => import("@/business/cluster"),
      name: "Dashboard",
      meta: {
        title: "Cluster",
        icon: "iconfont iconkubernets"
      },
    },
  ]
}

export default Cluster
