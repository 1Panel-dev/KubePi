import { get, del, post } from "@/plugins/request";

const deploymentUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/apis/apps/v1/deployments`;
};
const deploymentWithNsUrl = (cluster_name, namespaces) => {
return `/proxy/${cluster_name}/apis/apps/v1/namespaces/${namespaces}/deployments`;
};

export function listDeployments(cluster_name) {
  return get(`${deploymentUrl(cluster_name)}`);
}

export function getDeploymentByName(cluster_name, namespace, deployment) {
  return get(`${deploymentWithNsUrl(cluster_name, namespace)}/${deployment}`);
}

export function deleteDeployment(cluster_name, deployment) {
  return del(`${deploymentUrl(cluster_name)}/${deployment}`);
}

export function createDeployment(cluster_name, deployment) {
  return post(`${deploymentUrl(cluster_name)}/${deployment}`);
}