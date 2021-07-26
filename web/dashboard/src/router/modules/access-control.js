import Layout from "@/business/app-layout/horizontal-layout"

const AccessControl = {
  path: "/accesscontrol",
  parent: true,
  sort: 5,
  component: Layout,
  name: "Access Control",
  meta: {
    title: "business.access_control.access_control",
    icon: "iconfont iconaccesscontrol",
    global: false
  },
  children: [
    {
      path: "/clusterrolebindings",
      // requirePermission: {
      //   apiGroup: "rbac.authorization.k8s.io",
      //   resource: "clusterrolebindings",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/cluster-role-bindings"),
      name: "ClusterRoleBindings",
      meta: {
        title: "ClusterRoleBindings",
        global: false
      }
    },
    {
      path: "/clusterrolebindings/create",
      // requirePermission: {
      //   apiGroup: "rbac.authorization.k8s.io",
      //   resource: "clusterrolebindings",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/cluster-role-bindings/create"),
      name: "ClusterRoleBindingCreate",
      hidden: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false
      }
    },
    {
      path: "/clusterrolebindings/:name/detail",
      // requirePermission: {
      //   apiGroup: "rbac.authorization.k8s.io",
      //   resource: "clusterrolebindings",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/cluster-role-bindings/detail"),
      name: "ClusterRoleBindingDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false
      }
    },
    {
      path: "/clusterrolebindings/:name/edit",
      // requirePermission: {
      //   apiGroup: "rbac.authorization.k8s.io",
      //   resource: "clusterrolebindings",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/cluster-role-bindings/edit"),
      name: "ClusterRoleBindingEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false
      }
    },
    {
      path: "/clusterroles",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterroles",
        verb: "list",
      },
      component: () => import("@/business/access-control/cluster-roles"),
      name: "ClusterRoles",
      meta: {
        title: "ClusterRoles",
        global: false
      }
    },
    {
      path: "/rolebindings",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "rolebindings",
        verb: "list",
      },
      component: () => import("@/business/access-control/role-bindings"),
      name: "RoleBindings",
      meta: {
        title: "Role Bindings",
        global: false
      }
    },
    {
      path: "/roles",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "roles",
        verb: "list",
      },
      component: () => import("@/business/access-control/roles"),
      name: "Roles",
      meta: {
        title: "Roles",
        global: false
      }
    },
    {
      path: "/serviceaccounts",
      // requirePermission: {
      //   apiGroup: "",
      //   resource: "serviceaccounts",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/service-accounts"),
      name: "ServiceAccounts",
      meta: {
        title: "Service Accounts",
        global: false
      }
    },
    {
      path: "/serviceaccounts/create",
      // requirePermission: {
      //   apiGroup: "",
      //   resource: "serviceaccounts",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/service-accounts/create"),
      name: "ServiceAccountCreate",
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false
      }
    },
    {
      path: "/serviceaccounts/:namespace/:name/detail",
      // requirePermission: {
      //   apiGroup: "",
      //   resource: "serviceaccounts",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/service-accounts/detail"),
      name: "ServiceAccountDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false
      }
    },
    {
      path: "/serviceaccounts/:namespace/:name/edit",
      // requirePermission: {
      //   apiGroup: "",
      //   resource: "serviceaccounts",
      //   verb: "list",
      // },
      component: () => import("@/business/access-control/service-accounts/edit"),
      name: "ServiceAccountEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false
      }
    },
    {
      path: "/podsecuritypolicy",
      requirePermission: {
        apiGroup: "policy",
        resource: "podsecuritypolicies",
        verb: "list",
      },
      component: () => import("@/business/access-control/pod-security-policies"),
      name: "PodSecurityPolicy",
      meta: {
        title: "Pod Security Policy",
        global: false
      }
    },
  ]
}

export default AccessControl
