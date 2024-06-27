import Layout from "@/business/app-layout/horizontal-layout"

const Configuration = {
  path: "/configuration",
  parent: true,
  sort: 4,
  component: Layout,
  name: "Configuration",
  meta: {
    title: "business.configuration.configuration",
    icon: "iconfont iconconfiguration",
    scope:"namespace"
  },
  children: [
    {
      path: "/configmaps",
      requirePermission: {
        apiGroup: "",
        resource: "configmaps",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/config-maps"),
      name: "ConfigMaps",
      meta: {
        title: "Config Maps",
      }
    },
    {
      path: "/configmaps/create",
      requirePermission: {
        apiGroup: "",
        resource: "configmaps",
        verb: "create",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "configmaps",
        verb: "get",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/config-maps/detail"),
      name: "ConfigMapDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/configmaps"
      }
    },
    {
      requirePermission: {
        apiGroup: "",
        resource: "configmaps",
        verb: "update",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "secrets",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/secrets"),
      name: "Secrets",
      meta: {
        title: "Secrets",
      }
    },
    {
      path: "/secrets/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "",
        resource: "secrets",
        verb: "get",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "secrets",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/secrets/create"),
      name: "SecretCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/secrets"
      }
    },
    {
      path: "/secrets/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "",
        resource: "secrets",
        verb: "update",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "resourcequotas",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/resource-quotas"),
      name: "ResourceQuotas",
      meta: {
        title: "Resource Quotas",
      }
    },
    {
      path: "/:namespace/resourcequotas/detail/:name",
      requirePermission: {
        apiGroup: "",
        resource: "resourcequotas",
        verb: "get",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "resourcequotas",
        verb: "create",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "resourcequotas",
        verb: "update",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "limitranges",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/limit-ranges"),
      name: "LimitRanges",
      meta: {
        title: "Limit Ranges",
      }
    },
    {
      path: "/limitranges/create",
      requirePermission: {
        apiGroup: "",
        resource: "limitranges",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/limit-ranges/create"),
      name: "LimitRangeCreate",
      hidden: true,
      meta: {
        activeMenu: "/limitranges"
      }
    },
    {
      path: "/limitranges/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "",
        resource: "limitranges",
        verb: "get",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "",
        resource: "limitranges",
        verb: "update",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "autoscaling",
        resource: "horizontalpodautoscalers",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/hpa"),
      name: "HPA",
      meta: {
        title: "HPA",
      }
    },
    {
      path: "/horizontalpodautoscalers/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "autoscaling",
        resource: "horizontalpodautoscalers",
        verb: "get",
        scope:"namespace"
      },
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
      requirePermission: {
        apiGroup: "autoscaling",
        resource: "horizontalpodautoscalers",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/hpa/create"),
      name: "HPACreate",
      hidden: true,
      meta: {
        activeMenu: "/horizontalpodautoscalers"
      }
    },
    {
      path: "/horizontalpodautoscalers/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "autoscaling",
        resource: "horizontalpodautoscalers",
        verb: "update",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/hpa/edit"),
      name: "HPAEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/horizontalpodautoscalers"
      }
    },
    {
      path: "/poddisruptionbudgets",
      requirePermission: {
        apiGroup: "policy",
        resource: "poddisruptionbudgets",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/configuration/pdb"),
      name: "PDBs",
      meta: {
        title: "Pod Disruption Budgets",
      }
    },
    {
      path: "/poddisruptionbudgets/create",
      requirePermission: {
        apiGroup: "policy",
        resource: "poddisruptionbudgets",
        verb: "create",
        scope:"namespace"
      },
      hidden: true,
      props: true,
      component: () => import("@/business/configuration/pdb/create"),
      name: "PDBCreate",
      meta: {
        title: "Pod Disruption Budget",
        activeMenu: "/poddisruptionbudgets"
      }
    },
    {
      path: "/poddisruptionbudgets/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "policy",
        resource: "poddisruptionbudgets",
        verb: "update",
        scope:"namespace"
      },
      hidden: true,
      props: true,
      component: () => import("@/business/configuration/pdb/edit"),
      name: "PDBEdit",
      meta: {
        title: "Pod Disruption Budget",
        activeMenu: "/poddisruptionbudgets"
      }
    },
    {
      path: "/poddisruptionbudgets/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "policy",
        resource: "poddisruptionbudgets",
        verb: "get",
        scope:"namespace"
      },
      hidden: true,
      props: true,
      component: () => import("@/business/configuration/pdb/detail"),
      name: "PDBDetail",
      meta: {
        title: "Pod Disruption Budget",
        activeMenu: "/poddisruptionbudgets"
      }
    },
{
      path: "/mutatingwebhookconfigurations",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "mutatingwebhookconfigurations",
        verb: "list",
        scope:"cluster"
      },
      component: () => import("@/business/configuration/mutatingwebhookconfigurations"),
      name: "Mutatingwebhookconfigurations",
      meta: {
        title: "Mutatingwebhookconfigurations",
      }
    },
    {
      path: "/mutatingwebhookconfigurations/create",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "mutatingwebhookconfigurations",
        verb: "create",
      },
      component: () => import("@/business/configuration/mutatingwebhookconfigurations/create"),
      name: "MutatingWebhookConfigurationCreate",
      hidden: true,
      meta: {
        title: "MutatingWebhookConfiguration Create",
        activeMenu: "/mutatingwebhookconfigurations",
      }
    },
    {
      path: "/mutatingwebhookconfigurations/edit/:name",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "mutatingwebhookconfigurations",
        verb: "update",
      },
      component: () => import("@/business/configuration/mutatingwebhookconfigurations/edit"),
      name: "MutatingwebhookconfigurationEdit",
      hidden: true,
      props: true,
      meta: {
        title: "MutatingwebhookconfigurationEdit",
        activeMenu: "/mutatingwebhookconfigurations",
      }
    },
    {
      path: "/validatingwebhookconfigurations",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "validatingwebhookconfigurations",
        verb: "list",
        scope:"cluster"
      },
      component: () => import("@/business/configuration/validatingwebhookconfigurations"),
      name: "Validatingwebhookconfigurations",
      meta: {
        title: "Validatingwebhookconfigurations",
      }
    },
    {
      path: "/validatingwebhookconfigurations/create",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "validatingwebhookconfigurations",
        verb: "create",
      },
      component: () => import("@/business/configuration/validatingwebhookconfigurations/create"),
      name: "ValidatingwebhookconfigurationCreate",
      hidden: true,
      meta: {
        title: "Validatingwebhookconfiguration Create",
        activeMenu: "/validatingwebhookconfigurations",
      }
    },
    {
      path: "/validatingwebhookconfigurations/edit/:name",
      requirePermission: {
        apiGroup: "admissionregistration.k8s.io",
        resource: "validatingwebhookconfigurations",
        verb: "update",
      },
      component: () => import("@/business/configuration/validatingwebhookconfigurations/edit"),
      name: "ValidatingwebhookconfigurationEdit",
      hidden: true,
      props: true,
      meta: {
        title: "ValidatingwebhookconfigurationEdit",
        activeMenu: "/validatingwebhookconfigurations",
      }
    }
  ]
}

export default Configuration
