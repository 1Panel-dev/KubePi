import {del, get, patch, post} from "@/plugins/request"

const namespaceUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/namespaces`
}


export function listNamespace (cluster_name, limit, continueToken, search) {
  let url = namespaceUrl(cluster_name)
  url += "?limit=" + limit
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}

export function createNamespace (cluster_name, data) {
  return post(`${namespaceUrl(cluster_name)}`, data)
}

export function getNamespace (cluster_name, namespace) {
  return get(`${namespaceUrl(cluster_name)}/${namespace}`)
}

export function deleteNamespace (cluster_name, namespace) {
  return del(`${namespaceUrl(cluster_name)}/${namespace}`)
}

export function editNamespace (cluster_name, namespace, data) {
  return patch(`${namespaceUrl(cluster_name)}/${namespace}`, data)
}
