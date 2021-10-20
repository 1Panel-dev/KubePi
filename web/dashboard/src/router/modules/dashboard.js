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
                icon: "iconfont icongailan"
            },
        },
    ]
}

export default Dashboard
