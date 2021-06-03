/* 前后端不分离的登录方式*/
// import {updateUser} from "@/api/user"
// import {login, logout, isLogin, getSession} from "@/api/auth"
// import {resetRouter} from "@/router"
import {getLanguage, setLanguage} from "@/i18n"

const state = {
  login: false,
  name: "",
  currentProject: "",
  language: getLanguage(),
  roles: []
}

const mutations = {
  LOGIN: (state) => {
    state.login = true
  },
  LOGOUT: (state) => {
    state.login = false
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_LANGUAGE: (state, language) => {
    state.language = language
    setLanguage(language)
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  },
  SET_CURRENT_MENU: (state, menu) => {
    state.menu = menu
  },
}

const actions = {
  login ({ commit }, userInfo) {
    console.log(userInfo)
    // const { username, password, captchaId, code } = userInfo
    return new Promise((resolve) => {
      commit("LOGIN")
      const user = {
        name: "admin",
        roles: ['ADMIN'],
        menu: "global"
      }
      resolve(user)
      // login({ username: username.trim(), password: password, captchaId: captchaId, code: code }).then(response => {
      //   commit("LOGIN")
      //   resolve(response)
      // }).catch(error => {
      //   reject(error)
      // })
    })
  },

  isLogin () {
    return new Promise((resolve) => {
      if (state.login) {
        resolve(true)
      }else {
        resolve(false)
      }
      // isLogin().then((data) => {
      //   if (data.isLogin) {
      //     commit("LOGIN")
      //     resolve(true)
      //   } else {
      //     resolve(false)
      //   }
      // }).catch(() => {
      //   resolve(false)
      // })
    })
  },

  getCurrentUser ({ commit }) {
    return new Promise((resolve) => {
      const user = {
        name: "admin",
        roles: ['ADMIN'],
        menu: "cluster1"
      }
      commit("SET_CURRENT_MENU", user.menu)
      commit("SET_ROLES", user.roles)
      resolve(user)
      // getSession().then(data => {
      //   const user = data.user
      //   const { name, roles, language, currentProject } = user
      //   commit("SET_NAME", name)
      //   commit("SET_ROLES", roles)
      //   commit("SET_LANGUAGE", language)
      //   commit("SET_CURRENT_PROJECT", currentProject)
      //   resolve(user)
      // }).catch(error => {
      //   reject(error)
      // })
    })
  },
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
