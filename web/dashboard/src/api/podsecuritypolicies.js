import {del, get, put, post} from "@/plugins/request"


const pspUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/policy/v1beta1/podsecuritypolicies`
}

export function listPSPs (cluster_name, search, keywords, pageNum, pageSize) {
  let url = pspUrl(cluster_name)
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
export function getPSP (cluster_name, name) {
  return get(`${pspUrl(cluster_name)}/${name}`)
}

export function createPSP (cluster_name, data) {
  return post(`${pspUrl(cluster_name)}`, data)
}

export function deletePSP (cluster_name, name) {
  return del(`${pspUrl(cluster_name)}/${name}`)
}

export function updatePSP (cluster_name, name, data) {
  return put(`${pspUrl(cluster_name)}/${name}`, data)
}
