import {get} from "@/plugins/request"

const replicaSetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/replicasets`
}

const replicaSetNamespaceUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/replicasets`
}

export function listReplicaSets (cluster_name, limit, continueToken, search) {
  let url = replicaSetUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector = "metadata.namespace=" + search
  }
  return get(url, param)
}

export function listNsReplicaSets (cluster_name, namespace) {
  return get(`${replicaSetNamespaceUrl(cluster_name, namespace)}`)
}
