import {del, get, put, post} from "@/plugins/request"

const secretUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/secrets`
}

export function listSecretsWithLabelSelector (cluster_name,labelSelector) {
  return get(`${secretUrl(cluster_name)}`,{labelSelector:labelSelector})
}

