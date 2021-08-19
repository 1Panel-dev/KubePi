import {login, isLogin, logout, getCurrentClusterUser} from "@/api/auth"
import {resetRouter} from "@/router"
import {getLanguage, setLanguage} from "@/i18n"
import store from "../../store"
import {updateProfile} from "../../../../ekko/src/api/auth";

const state = {
    login: false,
    name: "",
    language: getLanguage(),
    roles: [],
    cluster: "",
    clusterRoles: [],
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
    SET_CURRENT_CLUSTER: (state, name) => {
        state.cluster = name
    },
    SET_CLUSTER_ROLES: (state, clusterRoles) => {
        state.clusterRoles = clusterRoles
    }
}

const actions = {
    setCurrentCluster({commit}, clusterName) {
        commit("SET_CURRENT_CLUSTER", clusterName)
    },
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
    setLanguage({commit}, lang) {
        return new Promise((resolve, reject) => {
            updateProfile({"language": lang}).then(() => {
                commit("SET_LANGUAGE", lang)
                resolve()
            }).catch(error => {
                reject(error)
            })
        })
    },
    getCurrentUser({commit}) {
        return new Promise((resolve, reject) => {
            const clusterName = store.getters.cluster
            getCurrentClusterUser(clusterName).then(data => {
                const user = data.data
                user["roles"] = ["ADMIN"]
                const {name, nickName, roles,language, clusterRoles} = user
                commit("SET_NAME", name)
                commit("SET_ROLES", roles)
                commit("SET_CLUSTER_ROLES", clusterRoles)
                commit("SET_NICK_NAME", nickName)
                commit("SET_LANGUAGE", language)

                resolve(user)
            }).catch(error => {
                reject(error)
            })
        })
    },
    logout({commit}) {
        logout().then(() => {
            commit("LOGOUT")
            commit("SET_ROLES", [])
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
