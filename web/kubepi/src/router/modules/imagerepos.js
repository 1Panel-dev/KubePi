import Layout from "@/business/app-layout/horizontal-layout"

const ImageRepos = {
  path: "/imagerepos",
  sort: 3,
  component: Layout,
  name: "ImageRepoManagement",
  requirePermission: {
    resource: "imagerepos",
    verb: "list"
  },
  meta: {
    title: "business.image_repos.list",
    icon: "el-icon-box",
  },
  children: [
    {
      path: "/imagerepos",
      component: () => import("@/business/imagerepo-management"),
      name: "ImageRepos",
      meta: {
        title: "business.image_repos.list",
        activeMenu: "/imagerepos",
      }
    },
    {
      path: "/imagerepos/create",
      component: () => import("@/business/imagerepo-management/operate"),
      name: "ImageRepoCreate",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/imagerepos",
      }
    }
  ]
}


export default ImageRepos

