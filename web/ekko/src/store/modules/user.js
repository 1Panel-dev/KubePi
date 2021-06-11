import {login, isLogin, logout, getCurrentUser} from "@/api/auth"
import {resetRouter} from "@/router"
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
    SET_NICK_NAME: (state, nickName) => {
        state.nickName = nickName
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
    login({commit}, userInfo) {
        const {username, password} = userInfo
        return new Promise((resolve, reject) => {
            commit("LOGIN")
            login({username: username.trim(), password: password}).then(response => {
                commit("LOGIN")
                resolve(response)
            }).catch(error => {
                reject(error)
            })
        })
    },

    isLogin({commit}) {
        return new Promise((resolve) => {
            if (state.login) {
                resolve(true)
            }
            isLogin().then((data) => {
                if (data.data) {
                    commit("LOGIN")
                    resolve(true)
                } else {
                    resolve(false)
                }
            }).catch(() => {
                resolve(false)
            })
        })
    },
    getCurrentUser({commit}) {
        return new Promise((resolve, reject) => {
            getCurrentUser().then(data => {

                const user = data.data
                user["roles"] = ["ADMIN"]
                user["menu"] = ["cluster1"]
                const {name, roles, menu, nickName} = user
                commit("SET_NAME", name)
                commit("SET_ROLES", roles)
                commit("SET_LANGUAGE", "zh-CN")
                commit("SET_CURRENT_MENU", menu)
                commit("SET_NICK_NAME", nickName)
                resolve(user)
            }).catch(error => {
                reject(error)
            })
        })
    },
    logout({commit}) {
        logout().then(() => {
            commit('LOGOUT')
            commit('SET_ROLES', [])
            resetRouter()
        })
    },
}
export default {
    namespaced: true,
    state,
    mutations,
    actions
}
