import {del, get, patch, post} from "@/plugins/request"

const clusterRoleBindingUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/clusterrolebindings`
}

export function listClusterRoleBindings (cluster_name, limit, continueToken, search) {
  let url = clusterRoleBindingUrl(cluster_name)
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

export function deleteClusterRoleBinding (cluster_name, name) {
  return del(`${clusterRoleBindingUrl(cluster_name)}/${name}`)
}

export function getClusterRoleBinding (cluster_name, name) {
  return get(`${clusterRoleBindingUrl(cluster_name)}/${name}`)
}

export function createClusterRoleBinding (cluster_name, data) {
  return post(`${clusterRoleBindingUrl(cluster_name)}`, data)
}

export function updateClusterRoleBinding (cluster_name, name, data) {
  return patch(`${clusterRoleBindingUrl(cluster_name)}/${name}`, data)
}

