import Layout from "@/business/app-layout/horizontal-layout"

const ClusterList = {
  path: "/userList",
  sort: 2,
  component: Layout,
  name: "UserList",
  meta: {
    title: "business.user.menu",
    icon: "el-icon-user",
    global: true
  },
  children: [
    {
      path: "/users",
      component: () => import("@/business/users/list"),
      name: "UserList",
      meta: {
        title: "business.user.user",
        global: true
      }
    },
    {
      path: "/roles",
      component: () => import("@/business/users/role"),
      name: "Roles",
      meta: {
        title: "business.user.role",
        global: true
      }
    },
  ]
}

export default ClusterList
