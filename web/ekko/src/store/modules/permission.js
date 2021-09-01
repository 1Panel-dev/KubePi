import {rolesRoutes, constantRoutes} from "@/router"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(user, route) {
    if (user.isAdministrator) {
        return true
    }
    if (route.requirePermission) {
        for (const resource of Object.keys(user.resourcePermissions)) {
            if (resource === "*") {
                return true
            }
            if (route.requirePermission && route.requirePermission.resource === resource) {
                for (const verb of user.resourcePermissions[resource]) {
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
        return new Promise(resolve => {
            let accessedRoutes
            accessedRoutes = filterRolesRoutes(rolesRoutes, user)
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
