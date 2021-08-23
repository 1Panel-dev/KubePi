import {del, get, put, post} from "@/plugins/request"

const secretUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/secrets`
}
const namespaceSecretUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/secrets`
}

export function listSecrets (cluster_name, search, keywords, pageNum, pageSize) {
  let url = secretUrl(cluster_name)
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

export function listSecretsWithNs (cluster_name, namespace) {
  return get(`${namespaceSecretUrl(cluster_name, namespace)}`)
}

export function deleteSecrets (cluster_name, namespace, name) {
  return del(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`)
}

export function getSecret (cluster_name, namespace, name) {
  return get(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`)
}

export function createSecret (cluster_name, namespace, data) {
  return post(`${namespaceSecretUrl(cluster_name, namespace)}`, data)
}

export function editSecret (cluster_name, namespace, name, data) {
  return put(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`, data)
}

