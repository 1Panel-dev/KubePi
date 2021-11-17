import Layout from "@/business/app-layout/horizontal-layout"

const UserManagement = {
  path: "/user-management",
  sort: 2,
  component: Layout,
  name: "UserManagement",
  requirePermission: {
    resource: "users",
    verb: "list"
  },
  meta: {
    title: "business.user.user_management",
    icon: "iconfont iconyonghuguanli",
    global: true
  },
  children: [
    {
      path: "users",
      component: () => import("@/business/user-management/user"),
      name: "Users",
      requirePermission: {
        resource: "users",
        verb: "list"
      },
      meta: {
        title: "business.user.user",
      }
    },
    {
      path: "users/create",
      component: () => import("@/business/user-management/user/create"),
      name: "UserCreate",
      hidden: true,
      meta: {
        activeMenu: "/user-management/users",
      }
    },
    {
      props: true,
      path: "users/edit/:name",
      component: () => import("@/business/user-management/user/edit"),
      name: "UserEdit",
      hidden: true,
      meta: {
        activeMenu: "/user-management/users",
      }
    },
    {
      path: "roles",
      component: () => import("@/business/user-management/role"),
      name: "Roles",
      requirePermission: {
        resource: "roles",
        verb: "list"
      },
      meta: {
        title: "business.user.role",
        global: true
      }
    },
    {
      path: "roles/create",
      component: () => import("@/business/user-management/role/create"),
      name: "RoleCreate",
      hidden: true,
      meta: {
        activeMenu: "/user-management/roles",
      }
    },
    {
      props: true,
      path: "roles/edit/:name",
      component: () => import("@/business/user-management/role/edit"),
      name: "RoleEdit",
      hidden: true,
      meta: {
        activeMenu: "/user-management/roles",
      },
    },
    {
      props: true,
      path: "roles/detail/:name",
      component: () => import("@/business/user-management/role/detail"),
      name: "RoleDetail",
      hidden: true,
      meta: {
        activeMenu: "/user-management/roles",
      },
    },
    {
      path: "ldap",
      component: () => import("@/business/user-management/ldap"),
      name: "LDAP",
      requirePermission: {
        resource: "ldap",
        verb: "list"
      },
      meta: {
        title: "business.user.ldap",
      },
    }

  ]
}

export default UserManagement
