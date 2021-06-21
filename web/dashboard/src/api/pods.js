import {get} from "@/plugins/request"

const podUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/pods`
}

export function listPods (cluster_name, limit, continueToken, search) {
  let url = podUrl(cluster_name)
  if (limit) {
    url += "?limit=" + limit
  }
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  if (search && search !== "") {
    url += "&fieldSelector=spec.nodeName=" + search
  }
  return get(url)
}

export function listPodsByNode (cluster_name, node) {
  let url = podUrl(cluster_name)
  if (node && node !== "") {
    url += "?fieldSelector=spec.nodeName=" + node
  }
  return get(url)
}