import Layout from "@/business/app-layout/horizontal-layout"

const Storage = {
  path: "/storage",
  sort: 3,
  parent: true,
  component: Layout,
  name: "Storage",
  meta: {
    title: "business.storage.storage",
    icon: "iconfont iconstorage"
  },
  children: [
    {
      path: "/persistentvolumes",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumes",
        verb: "list",
      },
      component: () => import("@/business/storage/pv"),
      name: "PersistentVolumes",
      meta: {
        title: "Persistent Volume",
      }
    },
    {
      path: "/persistentvolumes/create",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumes",
        verb: "create",
      },
      component: () => import("@/business/storage/pv/create"),
      name: "PersistentVolumeCreate",
      hidden: true,
      meta: {
        title: "Persistent Volume Create",
        activeMenu: "/persistentvolumes",
      }
    },
    {
      path: "/persistentvolumes/edit/:name",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumes",
        verb: "update",
      },
      component: () => import("@/business/storage/pv/edit"),
      name: "PersistentVolumeEdit",
      props: true,
      hidden: true,
      meta: {
        title: "Persistent Volume Edit",
        activeMenu: "/persistentvolumes"
      }
    },
    {
      path: "/persistentvolumes/detail/:name",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumes",
        verb: "get",
      },
      component: () => import("@/business/storage/pv/detail"),
      name: "PersistentVolumeDetail",
      props: true,
      hidden: true,
      meta: {
        title: "Persistent Volume Detail",
        activeMenu: "/persistentvolumes"
      }
    },
    {
      path: "/storageclasses",
      requirePermission: {
        apiGroup: "storage.k8s.io",
        resource: "storageclasses",
        verb: "list",
      },
      component: () => import("@/business/storage/sc"),
      name: "StorageClasses",
      meta: {
        title: "Storage Classes",
        activeMenu: "/storageclasses",
      }
    },
    {
      path: "/storageclasses/create",
      requirePermission: {
        apiGroup: "storage.k8s.io",
        resource: "storageclasses",
        verb: "crate",
      },
      component: () => import("@/business/storage/sc/create"),
      name: "StorageClassCreate",
      hidden: true,
      meta: {
        title: "Storage Classes Create",
        activeMenu: "/storageclasses",
      }
    },
    {
      path: "/storageclasses/edit/:name",
      requirePermission: {
        apiGroup: "storage.k8s.io",
        resource: "storageclasses",
        verb: "update",
      },
      component: () => import("@/business/storage/sc/edit"),
      name: "StorageClassEdit",
      hidden: true,
      meta: {
        title: "Storage Classes Edit",
        activeMenu: "/storageclasses",
      }
    },
    {
      path: "/storageclasses/detail/:name",
      requirePermission: {
        apiGroup: "storage.k8s.io",
        resource: "storageclasses",
        verb: "get",
      },
      component: () => import("@/business/storage/sc/detail"),
      name: "StorageClassDetail",
      hidden: true,
      meta: {
        title: "Storage Classes Detail",
        activeMenu: "/storageclasses",
      }
    }, {
      path: "/persistentVolumeClaim",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentVolumeClaim",
        verb: "list",
      },
      component: () => import("@/business/storage/pvc/"),
      name: "PersistentVolumeClaim",
      meta: {
        title: "PersistentVolumeClaim",
        activeMenu: "/persistentVolumeClaim",
      }
    }, {
      path: "/persistentVolumeClaim/create",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumeclaims",
        verb: "create",
      },
      component: () => import("@/business/storage/pvc/create"),
      name: "PersistentVolumeClaimCreate",
      hidden: true,
      meta: {
        title: "PersistentVolumeClaim Create",
        activeMenu: "/persistentVolumeClaim",
      }
    }, {
      path: "/persistentVolumeClaim/edit/:name",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumeclaims",
        verb: "update",
      },
      component: () => import("@/business/storage/pvc/edit"),
      name: "PersistentVolumeClaimEdit",
      hidden: true,
      meta: {
        title: "PersistentVolumeClaim Edit",
        activeMenu: "/persistentVolumeClaim",
      }
    }, {
      path: "/persistentVolumeClaim/detail/:name",
      requirePermission: {
        apiGroup: "core",
        resource: "persistentvolumeclaims",
        verb: "get",
      },
      component: () => import("@/business/storage/pvc/detail"),
      name: "PersistentVolumeClaimDetail",
      hidden: true,
      meta: {
        title: "PersistentVolumeClaim Detail",
        activeMenu: "/persistentVolumeClaim",
      }
    },
  ]
}

export default Storage
