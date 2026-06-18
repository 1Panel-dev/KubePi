import Layout from "@/business/app-layout/horizontal-layout"

const Monitor = {
    path: "/monitor",
    sort: 5,
    component: Layout,
    name: "Monitor",
    requirePermission: {
      resource: "monitor",
      verb: "list"
    },
    meta: {
      title: "business.monitor.name",
      icon: "iconfont iconjiankong",
      global: true
    },
    children: [
      {
        path: "/monitor/grafana",
        component: () => import("@/business/monitor/grafana"),
        name: "Grafana",
        requirePermission: {
          resource: "monitor",
          verb: "list"
        },
        meta: {
          title: "business.monitor.grafana.name",
        }
      },
      {
        path: "/monitor/metrics",
        component: () => import("@/business/monitor/metrics"),
        name: "Metrics",
        requirePermission: {
          resource: "monitor",
          verb: "list"
        },
        meta: {
          title: "business.monitor.metrics.name",
        }
      },
      {
        path: "/monitor/metrics/create",
        component: () => import("@/business/monitor/metrics/operate"),
        name: "MetricsCreate",
        hidden: true,
        props: true,
        meta: {
          activeMenu: "/metrics",
        }
      },
      {
        path: "/monitor/metrics/:name/edit",
        component: () => import("@/business/monitor/metrics/operate"),
        name: "MetricsEdit",
        hidden: true,
        props: true,
        meta: {
          activeMenu: "/metrics",
        }
      },
      {
        path: "/monitor/metrics/:instance/detail",
        component: () => import("@/business/monitor/metrics/detail"),
        name: "MetricsDetail",
        hidden: true,
        props: true,
        meta: {
          activeMenu: "/metrics",
        }
      },
    ]
}

export default Monitor;
