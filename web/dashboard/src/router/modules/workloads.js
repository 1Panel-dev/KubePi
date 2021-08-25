import Layout from "@/business/app-layout/horizontal-layout";

const Workloads = {
  path: "/workloads",
  parent: true,
  sort: 2,
  component: Layout,
  name: "Workloads",
  meta: {
    title: "business.workload.workload",
    icon: "iconfont iconworkload",
  },
  children: [
    {
      path: "/pods",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "list",
      },
      component: () => import("@/business/workloads/pods"),
      name: "Pods",
      meta: {
        title: "Pods",
      },
    },
    {
      path: "/pods/create",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "create",
      },
      name: "PodCreate",
      hidden: true,
      component: () => import("@/business/workloads/pods/create"),
      props: true,
      meta: {
        activeMenu: "/pods",
      },
    },
    {
      path: "pods/:namespace/:name",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "update",
      },
      name: "PodEdit",
      hidden: true,
      component: () => import("@/business/workloads/pods/edit"),
      props: true,
      meta: {
        activeMenu: "/pods",
      },
    },
    {
      path: "/pods/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "get",
      },
      name: "PodDetail",
      hidden: true,
      component: () => import("@/business/workloads/pods/detail"),
      props: true,
      meta: {
        activeMenu: "/pods",
      },
    },

    {
      path: "/deployments",
      requirePermission: {
        apiGroup: "apps",
        resource: "deployments",
        verb: "list",
      },
      component: () => import("@/business/workloads/deployments"),
      name: "Deployments",
      meta: {
        title: "Deployments",
      },
    },
    {
      path: "/deployments/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "deployments",
        verb: "get",
      },
      name: "DeploymentDetail",
      hidden: true,
      component: () => import("@/business/workloads/deployments/detail"),
      props: true,
      meta: {
        activeMenu: "/deployments",
      },
    },
    {
      path: "deployments/:operation",
      requirePermission: {
        apiGroup: "apps",
        resource: "deployments",
        verb: "create",
      },
      name: "DeploymentCreate",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/deployments",
      },
    },
    {
      path: "deployments/:operation/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "deployments",
        verb: "update",
      },
      name: "DeploymentEdit",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/deployments",
      },
    },

    {
      path: "/daemonsets",
      requirePermission: {
        apiGroup: "apps",
        resource: "daemonsets",
        verb: "list",
      },
      component: () => import("@/business/workloads/daemonsets"),
      name: "DaemonSets",
      meta: {
        title: "DaemonSets",
      },
    },
    {
      path: "/daemonsets/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "daemonsets",
        verb: "get",
      },
      name: "DaemonSetDetail",
      hidden: true,
      component: () => import("@/business/workloads/daemonsets/detail"),
      props: true,
      meta: {
        activeMenu: "/daemonsets",
      },
    },
    {
      path: "daemonsets/:operation",
      requirePermission: {
        apiGroup: "apps",
        resource: "daemonsets",
        verb: "create",
      },
      name: "DaemonSetCreate",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/daemonsets",
      },
    },
    {
      path: "daemonsets/:operation/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "daemonsets",
        verb: "update",
      },
      name: "DaemonSetEdit",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/daemonsets",
      },
    },

    {
      path: "/statefulsets",
      requirePermission: {
        apiGroup: "apps",
        resource: "statefulsets",
        verb: "list",
      },
      component: () => import("@/business/workloads/statefulsets"),
      name: "StatefulSets",
      meta: {
        title: "StatefulSets",
      },
    },
    {
      path: "/statefulsets/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "statefulsets",
        verb: "get",
      },
      name: "StatefulSetDetail",
      hidden: true,
      component: () => import("@/business/workloads/statefulsets/detail"),
      props: true,
      meta: {
        activeMenu: "/statefulsets",
      },
    },
    {
      path: "statefulsets/:operation",
      requirePermission: {
        apiGroup: "apps",
        resource: "statefulsets",
        verb: "create",
      },
      name: "StatefulSetCreate",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/statefulsets",
      },
    },
    {
      path: "statefulsets/:operation/:namespace/:name",
      requirePermission: {
        apiGroup: "apps",
        resource: "statefulsets",
        verb: "update",
      },
      name: "StatefulSetEdit",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/statefulsets",
      },
    },

    {
      path: "/jobs",
      requirePermission: {
        apiGroup: "batch",
        resource: "jobs",
        verb: "list",
      },
      component: () => import("@/business/workloads/jobs"),
      name: "Jobs",
      meta: {
        title: "Jobs",
      },
    },
    {
      path: "/jobs/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "batch",
        resource: "jobs",
        verb: "get",
      },
      name: "JobDetail",
      hidden: true,
      component: () => import("@/business/workloads/jobs/detail"),
      props: true,
      meta: {
        activeMenu: "/jobs",
      },
    },
    {
      path: "jobs/:operation",
      requirePermission: {
        apiGroup: "batch",
        resource: "jobs",
        verb: "create",
      },
      name: "JobCreate",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/jobs",
      },
    },
    {
      path: "jobs/:operation/:namespace/:name",
      requirePermission: {
        apiGroup: "batch",
        resource: "jobs",
        verb: "update",
      },
      name: "JobEdit",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/jobs",
      },
    },

    {
      path: "/cronjobs",
      requirePermission: {
        apiGroup: "batch",
        resource: "cronjobs",
        verb: "list",
      },
      component: () => import("@/business/workloads/cronjobs"),
      name: "CronJobs",
      meta: {
        title: "CronJobs",
      },
    },
    {
      path: "/cronjobs/detail/:namespace/:name",
      requirePermission: {
        apiGroup: "batch",
        resource: "cronjobs",
        verb: "get",
      },
      name: "CronJobDetail",
      hidden: true,
      component: () => import("@/business/workloads/cronjobs/detail"),
      props: true,
      meta: {
        activeMenu: "/cronjobs",
      },
    },
    {
      path: "cronjobs/:operation",
      requirePermission: {
        apiGroup: "batch",
        resource: "cronjobs",
        verb: "create",
      },
      name: "CronJobCreate",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/cronjobs",
      },
    },
    {
      path: "cronjobs/:operation/:namespace/:name",
      requirePermission: {
        apiGroup: "batch",
        resource: "cronjobs",
        verb: "update",
      },
      name: "CronJobEdit",
      hidden: true,
      component: () => import("@/business/workloads/index"),
      meta: {
        activeMenu: "/cronjobs",
      },
    },
  ],
};

export default Workloads;
