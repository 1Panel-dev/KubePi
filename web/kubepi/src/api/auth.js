import {post, del, get, put} from "@/plugins/request"

const authUrl = "/api/v1/sessions"

export function login(data, config) {
    return post(authUrl, data, undefined, config)
}

export function logout() {
    return del(authUrl)
}

export function getCurrentUser() {
    return get(authUrl)
}

export function isLogin(config) {
    return get(`${authUrl}/status`, undefined, undefined, config)
}

export function updateProfile(data) {
    return put(authUrl, data)
}

export function updatePassword(data, config) {
    return put(`${authUrl}/password`, data, undefined, config)
}
