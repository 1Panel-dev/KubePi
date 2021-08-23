import {del, get, put, post,} from "@/plugins/request"

const hpaUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/autoscaling/v2beta2/horizontalpodautoscalers`
}
const namespaceHpaUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/autoscaling/v2beta2/namespaces/${namespace}/horizontalpodautoscalers`
}

export function listHpas (cluster_name, search, keywords, pageNum, pageSize) {
  let url = hpaUrl(cluster_name)
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

export function deleteHpa (cluster_name, namespace, name) {
  return del(`${namespaceHpaUrl(cluster_name, namespace)}/${name}`)
}

export function getHpa (cluster_name, namespace, name) {
  return get(`${namespaceHpaUrl(cluster_name, namespace)}/${name}`)
}

export function createHpa (cluster_name, namespace, data) {
  return post(`${namespaceHpaUrl(cluster_name,namespace)}`, data)
}

export function updateHpa (cluster_name, namespace,name, data) {
  return put(`${namespaceHpaUrl(cluster_name,namespace)}/${name}`, data)
}
