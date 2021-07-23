import {del, get, patch, post} from "@/plugins/request"

const scUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/storage.k8s.io/v1/storageclasses`
}


export function listStorageClasses (cluster_name, limit, continueToken, search) {
  let url = scUrl(cluster_name)
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


export function deleteStorageClasses (cluster_name, name) {
  return del(`${scUrl(cluster_name,)}/${name}`)
}

export function getStorageClass (cluster_name, name) {
  return get(`${scUrl(cluster_name,)}/${name}`)
}

export function creategetStorageClass (cluster_name, data) {
  return post(`${scUrl(cluster_name,)}`, data)
}

export function editgetStorageClass (cluster_name, data) {
  return patch(`${scUrl(cluster_name,)}`, data)
}
