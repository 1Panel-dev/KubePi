import {del, get, put, post} from "@/plugins/request"


const networkPolicyUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/networkpolicies`
}

const namespaceNetworkPolicyUrlUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/namespaces/${namespace}/networkpolicies`
}

export function listNetworkPolicies (cluster_name,search, keywords, pageNum, pageSize) {
  let url = networkPolicyUrl(cluster_name)
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
  return get(url,params)
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
  return put(`${namespaceNetworkPolicyUrlUrl(cluster_name, namespace)}/${name}`, data)
}
