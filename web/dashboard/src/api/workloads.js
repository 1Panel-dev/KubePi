import { get } from "@/plugins/request";

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

