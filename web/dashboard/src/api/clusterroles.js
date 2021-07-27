import {del, get, patch, post} from "@/plugins/request"

const clusterRoleUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/clusterroles`
}

export function listClusterRoles (cluster_name, limit, continueToken, search) {
  let url = clusterRoleUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector = "metadata.namespace=" + search
  }
  return get(url, param)
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
  return patch(`${clusterRoleUrl(cluster_name)}/${name}`, data)
}

