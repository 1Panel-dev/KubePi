import {del, get, patch, post} from "@/plugins/request"


const configMapUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/configmaps`
}

const namespaceMapUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/configmaps`
}

export function listConfigMaps (cluster_name, limit, continueToken, search) {
  let url = configMapUrl(cluster_name)
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

export function listConfigMapsWithNs (cluster_name, namespace) {
  return get(`${namespaceMapUrl(cluster_name, namespace)}`);
}

export function getConfigMap (cluster_name, namespace, name) {
  return get(`${namespaceMapUrl(cluster_name, namespace)}/${name}`)
}

export function createConfigMap (cluster_name, namespace, data) {
  return post(`${namespaceMapUrl(cluster_name, namespace)}`, data)
}

export function deleteConfigMap(cluster_name,namespace, name) {
  return del(`${namespaceMapUrl(cluster_name, namespace)}/${name}`)
}

export function updateConfigMap(cluster_name, namespace,name, data) {
  return patch(`${namespaceMapUrl(cluster_name, namespace)}/${name}`, data)
}
