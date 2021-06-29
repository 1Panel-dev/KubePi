import {get} from "@/plugins/request"

const pvUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumes`
}


export function listPvs (cluster_name, limit, continueToken, search) {
  let url = pvUrl(cluster_name)
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