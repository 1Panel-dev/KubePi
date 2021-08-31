import {rolesRoutes, constantRoutes} from "@/router"
import {checkPermissions} from "@/utils/permission"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(clusterRoles, route) {
    if (route.requirePermission) {
        return checkPermissions(route.requirePermission)
    }
    return true
}


export function filterRolesRoutes(routes, clusterRoles) {
    const res = []
    routes.forEach(route => {
        const tmp = {...route}
        if (hasPermission(clusterRoles, tmp)) {
            if (tmp.children) {
                tmp.children = filterRolesRoutes(tmp.children, clusterRoles)
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
    generateRoutes({commit}, p) {
        return new Promise((resolve, reject) => {
            const clusterRoles = p
            let accessedRoutes
            try {
                accessedRoutes = filterRolesRoutes(rolesRoutes, clusterRoles)
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
