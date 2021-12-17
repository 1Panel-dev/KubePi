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
    icon: "iconfont iconjingxiangcangku",
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
    },
    {
      path: "/imagerepos/:name/edit",
      component: () => import("@/business/imagerepo-management/operate"),
      name: "ImageRepoEdit",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/imagerepos",
      }
    },
    {
      path: "/imagerepos/:repo/detail",
      component: () => import("@/business/imagerepo-management/detail"),
      name: "ImageRepoDetail",
      hidden: true,
      props: true,
      meta: {
        activeMenu: "/imagerepos",
      }
    }
  ]
}


export default ImageRepos

