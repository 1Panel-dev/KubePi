import {post} from "@/plugins/request"

const baseUrl = "/api/v1/users"

export function searchUsers(pageNum, pageSize, conditions) {
    return post(`${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`, conditions)
}