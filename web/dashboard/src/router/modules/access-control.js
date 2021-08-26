import Layout from "@/business/app-layout/horizontal-layout";

const AccessControl = {
  path: "/accesscontrol",
  parent: true,
  sort: 6,
  component: Layout,
  name: "Access Control",
  meta: {
    title: "business.access_control.access_control",
    icon: "iconfont iconaccesscontrol",
    global: false,
  },
  children: [
    {
      path: "/serviceaccounts",
      requirePermission: {
        apiGroup: "",
        resource: "serviceaccounts",
        verb: "list",
      },
      component: () => import("@/business/access-control/service-accounts"),
      name: "ServiceAccounts",
      meta: {
        title: "Service Accounts",
        global: false,
      },
    },
    {
      path: "/serviceaccounts/create",
      requirePermission: {
        apiGroup: "",
        resource: "serviceaccounts",
        verb: "create",
      },
      component: () =>
        import("@/business/access-control/service-accounts/create"),
      name: "ServiceAccountCreate",
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false,
      },
    },
    {
      path: "/serviceaccounts/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "",
        resource: "serviceaccounts",
        verb: "get",
      },
      component: () =>
        import("@/business/access-control/service-accounts/detail"),
      name: "ServiceAccountDetail",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false,
      },
    },
    {
      path: "/serviceaccounts/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "",
        resource: "serviceaccounts",
        verb: "update",
      },
      component: () =>
        import("@/business/access-control/service-accounts/edit"),
      name: "ServiceAccountEdit",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/serviceaccounts",
        global: false,
      },
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
        global: false,
      },
    },
    {
      path: "/roles/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "roles",
        verb: "get",
      },
      component: () => import("@/business/access-control/roles/detail"),
      name: "RoleDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/roles",
        global: false,
      },
    },
    {
      path: "/roles/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "roles",
        verb: "update",
      },
      component: () => import("@/business/access-control/roles/edit"),
      name: "RoleEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/roles",
        global: false,
      },
    },
    {
      path: "/roles/create",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "roles",
        verb: "create",
      },
      component: () => import("@/business/access-control/roles/create"),
      name: "RoleCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/roles",
        global: false,
      },
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
        global: false,
      },
    },
    {
      path: "/rolebindings/create",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "rolebindings",
        verb: "create",
      },
      component: () => import("@/business/access-control/role-bindings/create"),
      name: "RoleBindingCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/rolebindings",
        global: false,
      },
    },
    {
      path: "/rolebindings/:namespace/:name/edit",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "rolebindings",
        verb: "update",
      },
      component: () => import("@/business/access-control/role-bindings/edit"),
      name: "RoleBindingEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/rolebindings",
        global: false,
      },
    },
    {
      path: "/rolebindings/:namespace/:name/detail",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "rolebindings",
        verb: "get",
      },
      component: () => import("@/business/access-control/role-bindings/detail"),
      name: "RoleBindingDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/rolebindings",
        global: false,
      },
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
        title: "Cluster Roles",
        global: false,
      },
    },
    {
      path: "/clusterroles/:name/detail",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "get",
      },
      component: () => import("@/business/access-control/cluster-roles/detail"),
      name: "ClusterRoleDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterroles",
        global: false,
      },
    },
    {
      path: "/clusterroles/:name/edit",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "update",
      },
      component: () => import("@/business/access-control/cluster-roles/edit"),
      name: "ClusterRoleEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterroles",
        global: false,
      },
    },
    {
      path: "/clusterroles/create",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "create",
      },
      component: () => import("@/business/access-control/cluster-roles/create"),
      name: "ClusterRoleCreate",
      props: true,
      hidden: true,
      meta: {
        activeMenu: "/clusterroles",
        global: false,
      },
    },

    {
      path: "/clusterrolebindings",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "list",
      },
      component: () =>
        import("@/business/access-control/cluster-role-bindings"),
      name: "ClusterRoleBindings",
      meta: {
        title: "Cluster Role Bindings",
        global: false,
      },
    },
    {
      path: "/clusterrolebindings/create",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "create",
      },
      component: () =>
        import("@/business/access-control/cluster-role-bindings/create"),
      name: "ClusterRoleBindingCreate",
      hidden: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false,
      },
    },
    {
      path: "/clusterrolebindings/:name/detail",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "get",
      },
      component: () =>
        import("@/business/access-control/cluster-role-bindings/detail"),
      name: "ClusterRoleBindingDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false,
      },
    },
    {
      path: "/clusterrolebindings/:name/edit",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "update",
      },
      component: () =>
        import("@/business/access-control/cluster-role-bindings/edit"),
      name: "ClusterRoleBindingEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/clusterrolebindings",
        global: false,
      },
    },
  ],
};

export default AccessControl;
