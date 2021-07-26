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
      component: () => import("@/business/network/services/create"),
      name: "ServiceCreate",
      hidden: true,
      meta: {
        activeMenu: "/services"
      }
    },
    {
      path: "/services/:namespace/:name/edit",
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
    }
  ]
}

export default Network
