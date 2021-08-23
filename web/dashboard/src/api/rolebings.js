import {del, get, put, post,} from "@/plugins/request"

const roleBindingUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/rolebindings`
}
const namespaceRoleBindingUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/rolebindings`
}

export function listRoleBindings (cluster_name, search, keywords, pageNum, pageSize) {
  let url = roleBindingUrl(cluster_name)
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

export function deleteRoleBinding (cluster_name, namespace, name) {
  return del(`${namespaceRoleBindingUrl(cluster_name, namespace)}/${name}`)
}

export function getRoleBinding (cluster_name, namespace, name) {
  return get(`${namespaceRoleBindingUrl(cluster_name, namespace)}/${name}`)
}

export function createRoleBinding (cluster_name, namespace, data) {
  return post(`${namespaceRoleBindingUrl(cluster_name,namespace)}`, data)
}

export function updateRoleBinding (cluster_name, namespace,name, data) {
  return put(`${namespaceRoleBindingUrl(cluster_name,namespace)}/${name}`, data)
}
