import {rolesRoutes, constantRoutes} from "@/router"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(clusterRoles, route) {
    if (route.requirePermission) {
        for (const clusterRole of clusterRoles) {
            if (clusterRole.rules.length > 0) {
                for (const rule of clusterRole.rules) {
                    if (rule.apiGroups.includes("*") || rule.apiGroups.includes(route.requirePermission.apiGroup)) {
                        if (rule.resources.includes("*") || rule.resources.includes(route.requirePermission.resource)) {
                            if (rule.verbs.includes("*") || rule.verbs.includes(route.requirePermission.verb)) {
                                return true
                            }
                        }
                    }
                }
            }
        }
        return false
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
