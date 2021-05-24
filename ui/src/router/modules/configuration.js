import Layout from "@/business/app-layout/horizontal-layout"

const Configuration = {
  path: "/configuration",
  sort: 2,
  component: Layout,
  name: "Configuration",
  meta: {
    title: "business.configuration.configuration",
    icon: "iconfont iconconfiguration"
  },
  children: [
    {
      path: "/configmaps",
      component: () => import("@/business/configuration/config-maps"),
      name: "Configmaps",
      meta: {
        title: "Config Maps",
      }
    },
    {
      path: "/secrets",
      component: () => import("@/business/configuration/secrets"),
      name: "Secrets",
      meta: {
        title: "Secrets",
      }
    }
  ]
}

export default Configuration
