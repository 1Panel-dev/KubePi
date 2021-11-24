import {get, del, post, put} from "@/plugins/request"

const podUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/pods`
}
const podUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/pods`
}

const evictionUrl = (cluster_name,namespace,pod) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/pods/${pod}/eviction`
}

export function listPods (cluster_name, search) {
  let url = podUrl(cluster_name)
  const param = {}
  if (search && search !== "") {
    param.fieldSelector = "spec.nodeName=" + search
  }
  return get(url, param)
}

export function listPodsWithNsSelector (cluster_name, namespace, selectors, fieldSelectors) {
  let url = podUrl(cluster_name)
  if (namespace && namespace !== "") {
    url = podUrlWithNs(cluster_name, namespace)
  }
  const param = {}
  if (selectors && selectors !== "") {
    param.labelSelector = selectors
  }
  if (fieldSelectors && fieldSelectors !== "") {
    param.fieldSelector = fieldSelectors
  }
  return get(url, param)
}

export function getPodByName (cluster_name, namespace, pod) {
  return get(`${podUrlWithNs(cluster_name, namespace)}/${pod}`)
}

export function deletePod (cluster_name, pod) {
  return del(`${podUrl(cluster_name)}/${pod}`)
}

export function getPodLogsByName (cluster_name, namespace, pod, params) {
  return get(`${podUrlWithNs(cluster_name, namespace)}/${pod}/log`, params)
}

export function createPod (cluster_name, namespace, pod) {
  return post(`${podUrlWithNs(cluster_name, namespace)}`, pod)
}

export function updatePod (cluster_name, namespace,name, pod) {
  return put(`${podUrlWithNs(cluster_name, namespace)}/${name}`, pod)
}

export function evictionPod(cluster_name, namespace,name,data) {
  return post(`${evictionUrl(cluster_name, namespace,name)}`, data)
}
