import {rolesRoutes, constantRoutes} from "@/router"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(permissions, route) {
    if (route.requirePermission) {
        for (const resource of Object.keys(permissions)) {
            if (resource === "*") {
                return true
            }
            if (route.requirePermission && route.requirePermission.resource === resource) {
                for (const verb of permissions[resource]) {
                    if (verb === "*") {
                        return true
                    }
                    if (verb === route.requirePermission.verb) {
                        return true
                    }
                }
            }
        }
        return false
    }
    return true
}


export function filterRolesRoutes(routes, permissions) {
    const res = []
    routes.forEach(route => {
        const tmp = {...route}
        if (hasPermission(permissions, tmp)) {
            if (tmp.children) {
                tmp.children = filterRolesRoutes(tmp.children, permissions)
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
        return new Promise(resolve => {
            const permissions = p
            let accessedRoutes
            accessedRoutes = filterRolesRoutes(rolesRoutes, permissions)
            commit("SET_ROUTES", accessedRoutes)
            resolve(accessedRoutes)
        })
    }
}

const menuActions = {
    SET_MENU: (state, menu) => {
        state.menu = menu
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions,
    menuActions
}
