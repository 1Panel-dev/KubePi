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
                apiGroup: "",
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
              apiGroup: "",
              resource: "persistentvolumes",
              verb: "list",
            },
            component: () => import("@/business/storage/pv/create"),
            name: "PersistentVolumeCreate",
            hidden: true,
            meta: {
              title: "Persistent Volume Create",
            }
        },
        {
            path: "/persistentvolumes/edit/:name",
            requirePermission: {
              apiGroup: "",
              resource: "persistentvolumes",
              verb: "list",
            },
            component: () => import("@/business/storage/pv/edit"),
            name: "PersistentVolumeEdit",
            props: true,
            hidden: true,
            meta: {
              title: "Persistent Volume Edit",
            }
        },
        {
            path: "/persistentvolumes/detail/:name",
            requirePermission: {
              apiGroup: "",
              resource: "persistentvolumes",
              verb: "list",
            },
            component: () => import("@/business/storage/pv/detail"),
            name: "PersistentVolumeDetail",
            props: true,
            hidden: true,
            meta: {
              title: "Persistent Volume Detail",
            }
        },
        {
            path: "/persistentvolumeclaims",
            requirePermission: {
                apiGroup: "",
                resource: "persistentvolumeclaims",
                verb: "list",
            },
            component: () => import("@/business/storage/pvc"),
            name: "Pvcs",
            meta: {
                title: "Persistent Volume Claims",
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
            }
        },
        {
            path: "/storageclasses/create",
            requirePermission: {
                apiGroup: "storage.k8s.io",
                resource: "storageclasses",
                verb: "list",
            },
            component: () => import("@/business/storage/sc/create"),
            name: "StorageClassCreate",
            hidden: true,
            meta: {
                title: "Storage Classes Create",
            }
        },
        {
            path: "/storageclasses/edit/:name",
            requirePermission: {
                apiGroup: "storage.k8s.io",
                resource: "storageclasses",
                verb: "list",
            },
            component: () => import("@/business/storage/sc/edit"),
            name: "StorageClassEdit",
            hidden: true,
            meta: {
                title: "Storage Classes Edit",
            }
        },
        {
            path: "/storageclasses/detail/:name",
            requirePermission: {
                apiGroup: "storage.k8s.io",
                resource: "storageclasses",
                verb: "list",
            },
            component: () => import("@/business/storage/sc/detail"),
            name: "StorageClassDetail",
            hidden: true,
            meta: {
                title: "Storage Classes Detail",
            }
        }
    ]
}

export default Storage
