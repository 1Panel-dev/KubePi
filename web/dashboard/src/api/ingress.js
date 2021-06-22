import {get} from "@/plugins/request"


const ingressUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/apis/networking.k8s.io/v1/ingresses`
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