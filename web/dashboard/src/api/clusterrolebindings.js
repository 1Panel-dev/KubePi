import {del, get, put, post} from "@/plugins/request"

const clusterRoleBindingUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/rbac.authorization.k8s.io/v1/clusterrolebindings`
}

export function listClusterRoleBindings (cluster_name, search, keywords, pageNum, pageSize) {
  let url = clusterRoleBindingUrl(cluster_name)
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
  return put(`${clusterRoleBindingUrl(cluster_name)}/${name}`, data)
}

