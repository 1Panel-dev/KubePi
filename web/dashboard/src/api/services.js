import {del, get} from "@/plugins/request"

const serviceUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/services`
}

const namespaceServiceUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/services`
}


export function listServices (cluster_name, limit, continueToken, search) {
  let url = serviceUrl(cluster_name)
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

export function deleteService (cluster_name, namespace, name) {
  return del(`${namespaceServiceUrl(cluster_name, namespace)}/${name}`)
}
