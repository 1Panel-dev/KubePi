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
      path: "/deployments/from/:cluster/:namespace/:name",
      name: "DeploymentForm",
      hidden: true,
      component: () => import("@/business/workloads/deployments/form"),
      meta: {
        activeMenu: "/deployments",
      }
    },

    {
      path: "/deployments/yaml/:cluster/:namespace/:name",
      name: "DeploymentYaml",
      hidden: true,
      component: () => import("@/business/workloads/deployments/yaml"),
      meta: {
        activeMenu: "/deployments",
      }
    },
  ]
}

export default Workloads
