import {del, get, put, post} from "@/plugins/request"

const roleUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/roles`
}

const namespaceRoleUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/roles`
}

export function listRoles (cluster_name, namespace, search, keywords, pageNum, pageSize) {
  let url = roleUrl(cluster_name)
  if (namespace !== "") {
    url = namespaceRoleUrl(cluster_name, namespace)
  }
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

export function listNamespaceRoles(cluster_name, namespace ) {
  let url = namespaceRoleUrl(cluster_name, namespace)
  return get(url)
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
  return put(`${namespaceRoleUrl(cluster_name, namespace)}/${name}`, data)
}

