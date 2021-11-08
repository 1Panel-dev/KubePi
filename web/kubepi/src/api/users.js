import {post, del, get, put} from "@/plugins/request"

const baseUrl = "/api/v1/users"

export function searchUsers(pageNum, pageSize, conditions) {
    let url = `${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    return post(url, {conditions: conditions})
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

export function listUsers() {
    return get(`${baseUrl}`)
}