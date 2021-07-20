import {del, get} from "@/plugins/request"


const endpointUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/endpoints`
}

const namespaceEndpointUrl = (cluster_name,namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/endpoints`
}

export function listEndPoints (cluster_name, limit, continueToken, search) {
  let url = endpointUrl(cluster_name)
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

export function deleteEndPoint (cluster_name, name) {
  return del(`${namespaceEndpointUrl(cluster_name)}/${name}`)
}
