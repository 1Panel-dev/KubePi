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
        path: "grafana",
        component: () => import("@/business/monitor/grafana"),
        name: "Grafana",
        requirePermission: {
          resource: "grafana",
          verb: "list"
        },
        meta: {
          title: "business.monitor.grafana.name",
        }
      },
      {
        path: "metrics",
        component: () => import("@/business/monitor/metrics"),
        name: "Metrics",
        requirePermission: {
          resource: "metrics",
          verb: "list"
        },
        meta: {
          title: "business.monitor.metrics.name",
        }
      },
      {
        path: "alert-rules",
        component: () => import("@/business/monitor/alert-rules"),
        name: "Alert Rules",
        requirePermission: {
          resource: "alert-rules",
          verb: "list"
        },
        meta: {
          title: "business.monitor.alert_rules.name",
        }
      },
    ]
}

export default Monitor;
