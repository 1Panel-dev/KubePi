import {get} from "@/plugins/request"

const secretUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/secrets`
}
const namespaceSecretUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/secrets`
}

export function listSecrets (cluster_name, limit, continueToken, search) {
  let url = secretUrl(cluster_name)
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

export function listSecretsWithNs (cluster_name, namespace) {
  return get(`${namespaceSecretUrl(cluster_name, namespace)}`)
}

export function delSecrets () {
}

export function getSecret (cluster_name, namespace, name) {
  return get(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`)
}


