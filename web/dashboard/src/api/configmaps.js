import {del, get, put, post} from "@/plugins/request"


const configMapUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/configmaps`
}

const namespaceMapUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/configmaps`
}

export function listConfigMaps (cluster_name, search, keywords, pageNum, pageSize) {
  let url = configMapUrl(cluster_name)
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

export function listConfigMapsWithNs (cluster_name, namespace) {
  return get(`${namespaceMapUrl(cluster_name, namespace)}`)
}

export function getConfigMap (cluster_name, namespace, name) {
  return get(`${namespaceMapUrl(cluster_name, namespace)}/${name}`)
}

export function createConfigMap (cluster_name, namespace, data) {
  return post(`${namespaceMapUrl(cluster_name, namespace)}`, data)
}

export function deleteConfigMap (cluster_name, namespace, name) {
  return del(`${namespaceMapUrl(cluster_name, namespace)}/${name}`)
}

export function updateConfigMap (cluster_name, namespace, name, data) {
  return put(`${namespaceMapUrl(cluster_name, namespace)}/${name}`, data)
}
