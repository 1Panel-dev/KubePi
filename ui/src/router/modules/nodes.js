import Layout from "@/business/app-layout/horizontal-layout"

const Nodes = {
  path: "/nodes",
  sort: 2,
  component: Layout,
  name: "Nodes",
  children: [
    {
      path: "/nodes",
      component: () => import("@/business/node"),
      name: "NodeList",
      meta: {
        title: "Nodes",
        icon: "iconfont iconnode"
      },
    },
  ]
}

export default Nodes
