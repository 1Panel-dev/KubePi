import {del, get, put, post} from "@/plugins/request"

const serviceAccountUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/serviceaccounts`
}
const namespaceServiceAccountUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/serviceaccounts`
}

export function listServiceAccounts (cluster_name,search, keywords, pageNum, pageSize) {
  let url = serviceAccountUrl(cluster_name)
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

export function listServiceAccountsWithNs (cluster_name, namespace) {
  return get(`${namespaceServiceAccountUrl(cluster_name, namespace)}`)
}

export function deleteServiceAccount (cluster_name, namespace, name) {
  return del(`${namespaceServiceAccountUrl(cluster_name, namespace)}/${name}`)
}

export function getServiceAccount (cluster_name, namespace, name) {
  return get(`${namespaceServiceAccountUrl(cluster_name, namespace)}/${name}`)
}

export function createServiceAccount (cluster_name, namespace, data) {
  return post(`${namespaceServiceAccountUrl(cluster_name, namespace)}`, data)
}

export function updateServiceAccount (cluster_name, namespace, name, data) {
  return put(`${namespaceServiceAccountUrl(cluster_name, namespace)}/${name}`, data)
}

