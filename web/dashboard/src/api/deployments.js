import {get, patch} from "@/plugins/request"

const deploymentUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/deployments`
}
const deploymentWithNsUrl = (cluster_name, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespaces}/deployments`
}

export function listDeployments (cluster_name, currentPage, pageSize) {
  let url = deploymentUrl(cluster_name)
  if (currentPage && pageSize) {
    let params = {pageNum: currentPage, pageSize: pageSize }
    return get(url, params)
  }
  return get(url)
}

export function listDeploymentsByNs (cluster_name, namespace) {
  return get(`${deploymentWithNsUrl(cluster_name, namespace)}`)
}

export function scaleDeployment (cluster_name, namespace, deployment, data) {
  return patch(`${deploymentWithNsUrl(cluster_name, namespace)}/${deployment}/scale`, data)
}

export function patchDeployment (cluster_name, namespace, deployment, data) {
  return patch(`${deploymentWithNsUrl(cluster_name, namespace)}/${deployment}`, data)
}