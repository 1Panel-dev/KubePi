import Layout from "@/business/app-layout/horizontal-layout"

const Network = {
  path: "/network",
  sort: 3,
  parent: true,
  component: Layout,
  name: "Network",
  meta: {
    title: "business.network.service_discovery",
    icon: "iconfont iconnetwork",
    scope:"namespace"
  },
  children: [
    {
      path: "/services",
      requirePermission: {
        apiGroup: "",
        resource: "services",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/network/services"),
      name: "Services",
      meta: {
        title: "Services",
      }
    },
    {
      path: "/services/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "",
        resource: "services",
        verb: "get",
        scope:"namespace"
      },
      component: () => import("@/business/network/services/detail"),
      name: "ServiceDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/services"
      }
    },
    {
      path: "/services/create",
      requirePermission: {
        apiGroup: "",
        resource: "services",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/network/services/create"),
      name: "ServiceCreate",
      hidden: true,
      meta: {
        activeMenu: "/services"
      }
    },
    {
      path: "/services/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "",
        resource: "services",
        verb: "update",
        scope:"namespace"
      },
      component: () => import("@/business/network/services/edit"),
      name: "ServiceEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/services"
      }
    },

    {
      path: "/endpoints",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/network/endpoints"),
      name: "Endpoints",
      meta: {
        title: "Endpoints",
      }
    },
    {
      path: "/endpoints/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "get",
        scope:"namespace"
      },
      component: () => import("@/business/network/endpoints/detail"),
      name: "EndpointDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/endpoints"
      }
    },
    {
      path: "/endpoints/create",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/network/endpoints/create"),
      name: "EndpointCreate",
      hidden: true,
      meta: {
        activeMenu: "/endpoints"
      }
    },
    {
      path: "/endpoints/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "update",
        scope:"namespace"
      },
      component: () => import("@/business/network/endpoints/edit"),
      name: "EndpointEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/endpoints"
      }
    },

    {
      path: "/ingresses",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "ingresses",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/network/ingresses"),
      name: "Ingresses",
      meta: {
        title: "Ingresses",
      }
    },
    {
      path: "/ingresses/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "ingresses",
        verb: "get",
        scope:"namespace"
      },
      component: () => import("@/business/network/ingresses/detail"),
      name: "IngressDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/ingresses"
      }
    },
    {
      path: "/ingresses/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "ingresses",
        verb: "update",
        scope:"namespace"
      },
      component: () => import("@/business/network/ingresses/operate"),
      name: "IngressEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/ingresses"
      }
    },
    {
      path: "/ingresses/create",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "ingresses",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/network/ingresses/operate"),
      name: "IngressCreate",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/ingresses"
      }
    },

    {
      path: "/networkpolicies",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "networkpolicies",
        verb: "list",
        scope:"namespace"
      },
      component: () => import("@/business/network/network-policies"),
      name: "NetworkPolicies",
      meta: {
        title: "Network Policies",
      }
    },
    {
      path: "/networkpolicies/create",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "networkpolicies",
        verb: "create",
        scope:"namespace"
      },
      component: () => import("@/business/network/network-policies/create"),
      name: "NetworkPolicyCreate",
      hidden: true,
      meta: {
        activeMenu: "/networkpolicies"
      }
    },
    {
      path: "/networkpolicies/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "networkpolicies",
        verb: "get",
        scope:"namespace"
      },
      component: () => import("@/business/network/network-policies/detail"),
      name: "NetworkPolicyDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/networkpolicies"
      }
    },
    {
      path: "/networkpolicies/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "networkpolicies",
        verb: "update",
        scope:"namespace"
      },
      component: () => import("@/business/network/network-policies/edit"),
      name: "NetworkPolicyEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/networkpolicies"
      }
    },
  ]
}

export default Network
