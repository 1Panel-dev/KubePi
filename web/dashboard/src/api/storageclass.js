import {del, get, put, post, patch} from "@/plugins/request"

const scUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/storage.k8s.io/v1/storageclasses`
}


export function listStorageClasses (cluster_name, search, keywords, pageNum, pageSize) {
  let url = scUrl(cluster_name)
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


export function deleteStorageClasses (cluster_name, name) {
  return del(`${scUrl(cluster_name,)}/${name}`)
}

export function getStorageClass (cluster_name, name) {
  return get(`${scUrl(cluster_name,)}/${name}`)
}

export function createStorageClass (cluster_name, data) {
  return post(`${scUrl(cluster_name,)}`, data)
}

export function updateStorageClass (cluster_name, name, data) {
  return put(`${scUrl(cluster_name,)}/${name}`, data)
}

export function changeStorageClass(cluster_name,name,data) {
  return patch(`${scUrl(cluster_name,)}/${name}`, data)
}
