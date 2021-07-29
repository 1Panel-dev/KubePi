import Layout from "@/business/app-layout/horizontal-layout"

const Network = {
  path: "/network",
  sort: 3,
  parent: true,
  component: Layout,
  name: "Network",
  meta: {
    title: "business.network.network",
    icon: "iconfont iconnetwork"
  },
  children: [
    {
      path: "/services",
      requirePermission: {
        apiGroup: "",
        resource: "services",
        verb: "list",
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
      path: "/ingresses",
      requirePermission: {
        apiGroup: "networking.k8s.io",
        resource: "ingresses",
        verb: "list",
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
      },
      component: () => import("@/business/network/ingresses/edit"),
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
      },
      component: () => import("@/business/network/ingresses/create"),
      name: "IngressCreate",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/ingresses"
      }
    },
    {
      path: "/endpoints",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "list",
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
        verb: "update",
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
      },
      component: () => import("@/business/network/endpoints/create"),
      name: "EndpointCreate",
      hidden: true,
      meta: {
        activeMenu: "/ingresses"
      }
    },
    {
      path: "/endpoints/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "",
        resource: "endpoints",
        verb: "update",
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
      path: "/networkpolicies",
      requirePermission: {
        apiGroup: "",
        resource: "networkpolicies",
        verb: "list",
      },
      component: () => import("@/business/network/network-policies"),
      name: "NetworkPolicies",
      meta: {
        title: "Network Policies",
      }
    },
    {
      path: "/networkpolicies/create",
      component: () => import("@/business/network/network-policies/create"),
      name: "NetworkPolicyCreate",
      hidden: true,
      meta: {
        activeMenu: "/networkpolicies"
      }
    },
    {
      path: "/networkpolicies/:namespace/:name/detail",
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
