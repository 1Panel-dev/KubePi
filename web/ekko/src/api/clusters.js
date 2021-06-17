import {post, get, del} from "@/plugins/request"

const baseUrl = "/api/v1/clusters"

export function createCluster(data) {
    return post(baseUrl, data)
}

export function listClusters() {
    return get(baseUrl)
}

export function deleteCluster(name) {
    return del(`${baseUrl}/${name}`)
}

export function searchClusters(page, size, conditions) {
    return post(`${baseUrl}/search?pageNum=${page}&pageSize=${size}`, conditions)
}


export function listClusterMembers(name) {
    return get(`${baseUrl}/${name}/members`)
}

export function createClusterMember(name, member) {
    return post(`${baseUrl}/${name}/members`, member)
}

export function listClusterRoles(name) {
    return get(`${baseUrl}/${name}/clusterroles`)
}