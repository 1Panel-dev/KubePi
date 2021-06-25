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
      component: () => import("@/business/workloads/deployments"),
      name: "Deployments",
      meta: {
        title: "Deployments",
      }
    },
    {
      path: "/deployments/create",
      name: "DeploymentCreate",
      hidden: true,
      component: () => import("@/business/workloads/deployments/create"),
      meta: {
        activeMenu: "/deployments",
      }
    },
    {
      path: "/deployments/edit/:cluster/:namespace/:name/:yamlShow",
      name: "DeploymentEdit",
      hidden: true,
      component: () => import("@/business/workloads/deployments/edit"),
      meta: {
        activeMenu: "/deployments",
      }
    },
  ]
}

export default Workloads
