import router from "./router"
import NProgress from "nprogress"
import "nprogress/nprogress.css"
import store from "./store"

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ["/login"] // no redirect whitelist

const generateRoutes = async (to, from, next) => {

  try {
    const accessRoutes = await store.dispatch("permission/generateRoutes")
    router.addRoutes(accessRoutes)
    next({ ...to, replace: true })
  } catch (error) {
    NProgress.done()
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
    if (store.getters.permission_routes.length > 0 && to.name != null) {
      console.log(store.getters.permission_routes)
      next()
    } else {
      await generateRoutes(to, from, next)
    }
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
