import router from "./router"
import NProgress from "nprogress"
import "nprogress/nprogress.css"
import store from "./store"
import Layout from "@/business/app-layout/horizontal-layout"

NProgress.configure({showSpinner: false}) // NProgress Configuration

const whiteList = ["/login"] // no redirect whitelist

const generateRoutes = async (to, from, next) => {
    const hasRoles = store.getters.roles && store.getters.roles.length > 0
    if (hasRoles) {
        next()
    } else {
        try {
            const user = await store.dispatch("user/getCurrentUser")
            user.menu = 'global'
            const accessRoutes = await store.dispatch("permission/generateRoutes", user)
            if (user.menu === 'global') {
                router.addRoutes([{
                    path: "/",
                    component: Layout,
                    redirect: "/clusterlist",
                },])
            } else {
                router.addRoutes([{
                    path: "/",
                    component: Layout,
                    redirect: "/dashboard",
                },])
            }
            router.addRoutes(accessRoutes)
            next({...to, replace: true})
        } catch (error) {
            NProgress.done()
        }
    }
}

//路由前置钩子，根据实际需求修改
router.beforeEach(async (to, from, next) => {
    NProgress.start()
    const isLogin = await store.dispatch("user/isLogin")
    console.log(isLogin)
    if (isLogin) {
        if (to.path === "/login") {
            next({path: "/"})
            NProgress.done()
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
