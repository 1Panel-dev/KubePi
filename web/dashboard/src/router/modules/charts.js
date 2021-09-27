import Layout from "@/business/app-layout/horizontal-layout"

const Charts = {
  path: "",
  sort: 7,
  component: Layout,
  name: "ChartManagement",
  // requirePermission: {
  //   resource: "charts",
  //   verb: "list"
  // },
  meta: {
    title: "business.chart.marketPlace",
    icon: "el-icon-s-shop",
  },
  children: [
    {
      path: "/charts",
      component: () => import("@/business/market-place/chart"),
      name: "Charts",
      meta: {
        title: "business.chart.chart",
      }
    },
    {
      path: "create",
      component: () => import("@/business/market-place/chart/create"),
      name: "ChartCreate",
      hidden: true,
      meta: {
        activeMenu: "/charts",
      }
    },
    {
      path: "edit/:name",
      component: () => import("@/business/market-place/chart/edit"),
      name: "ChartEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/charts",
      }
    },
    {
      path: "/apps",
      component: () => import("@/business/market-place/apps"),
      name: "Apps",
      meta: {
        title: "business.chart.app",
      }
    },
  ]
}

export default Charts

