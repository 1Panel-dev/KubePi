import Layout from "@/business/app-layout/horizontal-layout"

const UserManagement = {
    path: "/user-management",
    sort: 2,
    component: Layout,
    name: "UserManagement",
    meta: {
        title: "business.user.user_management",
        icon: "el-icon-user",
        global: true
    },
    children: [
        {
            path: "users",
            component: () => import("@/business/user-management/user"),
            name: "UserList",
            meta: {
                title: "business.user.user",
                global: true
            }
        },
        {
            path: "roles",
            component: () => import("@/business/user-management/role"),
            name: "Roles",
            meta: {
                title: "business.user.role",
                global: true
            }
        },
    ]
}

export default UserManagement
