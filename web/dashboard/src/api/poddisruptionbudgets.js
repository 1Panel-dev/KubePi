import {del, get, put, post} from "@/plugins/request"


const pdbUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/policy/v1beta1/poddisruptionbudgets`
}

const namespacePDBUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/policy/v1beta1/namespaces/${namespace}/poddisruptionbudgets`
}

export function listPDBs (cluster_name, search, keywords, pageNum, pageSize) {
  let url = pdbUrl(cluster_name)
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

export function listPDBsWithNs (cluster_name, namespace) {
  return get(`${namespacePDBUrl(cluster_name, namespace)}`)
}

export function getPDB (cluster_name, namespace, name) {
  return get(`${namespacePDBUrl(cluster_name, namespace)}/${name}`)
}

export function createPDB (cluster_name,namespace, data) {
  return post(`${namespacePDBUrl(cluster_name,namespace)}`, data)
}

export function deletePDB (cluster_name, namespace, name) {
  return del(`${namespacePDBUrl(cluster_name, namespace)}/${name}`)
}

export function updatePDB (cluster_name, namespace, name, data) {
  return put(`${namespacePDBUrl(cluster_name, namespace)}/${name}`, data)
}
