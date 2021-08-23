import {del, get, put, post} from "@/plugins/request"

const clusterRoleUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/clusterroles`
}

export function listClusterRoles (cluster_name, search, keywords, pageNum, pageSize) {
  let url = clusterRoleUrl(cluster_name)
  const params = {}
  if (search) {
    params["search"] = true
    if (keywords) {
      params["keywords"] = keywords
    }
    if (pageNum && pageSize) {
      params["pageNum"] = pageNum
      params["pageSize"] = pageSize
    }
  }
  return get(url, params)
}

export function deleteClusterRole (cluster_name, name) {
  return del(`${clusterRoleUrl(cluster_name)}/${name}`)
}

export function getClusterRole (cluster_name, name) {
  return get(`${clusterRoleUrl(cluster_name)}/${name}`)
}

export function createClusterRole (cluster_name, data) {
  return post(`${clusterRoleUrl(cluster_name)}`, data)
}

export function updateClusterRole (cluster_name, name, data) {
  return put(`${clusterRoleUrl(cluster_name)}/${name}`, data)
}

