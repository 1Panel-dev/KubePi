import Layout from "@/business/app-layout/horizontal-layout"

const Workloads = {
  path: "/workloads",
  sort: 4,
  component: Layout,
  name: "Workloads",
  meta: {
    title: "business.workload.workload",
    icon: "iconfont iconworkload"
  },
  children: [
    {
      path: "/pods",
      component: () => import("@/business/workloads/pods"),
      name: "Pods",
      meta: {
        title: "Pods",
      }
    },
    {
      path: "/deployments",
      component: () => import("@/business/workloads/pods"),
      name: "Deployments",
      meta: {
        title: "Deployments",
      }
    }
  ]
}

export default Workloads
