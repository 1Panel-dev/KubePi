import Layout from "@/business/app-layout/horizontal-layout"

const Dashboard = {
    path: "/dashboard",
    sort: 0,
    global: true,
    component: Layout,
    name: "Dashboards",
    children: [
        {
            path: "/dashboard",
            component: () => import("@/business/dashboard"),
            name: "Dashboard",
            meta: {
                title: "business.dashboard.dashboard",
                icon: "el-icon-data-line",
                roles: ['ADMIN']
            },
        },
    ]
}

export default Dashboard
