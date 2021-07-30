import router from "./router"
import NProgress from "nprogress"
import "nprogress/nprogress.css"
import store from "./store"
import Layout from "@/business/app-layout/horizontal-layout"

NProgress.configure({ showSpinner: false }) // NProgress Configuration
const whiteList = ["/login"] // no redirect whitelist

const generateRoutes = async (to, from, next) => {
  const hasRoles = store.getters.clusterRoles && store.getters.clusterRoles.length > 0
  if (hasRoles) {
    next()
  } else {
    try {
      const { clusterRoles } = await store.dispatch("user/getCurrentUser")
      const accessRoutes = await store.dispatch("permission/generateRoutes", clusterRoles)
      if (accessRoutes.length > 0) {
        const root = {
          path: "/",
          component: Layout,
          redirect: accessRoutes[0].path,
        }
        router.addRoute(root)
      }
      router.addRoutes(accessRoutes)
      next({ ...to, replace: true })
    } catch (error) {
      NProgress.done()
    }
  }
}

//路由前置钩子，根据实际需求修改
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  const isLogin = await store.dispatch("user/isLogin")
  if (isLogin) {
    if (to.path === "/login") {
      next({ path: "/" })
      NProgress.done()
    }
    if (!to.query["cluster"]) {
      if (from.query["cluster"]) {
        await store.dispatch("user/setCurrentCluster", from.query["cluster"])
        const q = to.query
        q["cluster"] = from.query["cluster"]
        next({ path: to.path, query: q })
        NProgress.done()
      } else {
        console.log("no cluster")
      }
    } else {
      await store.dispatch("user/setCurrentCluster", to.query["cluster"])
    }
    await generateRoutes(to, from, next)
  } else {
    /* has not login*/
    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      next()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  NProgress.done()
})
