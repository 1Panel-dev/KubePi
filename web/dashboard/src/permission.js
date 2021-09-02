import router from "./router"
import NProgress from "nprogress"
import "nprogress/nprogress.css"
import store from "./store"
import Layout from "@/business/app-layout/horizontal-layout"

NProgress.configure({showSpinner: false}) // NProgress Configuration

const generateRoutes = async (to, from, next) => {
    const hasRoles = store.getters.isAdmin || (store.getters.clusterRoles && store.getters.clusterRoles.length > 0)
    if (hasRoles) {
        next()
    } else {
        try {
            const user = await store.dispatch("user/getCurrentUser")
            const accessRoutes = await store.dispatch("permission/generateRoutes", user)
            if (accessRoutes.length > 0) {
                const root = {
                    path: "/",
                    component: Layout,
                    redirect: accessRoutes[0].path,
                }
                router.addRoute(root)
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
    if (isLogin) {
        if (!to.query["cluster"]) {
            if (from.query["cluster"]) {
                await store.dispatch("user/setCurrentCluster", from.query["cluster"])
                const q = to.query
                q["cluster"] = from.query["cluster"]
                next({path: to.path, query: q})
                NProgress.done()
            } else {
                window.open("/ekko", '_self');
            }
        } else {
            await store.dispatch("user/setCurrentCluster", to.query["cluster"])
        }
        await generateRoutes(to, from, next)
    } else {
        window.open("/ekko", '_self');
        NProgress.done()
    }
})

router.afterEach(() => {
    NProgress.done()
})
