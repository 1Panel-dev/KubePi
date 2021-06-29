import {get} from "@/plugins/request"


const statefulSetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/statefulsets`
}

export function listStatefulSets (cluster_name, limit, continueToken, search) {
  let url = statefulSetUrl(cluster_name)
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