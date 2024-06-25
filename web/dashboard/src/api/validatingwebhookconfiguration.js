import {del, get, put, post, patch} from "@/plugins/request"

const vcUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/admissionregistration.k8s.io/v1/validatingwebhookconfigurations`
}


export function listValidatingwebhookconfigurations (cluster_name, search, keywords, pageNum, pageSize) {
  let url = vcUrl(cluster_name)
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


export function deleteValidatingwebhookconfigurations (cluster_name, name) {
  return del(`${vcUrl(cluster_name,)}/${name}`)
}

export function getValidatingwebhookconfiguration (cluster_name, name) {
  return get(`${vcUrl(cluster_name,)}/${name}`)
}

export function createValidatingwebhookconfiguration (cluster_name, data) {
  return post(`${vcUrl(cluster_name,)}`, data)
}

export function updateValidatingwebhookconfiguration (cluster_name, name, data) {
  return put(`${vcUrl(cluster_name,)}/${name}`, data)
}

export function changeValidatingwebhookconfiguration(cluster_name,name,data) {
  return patch(`${vcUrl(cluster_name,)}/${name}`, data)
}

