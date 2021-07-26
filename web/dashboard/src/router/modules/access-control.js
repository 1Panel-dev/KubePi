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
        global: false
      }
    },
    {
      path: "/clusterrolebindings",
      requirePermission: {
        apiGroup: "rbac.authorization.k8s.io",
        resource: "clusterrolebindings",
        verb: "list",
      },
      component: () => import("@/business/access-control/cluster-role-bindings"),
      name: "RoleBindings",
      meta: {
        title: "ClusterRole Bindings",
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
    }
  ]
}

export default AccessControl
