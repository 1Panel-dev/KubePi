import Layout from "@/business/app-layout/horizontal-layout";

const Workloads = {
  path: "/workloads",
  sort: 4,
  component: Layout,
  name: "Workloads",
  meta: {
    title: "business.workload.workload",
    icon: "iconfont iconworkload",
  },
  children: [
    {
      path: "/pods",
      component: () => import("@/business/workloads/pods"),
      name: "Pods",
      meta: {
        title: "Pods",
      },
    },
    {
      path: "/deployments",
      component: () => import("@/business/workloads/deployments"),
      name: "Deployments",
      meta: {
        title: "Deployments",
      },
    },
    {
      path: "/deployments/detail/:namespace/:name",
      name: "DeploymentDetail",
      hidden: true,
      component: () => import("@/business/workloads/deployments/detail"),
      meta: {
        activeMenu: "/deployments",
      },
    },
    {
      path: "/deployments/create",
      name: "DeploymentCreate",
      hidden: true,
      component: () => import("@/business/workloads/deployments/create"),
      meta: {
        activeMenu: "/deployments",
      },
    },
    {
      path: "/deployments/edit/:namespace/:name",
      name: "DeploymentEdit",
      hidden: true,
      component: () => import("@/business/workloads/deployments/edit"),
      meta: {
        activeMenu: "/deployments",
      },
    },

    {
      path: "/cronjobs",
      component: () => import("@/business/workloads/cronjobs"),
      name: "CronJobs",
      meta: {
        title: "CronJobs",
      },
    },
    {
      path: "/cronjobs/detail/:namespace/:name",
      name: "CronJobDetail",
      hidden: true,
      component: () => import("@/business/workloads/cronjobs/detail"),
      meta: {
        activeMenu: "/cronjobs",
      },
    },
    {
      path: "/cronjobs/create",
      name: "CronJobCreate",
      hidden: true,
      component: () => import("@/business/workloads/cronjobs/create"),
      meta: {
        activeMenu: "/cronjobs",
      },
    },
    {
      path: "/cronjobs/edit/:namespace/:name",
      name: "CronJobEdit",
      hidden: true,
      component: () => import("@/business/workloads/cronjobs/edit"),
      meta: {
        activeMenu: "/cronjobs",
      },
    },

    {
      path: "/jobs",
      component: () => import("@/business/workloads/jobs"),
      name: "Jobs",
      meta: {
        title: "Jobs",
      },
    },
    {
      path: "/jobs/detail/:namespace/:name",
      name: "JobDetail",
      hidden: true,
      component: () => import("@/business/workloads/jobs/detail"),
      meta: {
        activeMenu: "/jobs",
      },
    },
    {
      path: "/jobs/create",
      name: "JobCreate",
      hidden: true,
      component: () => import("@/business/workloads/jobs/create"),
      meta: {
        activeMenu: "/jobs",
      },
    },
    {
      path: "/jobs/edit/:namespace/:name",
      name: "JobEdit",
      hidden: true,
      component: () => import("@/business/workloads/jobs/edit"),
      meta: {
        activeMenu: "/jobs",
      },
    },
  ],
};

export default Workloads;
