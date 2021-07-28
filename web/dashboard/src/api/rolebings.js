import {del, get, patch, post,} from "@/plugins/request"

const roleBindingUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/rolebindings`
}
const namespaceRoleBindingUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/namespaces/${namespace}/rolebindings`
}

export function listRoleBindings (cluster_name, limit, continueToken, search) {
  let url = roleBindingUrl(cluster_name)
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
  return patch(`${namespaceRoleBindingUrl(cluster_name,namespace)}/${name}`, data)
}
