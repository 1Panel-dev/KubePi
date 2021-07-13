import {del, get, patch, post} from "@/plugins/request"

const pvUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumes`
}


export function listPvs (cluster_name, limit, continueToken, search) {
  let url = pvUrl(cluster_name)
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

export function listPvsWithNs (cluster_name, namespace) {
  return get(`${pvUrl(cluster_name, namespace)}`)
}

export function deletePvs (cluster_name, namespace, name) {
  return del(`${pvUrl(cluster_name, namespace)}/${name}`)
}

export function getSecret (cluster_name, namespace, name) {
  return get(`${pvUrl(cluster_name, namespace)}/${name}`)
}

export function createSecret (cluster_name, namespace, data) {
  return post(`${pvUrl(cluster_name, namespace)}`, data)
}

export function editSecret (cluster_name, namespace, data) {
  return patch(`${pvUrl(cluster_name, namespace)}`, data)
}
