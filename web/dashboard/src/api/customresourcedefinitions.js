import {del, get, put, post} from "@/plugins/request"

const customResourceDefinitionUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apiextensions.k8s.io/v1/customresourcedefinitions`
}

export function listCustomResourceDefinitions (cluster_name, search, keywords, pageNum, pageSize) {
  let url = customResourceDefinitionUrl(cluster_name)
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

export function deleteCustomResourceDefinition (cluster_name, name) {
  return del(`${customResourceDefinitionUrl(cluster_name)}/${name}`)
}

export function getCustomResourceDefinition (cluster_name, name) {
  return get(`${customResourceDefinitionUrl(cluster_name)}/${name}`)
}

export function createCustomResourceDefinition (cluster_name, data) {
  return post(`${customResourceDefinitionUrl(cluster_name)}`, data)
}

export function updateCustomResourceDefinitione (cluster_name, name, data) {
  return put(`${customResourceDefinitionUrl(cluster_name)}/${name}`, data)
}

const resourceUrl = (cluster_name, version, group, resource_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/${group}/${version}/${resource_name}`
}
const resourceNSUrl = (cluster_name, version, group, resource_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/${group}/${version}/namespaces/${namespace}/${resource_name}`
}

export function listResourceByGroup (cluster_name, version, group, resource_name, search, keywords, pageNum, pageSize) {

  let url = resourceUrl(cluster_name, version, group, resource_name)
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

export function deleteResource (cluster_name, version, group, names, namespace, name) {
  let url = resourceUrl(cluster_name, version, group, names)
  if (namespace && namespace !== "") {
    url = resourceNSUrl(cluster_name, version, group, names, namespace)
  }
  return del(`${url}/${name}`)
}

export function getResource (cluster_name, version, group, names, namespace, name) {
  let url = resourceUrl(cluster_name, version, group, names)
  if (namespace && namespace !== "") {
    url = resourceNSUrl(cluster_name, version, group, names, namespace)
  }
  return get(`${url}/${name}`)
}

export function updateResource(cluster_name, version, group, names, namespace, name,data) {
  let url = resourceUrl(cluster_name, version, group, names)
  if (namespace && namespace !== "") {
    url = resourceNSUrl(cluster_name, version, group, names, namespace)
  }
  return put(`${url}/${name}`,data)
}
