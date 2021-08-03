import {get, del, post, put} from "@/plugins/request"

const deploymentUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/deployments`
}
const deploymentWithNsUrl = (cluster_name, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespaces}/deployments`
}

export function listDeployments (cluster_name, currentPage, pageSize) {
  let url = deploymentUrl(cluster_name)
  if (currentPage && pageSize) {
    return get(`${url}?pageNum=${currentPage}&pageSize=${pageSize}`)
  }
  return get(url)
}

export function listDeploymentsByNs (cluster_name, namespace) {
  return get(`${deploymentWithNsUrl(cluster_name, namespace)}`)
}

export function getDeploymentByName (cluster_name, namespace, deployment) {
  return get(`${deploymentWithNsUrl(cluster_name, namespace)}/${deployment}`)
}

export function deleteDeployment (cluster_name, deployment) {
  return del(`${deploymentUrl(cluster_name)}/${deployment}`)
}

export function createDeployment (cluster_name, deployment) {
  return post(`${deploymentUrl(cluster_name)}/${deployment}`)
}

export function updateDeployment (cluster_name, deployment) {
  return put(`${deploymentWithNsUrl(cluster_name)}/${deployment}`)
}
