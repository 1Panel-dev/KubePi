import {del, get, put, post} from "@/plugins/request"

const serviceUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/services`
}

const namespaceServiceUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/services`
}


export function listServices (cluster_name, search, keywords, pageNum, pageSize) {
  let url = serviceUrl(cluster_name)
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
  return get(url,params)
}

export function listServicesWithNs (cluster_name, namespace) {
  let url = namespaceServiceUrl(cluster_name, namespace)
  return get(url)
}

export function deleteService (cluster_name, namespace, name) {
  return del(`${namespaceServiceUrl(cluster_name, namespace)}/${name}`)
}

export function getService (cluster_name, namespace, name) {
  return get(`${namespaceServiceUrl(cluster_name, namespace)}/${name}`)
}

export function createService (cluster_name, namespace, data) {
  return post(`${namespaceServiceUrl(cluster_name, namespace)}`, data)
}

export function updateService (cluster_name, namespace, name, data) {
  return put(`${namespaceServiceUrl(cluster_name, namespace)}/${name}`, data)
}
