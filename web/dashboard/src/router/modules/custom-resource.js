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
      component: () => import("@/business/custom-resource"),
      name: "CustomResources",
      meta: {
        title: "business.custom_resource.custom_resource",
        icon: "iconfont iconcustom",
        global: false
      },
    },
  ]
}

export default CustomResource



