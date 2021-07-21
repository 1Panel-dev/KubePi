import {get} from "@/plugins/request"

const scUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/storage.k8s.io/v1/storageclasses`
}

export function listStorageClass (cluster_name, limit, continueToken) {
  let url = scUrl(cluster_name)
  if (limit) {
    url += "?limit=" + limit
  }
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  return get(url)
}