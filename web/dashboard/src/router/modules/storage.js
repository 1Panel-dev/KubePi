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
            component: () => import("@/business/storage/pvs"),
            name: "Pvs",
            meta: {
                title: "Persistent Volume",
            }
        },
        {
            path: "/persistentvolumeclaims",
            requirePermission: {
                apiGroup: "",
                resource: "persistentvolumeclaims",
                verb: "list",
            },
            component: () => import("@/business/storage/pvcs"),
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
            component: () => import("@/business/storage/storage-classes"),
            name: "StorageClasses",
            meta: {
                title: "Storage Classes",
            }
        }
    ]
}

export default Storage
