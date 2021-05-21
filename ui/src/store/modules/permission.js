import {rolesRoutes, constantRoutes} from "@/router"

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  generateRoutes ({ commit }) {
    return new Promise(resolve => {
      let accessedRoutes
      accessedRoutes = rolesRoutes
      commit("SET_ROUTES", accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
