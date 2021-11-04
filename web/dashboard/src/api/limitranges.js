import {del, get, put, post} from "@/plugins/request"


const limitRangeUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/limitranges`
}

const namespaceLimitRangeUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/limitranges`
}

export function listLimitRanges (cluster_name, search, keywords, pageNum, pageSize) {
  let url = limitRangeUrl(cluster_name)
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
  return put(`${namespaceLimitRangeUrl(cluster_name, namespace)}/${name}`, data)
}

export function listLimitRangeByNamespace(cluster_name,namespace) {
  return get(`${namespaceLimitRangeUrl(cluster_name, namespace)}`)
}
