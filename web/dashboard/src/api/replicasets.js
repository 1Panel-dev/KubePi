import {get} from "@/plugins/request"

const replicaSetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/replicasets`
}

const replicaSetNamespaceUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/replicasets`
}

export function listReplicaSets (cluster_name, search, keywords, pageNum, pageSize) {
  let url = replicaSetUrl(cluster_name)
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
  return get(url, params)
}

export function listNsReplicaSets (cluster_name, namespace) {
  return get(`${replicaSetNamespaceUrl(cluster_name, namespace)}`)
}
