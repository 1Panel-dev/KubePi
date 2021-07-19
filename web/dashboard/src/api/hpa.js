import {del, get, post,} from "@/plugins/request"

const hpaUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/autoscaling/v2beta2/horizontalpodautoscalers`
}
const namespaceHpaUrl = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/autoscaling/v2beta2/namespaces/${namespace}/horizontalpodautoscalers`
}

export function listHpas (cluster_name, limit, continueToken, search) {
  let url = hpaUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector = "metadata.namespace=" + search
  }
  return get(url, param)
}

export function deleteHpa (cluster_name, namespace, name) {
  return del(`${namespaceHpaUrl(cluster_name, namespace)}/${name}`)
}

export function getHpa (cluster_name, namespace, name) {
  return get(`${namespaceHpaUrl(cluster_name, namespace)}/${name}/status`)
}

export function createHpa (cluster_name, namespace, data) {
  return post(`${namespaceHpaUrl(cluster_name,namespace)}`, data)
}
