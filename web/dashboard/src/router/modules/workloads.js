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
    scope:"namespace"
  },
  children: [
    {
      path: "/pods",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "list",
        scope:"namespace"
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
        scope:"namespace"
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
      path: "/pods/top",
      requirePermission: {
        apiGroup: "",
        resource: "pods",
        verb: "list",
        scope:"namespace"
      },
      name: "PodTop",
      hidden: true,
      component: () => import("@/business/workloads/pods/top"),
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
        scope:"namespace"
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
