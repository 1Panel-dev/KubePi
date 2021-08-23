import {del, get, put, post} from "@/plugins/request"


const endpointUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/endpoints`
}

const namespaceEndpointUrl = (cluster_name,namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/endpoints`
}

export function listEndPoints (cluster_name, search, keywords, pageNum, pageSize) {
  let url = endpointUrl(cluster_name)
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

export function deleteEndPoint (cluster_name, namespace, name) {
  return del(`${namespaceEndpointUrl(cluster_name, namespace)}/${name}`)
}

export function createEndPoint(cluster_name, namespace, data) {
  return post(`${namespaceEndpointUrl(cluster_name, namespace)}`, data)
}

export function getEndPoint (cluster_name, namespace, name) {
  return get(`${namespaceEndpointUrl(cluster_name, namespace)}/${name}`)
}

export function updateEndPoint (cluster_name, namespace, name, data) {
  return put(`${namespaceEndpointUrl(cluster_name, namespace)}/${name}`, data)
}
