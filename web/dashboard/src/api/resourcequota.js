import {del, get, put, post} from "@/plugins/request"

const resourceQuotaUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/resourcequotas`
}

const namespaceResourceQuotaUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/resourcequotas`
}

export function listResourceQuotas (cluster_name, search, keywords, pageNum, pageSize) {
  let url = resourceQuotaUrl(cluster_name)
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
  return put(`${namespaceResourceQuotaUrl(cluster_name, namespace)}/${name}`, data)
}

export function listResourceQuotaByNamespace(cluster_name,namespace) {
  return get(`${namespaceResourceQuotaUrl(cluster_name, namespace)}`)

}
