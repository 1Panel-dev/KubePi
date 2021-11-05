import {del, get, put, post} from "@/plugins/request"


const ingressUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/ingresses`
}

const namespaceIngressUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/namespaces/${namespace}/ingresses`
}

export function listIngresses (cluster_name, search, keywords, pageNum, pageSize) {
  let url = ingressUrl(cluster_name)
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

export function deleteIngress (cluster_name, namespace, name) {
  return del(`${namespaceIngressUrl(cluster_name, namespace)}/${name}`)
}

export function listIngressWithNs (cluster_name, namespace) {
  return get(`${namespaceIngressUrl(cluster_name, namespace)}`)
}

export function getIngress (cluster_name, namespace, name) {
  return get(`${namespaceIngressUrl(cluster_name, namespace)}/${name}`)
}

export function createIngress (cluster_name, namespace, data) {
  return post(`${namespaceIngressUrl(cluster_name, namespace)}`, data)
}

export function updateIngress (cluster_name, namespace, name, data) {
  return put(`${namespaceIngressUrl(cluster_name, namespace)}/${name}`, data)
}
