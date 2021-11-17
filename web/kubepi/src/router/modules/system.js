import Layout from "@/business/app-layout/horizontal-layout";

const SystemLog = {
  path: "/systemlogs",
  sort: 3,
  component: Layout,
  name: "Systems",
  children: [
    {
      path: "/systemlogs",
      component: () => import("@/business/system-log"),
      name: "SystemLog",
      meta: {
        title: "business.system.system_log",
        icon: "iconfont iconxitongrizhi",
      },
    },
  ],
}

export default SystemLog;
