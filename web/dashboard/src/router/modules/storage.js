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
      path: "/persistent-volumes",
      component: () => import("@/business/storage/pv"),
      name: "PersistentVolumes",
      meta: {
        title: "Persistent Volumes",
      }
    },
    {
      path: "/persistent-volumes/create",
      component: () => import("@/business/storage/pv/create"),
      name: "PersistentVolumeCreate",
      hidden: true,
      meta: {
        title: "Persistent Volume Create",
      }
    },
    {
      path: "/persistent-volumes/edit/:name",
      component: () => import("@/business/storage/pv/edit"),
      name: "PersistentVolumeEdit",
      props: true,
      hidden: true,
      meta: {
        title: "Persistent Volume Edit",
      }
    },
    {
      path: "/persistent-volumes/detail/:name",
      component: () => import("@/business/storage/pv/detail"),
      name: "PersistentVolumeDetail",
      props: true,
      hidden: true,
      meta: {
        title: "Persistent Volume Detail",
      }
    },
    {
      path: "/pvc",
      component: () => import("@/business/storage/pvc"),
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
