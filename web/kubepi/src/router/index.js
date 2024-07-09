import Vue from "vue"
import Router from "vue-router"
import Layout from "@/business/app-layout/horizontal-layout"
import Error403 from "@/business/error-page/403"

// 加载modules中的路由
const modules = require.context("./modules", true, /\.js$/)

// 修复路由变更后报错的问题
const routerPush = Router.prototype.push
Router.prototype.push = function push(location) {
    return routerPush.call(this, location).catch(error => error)
}

Vue.use(Router)
export const constantRoutes = [
    {
        path: "/redirect",
        component: Layout,
        hidden: true,
        name: "redirect",
        children: [
            {
                path: "/redirect/:path(.*)",
                component: () => import("@/components/redirect")
            }
        ]
    },
    {
        path: "/profile",
        component: Layout,
        hidden: true,
        children: [
            {
                path: "",
                component: () => import("@/business/profile")
            }
        ]
    },
    {
        path: "/login",
        component: () => import("@/business/login"),
        hidden: true
    },
    {
        path: "/sso",
        component: () => import("@/business/sso"),
        hidden: true
    },
    {
        path: "/403",
        component: Error403,
        hidden: true
    }
]

/**
 * 用户登录后根据角色加载的路由
 */
export const rolesRoutes = [
    // 先按sort排序
    ...modules.keys().map(key => modules(key).default).sort((r1, r2) => {
        r1.sort ??= Number.MAX_VALUE
        r2.sort ??= Number.MAX_VALUE
        return r1.sort - r2.sort
    }),
    {path: "*", redirect: "/", hidden: true}
]

const createRouter = () => new Router({
    scrollBehavior: () => ({y: 0}),
    routes: constantRoutes,
    mode: 'history',
    base: "kubepi"
})

const router = createRouter()


export function resetRouter() {
    const newRouter = createRouter()
    router.matcher = newRouter.matcher // reset router
}

export default router

