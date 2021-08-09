import {del, get, patch, post} from "@/plugins/request"

const customResourceUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apiextensions.k8s.io/v1/customresourcedefinitions`
}

export function listCustomResources (cluster_name, search, keywords, pageNum, pageSize) {
  let url = customResourceUrl(cluster_name)
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

export function deleteCustomResource (cluster_name, name) {
  return del(`${customResourceUrl(cluster_name)}/${name}`)
}

export function getCustomResource (cluster_name, name) {
  return get(`${customResourceUrl(cluster_name)}/${name}`)
}

export function createCustomResource (cluster_name, data) {
  return post(`${customResourceUrl(cluster_name)}`, data)
}

export function updateCustomResource (cluster_name, name, data) {
  return patch(`${customResourceUrl(cluster_name)}/${name}`, data)
}

