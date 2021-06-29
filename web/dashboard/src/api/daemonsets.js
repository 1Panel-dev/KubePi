import {get} from "@/plugins/request"


const daemonSetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/daemonsets`
}

export function listDaemonSets (cluster_name, limit, continueToken, search) {
  let url = daemonSetUrl(cluster_name)
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