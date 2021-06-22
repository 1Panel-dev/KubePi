import {get} from "@/plugins/request"

const nodeUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/nodes`
}


export function listNodes (cluster_name, limit, continueToken, search) {
  let url = nodeUrl(cluster_name)
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

export function getNode (cluster_name, name) {
  return get(`${nodeUrl(cluster_name)}/${name}`)
}