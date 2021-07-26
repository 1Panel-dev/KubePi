import {del, get, patch, post} from "@/plugins/request"


const ingressUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/ingresses`
}

const namespaceIngressUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/namespaces/${namespace}/ingresses`
}

export function listIngresses (cluster_name, limit, continueToken, search) {
  let url = ingressUrl(cluster_name)
  if (limit) {
    url += "?limit=" + limit
  }
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}

export function deleteIngress (cluster_name, name) {
  return del(`${namespaceIngressUrl(cluster_name)}/${name}`)
}

export function getIngress (cluster_name, namespace, name) {
  return get(`${namespaceIngressUrl(cluster_name, namespace)}/${name}`)
}

export function createIngress (cluster_name, namespace, data) {
  return post(`${namespaceIngressUrl(cluster_name, namespace)}`, data)
}

export function updateIngress (cluster_name, namespace, name, data) {
  return patch(`${namespaceIngressUrl(cluster_name, namespace)}/${name}`, data)
}
