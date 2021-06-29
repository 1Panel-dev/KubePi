import {get} from "@/plugins/request"

const secretUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/secrets`
}
const secretUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/secrets`
}

export function listSecrets (cluster_name) {
  return get(`${secretUrl(cluster_name)}`);
}

export function listSecretsWithNs (cluster_name, namespace) {
  return get(`${secretUrlWithNs(cluster_name, namespace)}`);
}