import {del, get, patch, post} from "@/plugins/request"

const resourceQuotaUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/resourcequotas`
}

const namespaceResourceQuotaUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/resourcequotas`
}

export function listResourceQuotas (cluster_name, limit, continueToken, search) {
  let url = resourceQuotaUrl(cluster_name)
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

export function deleteResourceQuota (cluster_name, namespace, name) {
  return del(`${namespaceResourceQuotaUrl(cluster_name, namespace)}/${name}`)
}

export function getResourceQuota (cluster_name, namespace, name) {
  return get(`${namespaceResourceQuotaUrl(cluster_name, namespace)}/${name}`)
}

export function createResourceQuota (cluster_name, namespace, data) {
  return post(`${namespaceResourceQuotaUrl(cluster_name, namespace)}`, data)
}

export function updateResourceQuota (cluster_name, namespace,name, data) {
  return patch(`${namespaceResourceQuotaUrl(cluster_name, namespace)}/${name}`, data)
}