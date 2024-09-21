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
            for (const i in accessRoutes) {
              let father = accessRoutes[i]
              if (father.children?.length >0) {
                  for (const j in father.children) {
                    const route = father.children[j]
                    if (route.path.indexOf("create") > -1 || route.path.indexOf("operation") >-1) {
                      father.children.push(           {
                        path: route.path+"/yaml",
                        hidden: true,
                        name: route.name+"Yaml",
                        props: true,
                        component: () => import("@/business/yaml"),
                        meta: {
                          activeMenu: route.meta?.activeMenu
                        }
                      })
                    }
                  }
              }
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
                // if (to.path.indexOf("yaml") > -1) {
                //   router.addRoute(

              //   )
                // }
                next({path: to.path, query: q})
                NProgress.done()
                document.title = q["cluster"] || "KubePi"
            } else {
                window.open("/kubepi", '_self');
            }
        } else {
            await store.dispatch("user/setCurrentCluster", to.query["cluster"])
            document.title = to.query["cluster"] || "KubePi"
        }
        await generateRoutes(to, from, next)
    } else {
        window.open("/kubepi", '_self');
        NProgress.done()
    }
})

router.afterEach(() => {
    NProgress.done()
})
