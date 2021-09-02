import {rolesRoutes, constantRoutes} from "@/router"
import {checkPermissions} from "@/utils/permission"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(user, route) {
    if (route.requirePermission) {
        return checkPermissions(route.requirePermission)
    }
    return true
}


export function filterRolesRoutes(routes, user) {
    const res = []
    routes.forEach(route => {
        const tmp = {...route}
        if (hasPermission(user, tmp)) {
            if (tmp.children) {
                tmp.children = filterRolesRoutes(tmp.children, user)
            }
            res.push(tmp)
        }
    })

    return res
}

const mutations = {
    SET_ROUTES: (state, routes) => {
        state.addRoutes = routes
        state.routes = constantRoutes.concat(routes)
    }
}

const actions = {
    generateRoutes({commit}, user) {
        return new Promise((resolve, reject) => {
            let accessedRoutes
            try {
                accessedRoutes = filterRolesRoutes(rolesRoutes, user)
                for (const route of accessedRoutes) {
                    if (route.parent) {
                        let hidden = true
                        if (route.children.length > 0) {
                            for (const childRoute of route.children) {
                                hidden = hidden && childRoute.hidden
                            }
                            if (hidden) {
                                route.hidden = true
                            }
                        } else {
                            route.hidden = true
                        }
                    }
                }

                commit("SET_ROUTES", accessedRoutes)
                resolve(accessedRoutes)
            } catch (error) {
                reject(error)
            }
        })
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions,
}
