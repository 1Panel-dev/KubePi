import {post, del, get, put} from "@/plugins/request"

const baseUrl = "/api/v1/users"

export function searchUsers(pageNum, pageSize, conditions) {
    return post(`${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`, conditions)
}

export function createUser(user) {
    return post(baseUrl, user)
}

export function deleteUser(name) {
    return del(`${baseUrl}/${name}`)
}

export function getUser(name) {
    return get(`${baseUrl}/${name}`)
}

export function updateUser(name, user) {

    return put(`${baseUrl}/${name}`, user)
}