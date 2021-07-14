import {get} from "@/plugins/request"

const eventUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/events`
}
const eventUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/events`
}


export function listEvents (cluster_name, limit, continueToken, search) {
  let url = eventUrl(cluster_name)
  if (limit) {
    url += "?limit=" + limit
  }
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  if (search && search !== "") {
    url += "&fieldSelector=metadata.namespace=" + search
  }
  return get(url)
}

export function listEventsWithNsSelector (cluster_name, namespace, selectors) {
  let url = eventUrlWithNs(cluster_name, namespace)
  const param = {}
  if (selectors && selectors !== "") {
    param.fieldSelector = selectors
  }
  return get(url,param)
}
