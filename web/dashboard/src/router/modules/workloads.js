import Layout from "@/business/app-layout/horizontal-layout";

const Workloads = {
    path: "/workloads",
    parent: true,
    sort: 4,
    component: Layout,
    name: "Workloads",
    meta: {
        title: "business.workload.workload",
        icon: "iconfont iconworkload",
    },
    children: [
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
            name: "DaemonSetDetail",
            hidden: true,
            component: () => import("@/business/workloads/daemonsets/detail"),
            meta: {
                activeMenu: "/daemonsets",
            },
        },
        {
            path: "/daemonsets/create",
            name: "DaemonSetCreate",
            hidden: true,
            component: () => import("@/business/workloads/daemonsets/create"),
            meta: {
                activeMenu: "/daemonsets",
            },
        },
        {
            path: "/daemonsets/edit/:namespace/:name",
            name: "DaemonSetEdit",
            hidden: true,
            component: () => import("@/business/workloads/daemonsets/edit"),
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
            name: "StatefulSetDetail",
            hidden: true,
            component: () => import("@/business/workloads/statefulsets/detail"),
            meta: {
                activeMenu: "/statefulsets",
            },
        },
        {
            path: "/statefulsets/create",
            name: "StatefulSetCreate",
            hidden: true,
            component: () => import("@/business/workloads/statefulsets/create"),
            meta: {
                activeMenu: "/statefulsets",
            },
        },
        {
            path: "/statefulsets/edit/:namespace/:name",
            name: "StatefulSetEdit",
            hidden: true,
            component: () => import("@/business/workloads/statefulsets/edit"),
            meta: {
                activeMenu: "/statefulsets",
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
            path: "/pods/detail/:namespace/:name",
            name: "PodDetail",
            hidden: true,
            component: () => import("@/business/workloads/pods/detail"),
            meta: {
                activeMenu: "/pods",
            },
        },
    ],
};

export default Workloads;
