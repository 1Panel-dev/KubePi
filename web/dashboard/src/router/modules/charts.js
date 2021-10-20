import Layout from "@/business/app-layout/horizontal-layout"

const Charts = {
  path: "/charts",
  sort: 7,
  component: Layout,
  name: "ChartManagement",
  requirePermission: {
    apiGroup: "kubepi.io",
    resource: "appmarkets",
    verb: "list",
    scope:"cluster"
  },
  meta: {
    title: "business.chart.marketPlace",
    icon: "iconfont iconmarket",
  },
  children: [
    {
      path: "/repos",
      component: () => import("@/business/market-place/repo"),
      name: "Repos",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "list",
        scope:"cluster"
      },
      meta: {
        title: "business.chart.chart",
      }
    },
    {
      path: "/repos/create",
      component: () => import("@/business/market-place/repo/create"),
      name: "RepoCreate",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "create",
        scope:"cluster"
      },
      hidden: true,
      meta: {
        activeMenu: "/repos",
      }
    },
    {
      path: "/repos/edit/:name",
      component: () => import("@/business/market-place/repo/edit"),
      name: "RepoEdit",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "update",
        scope:"cluster"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/repos",
      }
    },
    {
      path: "/charts",
      component: () => import("@/business/market-place/chart"),
      name: "Charts",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "list",
        scope:"cluster"
      },
      meta: {
        title: "business.chart.app",
      }
    },
    {
      path: "/charts/:repo/:name",
      component: () => import("@/business/market-place/chart/detail"),
      name: "ChartDetail",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "get",
        scope:"cluster"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/charts",
      }
    },
    {
      path: "/apps",
      component: () => import("@/business/market-place/app"),
      name: "Apps",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "list",
        scope:"cluster"
      },
      meta: {
        title: "business.chart.app_installed",
      }
    },
    {
      path: "/apps/:name",
      component: () => import("@/business/market-place/app/update"),
      name: "AppUpgrade",
      requirePermission: {
        apiGroup: "kubepi.io",
        resource: "appmarkets",
        verb: "update",
        scope:"cluster"
      },
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/apps",
      }
    },
  ]
}

export default Charts

