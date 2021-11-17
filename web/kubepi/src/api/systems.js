import {post} from "@/plugins/request"

const baseUrl = "/api/v1/systems"

export function searchSystemLog(pageNum, pageSize, conditions) {
    let url = `${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    return post(url, {conditions: conditions})
}
