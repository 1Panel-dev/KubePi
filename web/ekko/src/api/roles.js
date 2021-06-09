import {get, post} from "@/plugins/request"

const baseUrl = "/api/v1/roles"

export function searchRoles(pageNum, pageSize, conditions) {
    return post(`${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`, conditions)
}

export function listRoles() {
    return get(baseUrl)
}

export function createRole(role) {
    return post(baseUrl, role)
}