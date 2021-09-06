import {post, del, get, put} from "@/plugins/request"

const authUrl = "/api/v1/sessions"

export function login(data) {
    return post(authUrl, data)
}

export function logout() {
    return del(authUrl)
}

export function getCurrentUser() {
    return get(authUrl)
}

export function isLogin() {
    return get(`${authUrl}/status`)
}

export function updateProfile(data) {
    return put(authUrl, data)
}

export function updatePassword(data) {
    return put(`${authUrl}/password`, data)
}
