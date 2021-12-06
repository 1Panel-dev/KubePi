import Layout from "@/business/app-layout/horizontal-layout";

const SystemLog = {
  path: "/systemlogs",
  sort: 4,
  component: Layout,
  name: "Systems",
  requirePermission: {
    resource: "systems",
    verb: "list"
  },
  meta: {
    title: "business.system.system_log",
    icon: "iconfont iconxitongrizhi",
  },
  children: [
    {
      path: "/loginlogs",
      component: () => import("@/business/system-log/login/index"),
      name: "LoginLog",
      meta: {
        title: "business.system.login_log",
        activeMenu: "/loginlogs",
      },
    },
    {
      path: "/operationlogs",
      component: () => import("@/business/system-log/operation/index"),
      name: "OperationLog",
      meta: {
        title: "business.system.operation_log",
        activeMenu: "/operationlogs",
      },
    },
  ],
}

export default SystemLog;
