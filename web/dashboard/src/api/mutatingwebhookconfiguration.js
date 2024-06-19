import {del, get, put, post, patch} from "@/plugins/request"

const mcUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/admissionregistration.k8s.io/v1/mutatingwebhookconfigurations`
}


export function listMutatingwebhookconfigurations (cluster_name, search, keywords, pageNum, pageSize) {
  let url = mcUrl(cluster_name)
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


export function deleteMutatingwebhookconfigurations (cluster_name, name) {
  return del(`${mcUrl(cluster_name,)}/${name}`)
}

export function getMutatingwebhookconfiguration (cluster_name, name) {
  return get(`${mcUrl(cluster_name,)}/${name}`)
}

export function createMutatingwebhookconfiguration (cluster_name, data) {
  return post(`${mcUrl(cluster_name,)}`, data)
}

export function updateMutatingwebhookconfiguration (cluster_name, name, data) {
  return put(`${mcUrl(cluster_name,)}/${name}`, data)
}

export function changeMutatingwebhookconfiguration(cluster_name,name,data) {
  return patch(`${mcUrl(cluster_name,)}/${name}`, data)
}

