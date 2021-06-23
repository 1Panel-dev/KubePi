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
      name: "ConfigMaps",
      meta: {
        title: "Config Maps",
      }
    },
    {
      path: "/configmaps/:namespace/create",
      component: () => import("@/business/configuration/config-maps/create"),
      name: "ConfigMapCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      path: "/configmaps/:namespace/:name",
      component: () => import("@/business/configuration/config-maps/detail"),
      name: "ConfigMapDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      path: "/configmaps/:namespace/:name/edit",
      component: () => import("@/business/configuration/config-maps/edit"),
      name: "ConfigMapEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      path: "/secrets",
      component: () => import("@/business/configuration/secrets"),
      name: "Secrets",
      meta: {
        title: "Secrets",
      }
    },
    {
      path: "/resourcequotas",
      component: () => import("@/business/configuration/resource-quotas"),
      name: "ResourceQuotas",
      meta: {
        title: "Resource Quotas",
      }
    },
    {
      path: "/limitranges",
      component: () => import("@/business/configuration/limit-ranges"),
      name: "LimitRanges",
      meta: {
        title: "Limit Ranges",
      }
    },
    {
      path: "/hpa",
      component: () => import("@/business/configuration/hpa"),
      name: "HPA",
      meta: {
        title: "Horizontal Pod Autoscaler",
      }
    }
  ]
}

export default Configuration
