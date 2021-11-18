import {post} from "@/plugins/request"

const baseUrl = "/api/v1/systems"

export function searchLoginLogs(pageNum, pageSize, conditions) {
    let url = `${baseUrl}/login/logs/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    return post(url, {conditions: conditions})
}

export function searchOperationLogs(pageNum, pageSize, conditions) {
    let url = `${baseUrl}/operation/logs/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    return post(url, {conditions: conditions})
}
