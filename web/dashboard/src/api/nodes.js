import {get} from "@/plugins/request"

const nodeUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/nodes`
}


export function listNodes (cluster_name, limit, continueToken, search) {
  let url = nodeUrl(cluster_name)
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}

export function getNode (cluster_name, name) {
  return get(`${nodeUrl(cluster_name)}/${name}`)
}