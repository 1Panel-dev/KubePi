import Layout from "@/business/app-layout/horizontal-layout"

const Storage = {
  path: "/storage",
  sort: 3,
  component: Layout,
  name: "Storage",
  meta: {
    title: "business.storage.storage",
    icon: "iconfont iconstorage"
  },
  children: [
    {
      path: "/pv",
      component: () => import("@/business/storage/pvs"),
      name: "Pvs",
      meta: {
        title: "Persistent Volume",
      }
    },
    {
      path: "/pvcs",
      component: () => import("@/business/storage/pvcs"),
      name: "Pvcs",
      meta: {
        title: "Persistent Volume Claims",
      }
    },
    {
      path: "/storageclasses",
      component: () => import("@/business/storage/storage-classes"),
      name: "StorageClasses",
      meta: {
        title: "Storage Classes",
      }
    }
  ]
}

export default Storage
