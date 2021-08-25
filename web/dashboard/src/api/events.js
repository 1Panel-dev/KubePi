import {get} from "@/plugins/request"

const eventUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/events`
}
const eventUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/events`
}


export function listEvents (cluster_name, search, keywords, pageNum, pageSize) {
  let url = eventUrl(cluster_name)
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
  params["sortBy"] = ".metadata.creationTimestamp"
  return get(url, params)
}

export function listEventsWithNsSelector (cluster_name, namespace, selectors) {
  let url = eventUrlWithNs(cluster_name, namespace)
  const param = {}
  if (selectors && selectors !== "") {
    param.fieldSelector = selectors
  }
  return get(url, param)
}
