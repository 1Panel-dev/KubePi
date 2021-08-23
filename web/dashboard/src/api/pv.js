import {del, get, put, post} from "@/plugins/request"

const pvUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumes`
}


export function listPvs (cluster_name, search, keywords, pageNum, pageSize) {
  let url = pvUrl(cluster_name)
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

export function listPvsWithNs (cluster_name, namespace) {
  return get(`${pvUrl(cluster_name, namespace)}`)
}

export function deletePvs (cluster_name, name) {
  return del(`${pvUrl(cluster_name,)}/${name}`)
}

export function getPv (cluster_name, name) {
  return get(`${pvUrl(cluster_name,)}/${name}`)
}

export function createPv (cluster_name, data) {
  return post(`${pvUrl(cluster_name,)}`, data)
}

export function updatePv (cluster_name, data, name) {
  return put(`${pvUrl(cluster_name,)}/${name}`, data)
}
