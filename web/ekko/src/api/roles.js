import {get, post, del, put} from "@/plugins/request"

const baseUrl = "/api/v1/roles"

export function searchRoles(pageNum, pageSize, pattern) {
    console.log(pattern)
    let url = `${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    if (pattern) {
        url = `${url}&pattern=${pattern}`
    }
    return get(url)
}

export function listRoles() {
    return get(baseUrl)
}

export function createRole(role) {
    return post(baseUrl, role)
}

export function deleteRole(name) {
    return del(`${baseUrl}/${name}`)
}

export function getRole(name) {
    return get(`${baseUrl}/${name}`)
}

export function updateRole(name, role) {
    return put(`${baseUrl}/${name}`, role)
}
