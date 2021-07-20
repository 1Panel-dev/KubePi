import {del, get} from "@/plugins/request"


const networkPolicyUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/networkpolicies`
}

const namespaceNetworkPolicyUrlUrl = (cluster_name,namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/networking.k8s.io/v1/namespaces/${namespace}/networkpolicies`
}

export function listNetworkPolicies (cluster_name, limit, continueToken, search) {
  let url = networkPolicyUrl(cluster_name)
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

export function deletePolicyUrl (cluster_name, name) {
  return del(`${namespaceNetworkPolicyUrlUrl(cluster_name)}/${name}`)
}
