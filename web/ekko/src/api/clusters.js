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

export function createClusterRole(name, clusterRole) {
    return post(`${baseUrl}/${name}/clusterroles`, clusterRole)
}

export function listClusterApiGroups(name) {
    return get(`${baseUrl}/${name}/apigroups`)
}

export function listClusterResourceByGroupVersion(name, groupVersion) {
    return get(`${baseUrl}/${name}/apigroups/${groupVersion}`)
}

export function deleteClusterRole(name, clusterRoleName) {
    return del(`${baseUrl}/${name}/clusterroles/${clusterRoleName}`)
}
