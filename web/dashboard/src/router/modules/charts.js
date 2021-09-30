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
      path: "/repos",
      component: () => import("@/business/market-place/repo"),
      name: "Repos",
      meta: {
        title: "business.chart.chart",
      }
    },
    {
      path: "create",
      component: () => import("@/business/market-place/repo/create"),
      name: "RepoCreate",
      hidden: true,
      meta: {
        activeMenu: "/repos",
      }
    },
    {
      path: "edit/:name",
      component: () => import("@/business/market-place/repo/edit"),
      name: "RepoEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/repos",
      }
    },
    {
      path: "/apps",
      component: () => import("@/business/market-place/chart"),
      name: "Charts",
      meta: {
        title: "business.chart.app",
      }
    },
    {
      path: "/apps/:repo/:name",
      component: () => import("@/business/market-place/chart/detail"),
      name: "ChartDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/apps",
      }
    },
  ]
}

export default Charts

