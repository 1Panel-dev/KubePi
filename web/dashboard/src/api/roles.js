import {del, get, patch, post} from "@/plugins/request"

const roleUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/roles`
}

const namespaceRoleUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/roles`
}

export function listRoles (cluster_name,  search, keywords, pageNum, pageSize) {
  let url = roleUrl(cluster_name)
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

export function deleteRole (cluster_name, namespace, name) {
  return del(`${namespaceRoleUrl(cluster_name, namespace)}/${name}`)
}

export function getRole (cluster_name, namespace, name) {
  return get(`${namespaceRoleUrl(cluster_name, namespace)}/${name}`)
}

export function createRole (cluster_name, namespace, data) {
  return post(`${namespaceRoleUrl(cluster_name, namespace)}`, data)
}

export function updateRole (cluster_name, namespace, name, data) {
  return patch(`${namespaceRoleUrl(cluster_name, namespace)}/${name}`, data)
}

