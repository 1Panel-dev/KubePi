import Layout from "@/business/app-layout/horizontal-layout"

const Accesscontrol = {
  path: "/accesscontrol",
  sort: 5,
  component: Layout,
  name: "Access Control",
  meta: {
    title: "business.access_control.access_control",
    icon: "iconfont iconaccesscontrol"
  },
  children: [
    {
      path: "/rolebindings",
      component: () => import("@/business/access-control/role-bindings"),
      name: "RoleBindings",
      meta: {
        title: "Role Bindings",
      }
    },
    {
      path: "/serviceaccounts",
      component: () => import("@/business/access-control/service-accounts"),
      name: "ServiceAccounts",
      meta: {
        title: "Service Accounts",
      }
    }
  ]
}

export default Accesscontrol
