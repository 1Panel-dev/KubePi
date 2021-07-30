import Layout from "@/business/app-layout/horizontal-layout"

const CustomResource = {
  path: "/customResources",
  sort: 6,
  global: true,
  component: Layout,
  name: "CustomResources",
  children: [
    {
      path: "/customResources",
      requirePermission: {
        apiGroup: "apiextensions.k8s.io",
        resource: "customresourcedefinitions",
        verb: "list",
      },
      component: () => import("@/business/custom-resource"),
      name: "CustomResourceDefinitions",
      meta: {
        title: "business.custom_resource.custom_resource",
        icon: "iconfont iconcustom",
        global: false
      },
    },
    {
      path: "/customResources/:name/detail",
      component: () => import("@/business/custom-resource/detail"),
      hidden: true,
      props: true,
      name: "CustomResourceDefinitionDetail",
      meta: {
        activeMenu: "/customResources",
        global: false
      },
    },
    {
      path: "/customResources/create",
      component: () => import("@/business/custom-resource/create"),
      hidden: true,
      name: "CustomResourceDefinitionCreate",
      meta: {
        activeMenu: "/customResources",
        global: false
      },
    },
    {
      path: "/customResources/:name/edit",
      component: () => import("@/business/custom-resource/edit"),
      hidden: true,
      props: true,
      name: "CustomResourceDefinitionEdit",
      meta: {
        activeMenu: "/customResources",
        global: false
      },
    },
  ]
}

export default CustomResource



