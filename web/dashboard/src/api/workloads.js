import {get, del, post, put} from "@/plugins/request"

const apiV1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/${type}`
}
const apiV1UrlWithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespaces}/${type}`
}
const appsV1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/${type}`
}
const appsV1UrlWithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespaces}/${type}`
}
const batchV1beta1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/${type}`
}
const batchV1beta1WithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/namespaces/${namespaces}/${type}`
}
const batchV1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/${type}`
}
const batchV1WithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/namespaces/${namespaces}/${type}`
}

export function listWorkLoads (
  cluster_name,
  type,
  search,
  keywords,
  pageNum,
  pageSize
) {
  let url
  let params = {}
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      url = appsV1Url(cluster_name, type)
      break
    case "pods":
      url = apiV1Url(cluster_name, type)
      break
    case "cronjobs":
      url = batchV1beta1Url(cluster_name, type)
      break
    case "jobs":
      url = batchV1Url(cluster_name, type)
      break
  }
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

export function getWorkLoadByName (cluster_name, type, namespace, name) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return get(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`
      )
    case "cronjobs":
      return get(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`
      )
    case "pods":
      return get(`${apiV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`)
    case "jobs":
      return get(`${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`)
  }
}

export function deleteWorkLoad (cluster_name, type, namespace, name) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return del(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`
      )
    case "cronjobs":
      return del(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`
      )
    case "pods":
      return del(`${apiV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`)
    case "jobs":
      return del(`${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`)
  }
}

export function createWorkLoad (cluster_name, type, namespace, data) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return post(`${appsV1UrlWithNsUrl(cluster_name, type, namespace)}`, data)
    case "cronjobs":
      return post(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}`,
        data
      )
    case "pods":
    case "services":
      return post(`${apiV1UrlWithNsUrl(cluster_name, type, namespace)}`, data)
    case "jobs":
      return post(`${batchV1WithNsUrl(cluster_name, type, namespace)}`, data)
  }
}

export function updateWorkLoad (cluster_name, type, namespace, name, data) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return put(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      )
    case "cronjobs":
      return put(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      )
    case "pods":
      return put(`${apiV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`)
    case "jobs":
      return put(
        `${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      )
  }
}
