import {post, get, del} from "@/plugins/request"

const baseUrl = "/api/v1/clusters"

export function createCluster(data) {
    return post(baseUrl, data)
}

export function listClusters() {
    return get(baseUrl)
}

export function deleteBy(name) {
    return del(`${baseUrl}/${name}`)
}

export function searchClusters(page, size, conditions) {
    return post(`${baseUrl}/search?pageNum=${page}&pageSize=${size}`, conditions)
}

