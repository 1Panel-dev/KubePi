import {post, get, del, put} from "@/plugins/request"

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
    let url = `${baseUrl}/search?pageNum=${page}&pageSize=${size}&showExtra=true`
    return post(url, {conditions: conditions})
}


export function listClusterMembers(name) {
    return get(`${baseUrl}/${name}/members`)
}

export function createClusterMember(name, member) {
    return post(`${baseUrl}/${name}/members`, member)
}

export function updateCluster(name, req) {
    return put(`${baseUrl}/${name}`, req)
}


export function updateClusterMember(name, memberName, member) {
    return put(`${baseUrl}/${name}/members/${memberName}`, member)
}


export function getClusterMember(name, memberName) {
    return get(`${baseUrl}/${name}/members/${memberName}`)
}

export function deleteClusterMember(name, memberName, kind) {
    return del(`${baseUrl}/${name}/members/${memberName}?kind=${kind}`)
}

export function listClusterRoles(name, scope) {

    let url = `${baseUrl}/${name}/clusterroles`
    if (scope) {
        url += `?scope=${scope}`
    }
    return get(url)
}

export function createClusterRole(name, clusterRole) {
    return post(`${baseUrl}/${name}/clusterroles`, clusterRole)
}

export function updateClusterRole(name, clusterRoleName, clusterRole) {
    return put(`${baseUrl}/${name}/clusterroles/${clusterRoleName}`, clusterRole)
}

export function listClusterApiGroups(name) {
    return get(`${baseUrl}/${name}/apigroups`)
}

export function listClusterResourceByGroupVersion(name, groupVersion, scope) {
    return get(`${baseUrl}/${name}/apigroups/${groupVersion}?scope=${scope}`)
}

export function deleteClusterRole(name, clusterRoleName) {
    return del(`${baseUrl}/${name}/clusterroles/${clusterRoleName}`)
}

export function listNamespaces(name) {
    return get(`${baseUrl}/${name}/namespaces`)
}

export function listClusterRepos(name) {
  return get(`${baseUrl}/${name}/repos`)
}

export function listClusterReposDetail(name) {
return get(`${baseUrl}/${name}/repos/detail`)
}

export function addClusterRepo(name,data) {
  return post(`${baseUrl}/${name}/repos`, data)
}

export function deleteClusterRepo(cluster,repo) {
  return del(`${baseUrl}/${cluster}/repos/${repo}`)
}
