import {get} from "@/plugins/request"

const apiUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis`
}

export function listApis (cluster_name) {
  return get(`${apiUrl(cluster_name)}`)
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
