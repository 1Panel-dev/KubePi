import Layout from "@/business/app-layout/horizontal-layout"

const Clusters = {
  path: "/cluster",
  sort: 1,
  component: Layout,
  name: "Cluster",
  meta: {
    title: "business.cluster.cluster",
    icon: "iconfont iconkubernets",
  },
  children: [
    {
      path: "/nodes",
      component: () => import("@/business/cluster/nodes"),
      name: "Nodes",
      meta: {
        title: "Nodes",
      }
    },
    {
      path: "/namespaces",
      component: () => import("@/business/cluster/namespaces"),
      name: "Namespaces",
      meta: {
        title: "Namespaces"
      }
    },
    {
      path: "/namespaces/create",
      component: () => import("@/business/cluster/namespaces/create"),
      name: "NamespaceCreate",
      hidden: true,
      meta: {
        activeMenu: "/namespaces"
      }
    },
    {
      path: "/namespaces/detail/:name",
      props: true,
      component: () => import("@/business/cluster/namespaces/detail"),
      name: "NamespaceDetail",
      hidden: true,
      meta: {
        activeMenu: "/namespaces"
      }
    },
    {
      path: "/events",
      component: () => import("@/business/cluster/events"),
      name: "events",
      meta: {
        title: "Events",
      }
    }
  ]
}

export default Clusters
