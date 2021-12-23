import {get, post, patch} from "@/plugins/request"

const apiUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis`
}
const apiServerUrl = (cluster_name) => {
  return `${apiUrl(cluster_name)}/apiregistration.k8s.io/v1/apiservices`
}

export function listApis (cluster_name) {
  return get(`${apiUrl(cluster_name)}`)
}

export function createApiServer (cluster_name, data) {
  return post(`${apiServerUrl(cluster_name)}`, data)
}
export function updateApiServer (cluster_name, name, data) {
  return patch(`${apiServerUrl(cluster_name)}/${name}`, data)
}
export function getApiService (cluster_name, name) {
  return get(`${apiServerUrl(cluster_name)}/${name}`)
}
export function listApiServices (cluster_name) {
  return get(`${apiServerUrl(cluster_name)}`)
}

export function listResourceApi (cluster_name, apiGroup, version) {
  return get(`${apiUrl(cluster_name)}/${apiGroup}/${version}`)
}

export function listPodMetrics (cluster_name, namespace,labelSelector,name) {

  const param = {}
  if (namespace && namespace !== 'All Namespaces') {
    param.fieldSelector = "metadata.namespace=" + namespace
  }
  if (name) {
    param.fieldSelector = "metadata.name=" + name
  }

  if (labelSelector) {
    param.labelSelector =  labelSelector
  }

  return get(`${apiUrl(cluster_name)}/metrics.k8s.io/v1beta1/pods`, param)
}


export function listNodeMetrics (cluster_name) {
  return get(`${apiUrl(cluster_name)}/metrics.k8s.io/v1beta1/nodes`)
}