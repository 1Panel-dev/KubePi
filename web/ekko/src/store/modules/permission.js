import {rolesRoutes, constantRoutes} from "@/router"

const state = {
    routes: [],
    addRoutes: []
}

function hasPermission(user, route) {
    if (user && route) {
        return true
    }
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
    generateRoutes({commit}, p) {
        return new Promise(resolve => {
            const user = p
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
