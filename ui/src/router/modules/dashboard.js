import Layout from "@/business/app-layout/horizontal-layout"

const Dashboard = {
  path: "/dashboard",
  sort: 0,
  component: Layout,
  name: "Dashboard",
  children: [
    {
      path: "/dashboard",
      component: () => import("@/business/dashboard"),
      name: "Dashboard",
      meta: {
        title: "business.dashboard.dashboard",
        icon: "el-icon-data-line"
      },
    },
  ]
}

export default Dashboard
