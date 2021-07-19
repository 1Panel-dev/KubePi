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
      path: "/configmaps/create",
      component: () => import("@/business/configuration/config-maps/create"),
      name: "ConfigMapCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      path: "/:namespace/configmaps/detail/:name",
      component: () => import("@/business/configuration/config-maps/detail"),
      name: "ConfigMapDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      path: "/configmaps/edit/:namespace/:name",
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
      path: "/secrets/detail/:namespace/:name",
      component: () => import("@/business/configuration/secrets/detail"),
      name: "SecretDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/secrets"
      }
    },
    {
      path: "/secrets/create",
      component: () => import("@/business/configuration/secrets/create"),
      name: "SecretCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/secrets"
      }
    },
    {
      path: "/:namespace/secrets/edit/:name",
      component: () => import("@/business/configuration/secrets/edit"),
      name: "SecretEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/secrets"
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
      path: "/:namespace/resourcequotas/detail/:name",
      component: () => import("@/business/configuration/resource-quotas/detail"),
      name: "ResourceQuotaDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/resourcequotas"
      }
    },
    {
      path: "/resourcequotas/create",
      component: () => import("@/business/configuration/resource-quotas/create"),
      name: "ResourceQuotaCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/resourcequotas"
      }
    },
    {
      path: "/resourcequotas/edit/:namespace/:name",
      component: () => import("@/business/configuration/resource-quotas/edit"),
      name: "ResourceQuotaEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/resourcequotas"
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
      path: "/limitranges/create",
      component: () => import("@/business/configuration/limit-ranges/create"),
      name: "LimitRangeCreate",
      hidden: true,
      meta: {
        activeMenu: "/limitranges"
      }
    },
    {
      path: "/limitranges/detail/:namespace/:name",
      component: () => import("@/business/configuration/limit-ranges/detail"),
      name: "LimitRangeDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/limitranges"
      }
    },
    {
      path: "/limitranges/edit/:namespace/:name",
      component: () => import("@/business/configuration/limit-ranges/edit"),
      name: "LimitRangeEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/limitranges"
      }
    },
    {
      path: "/horizontalpodautoscalers",
      component: () => import("@/business/configuration/hpa"),
      name: "HPA",
      meta: {
        title: "Horizontal Pod Autoscaler",
      }
    },
    {
      path: "/horizontalpodautoscalers/:namespace/:name/detail",
      component: () => import("@/business/configuration/hpa/detail"),
      name: "HPADetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/horizontalpodautoscalers"
      }
    },
    {
      path: "/horizontalpodautoscalers/create",
      component: () => import("@/business/configuration/hpa/create"),
      name: "HPACreate",
      hidden: true,
      meta: {
        activeMenu: "/horizontalpodautoscalers"
      }
    }
  ]
}

export default Configuration
