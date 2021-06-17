import {get} from "@/plugins/request"

const eventUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/events`
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
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}