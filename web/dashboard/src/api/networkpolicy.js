import {del, get, patch, post} from "@/plugins/request"


const networkPolicyUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/networkpolicies`
}

const namespaceNetworkPolicyUrlUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/namespaces/${namespace}/networkpolicies`
}

export function listNetworkPolicies (cluster_name,search) {
  let url = networkPolicyUrl(cluster_name)
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}

export function deletePolicyUrl (cluster_name, namespace, name) {
  return del(`${namespaceNetworkPolicyUrlUrl(cluster_name, namespace)}/${name}`)
}

export function createNetworkPolicy (cluster_name, namespace, data) {
  return post(`${namespaceNetworkPolicyUrlUrl(cluster_name, namespace)}`, data)
}

export function getNetworkPolicy (cluster_name, namespace, name) {
  return get(`${namespaceNetworkPolicyUrlUrl(cluster_name, namespace)}/${name}`)
}

export function updateNetworkPolicy (cluster_name, namespace, name, data) {
  return patch(`${namespaceNetworkPolicyUrlUrl(cluster_name, namespace)}/${name}`, data)
}
