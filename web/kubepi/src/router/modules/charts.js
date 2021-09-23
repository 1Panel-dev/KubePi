import Layout from "@/business/app-layout/horizontal-layout"

const Charts = {
  path: "/charts",
  sort: 3,
  component: Layout,
  name: "ChartManagement",
  // requirePermission: {
  //   resource: "charts",
  //   verb: "list"
  // },
  meta: {
    title: "business.chart.chart",
    icon: "el-icon-house",
  },
  children: [
    {
      path: "",
      component: () => import("@/business/chart-management"),
      name: "Charts",
      meta: {
        title: "business.chart.chart",
        activeMenu: "/charts",
      }
    },
    {
      path: "create",
      component: () => import("@/business/chart-management/create"),
      name: "ChartCreate",
      hidden: true,
      meta: {
        activeMenu: "/charts",
      }
    },
    {
      path: "edit/:name",
      component: () => import("@/business/chart-management/edit"),
      name: "ChartEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/charts",
      }
    },
  ]
}

export default Charts

