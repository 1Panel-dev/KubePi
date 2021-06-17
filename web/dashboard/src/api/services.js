import {get} from "@/plugins/request"

const serviceUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/services`
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