import {del, get, patch, post} from "@/plugins/request"


const limitRangeUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/limitranges`
}

const namespaceLimitRangeUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/limitranges`
}

export function listLimitRanges (cluster_name, limit, continueToken, search) {
  let url = limitRangeUrl(cluster_name)
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

export function getLimitRange (cluster_name, namespace, name) {
  return get(`${namespaceLimitRangeUrl(cluster_name, namespace)}/${name}`)
}

export function deleteLimitRange (cluster_name, namespace, name) {
  return del(`${namespaceLimitRangeUrl(cluster_name, namespace)}/${name}`)
}

export function createLimitRange (cluster_name, namespace, data) {
  return post(`${namespaceLimitRangeUrl(cluster_name, namespace)}`, data)
}

export function updateLimitRange (cluster_name, namespace, name,data) {
  return patch(`${namespaceLimitRangeUrl(cluster_name, namespace)}/${name}`, data)
}
