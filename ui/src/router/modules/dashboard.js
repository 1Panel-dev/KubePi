import Layout from "@/business/app-layout/horizontal-layout"

const Dashboard = {
  path: "/dashboard",
  sort: 0,
  component: Layout,
  name: "Dashboard",
  children: [
    {
      path: "/dashboard",
      component: () => import("@/business/dashborad"),
      name: "Dashboard",
      meta: {
        title: "business.dashboard.dashboard",
        icon: "iconfont iconkubernets"
      },
    },
  ]
}

export default Dashboard
