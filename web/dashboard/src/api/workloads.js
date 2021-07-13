import { get, del, post } from "@/plugins/request";

const deploymentUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/deployments`;
};
const deploymentWithNsUrl = (cluster_name, namespaces) => {
return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespaces}/deployments`;
};

const cornJobUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/cronjobs`;
};
const cornJobWithNsUrl = (cluster_name, namespaces) => {
return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/namespaces/${namespaces}/cronjobs`;
};

// deployment
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

// cornjob
export function listCronJobs(cluster_name) {
  return get(`${cornJobUrl(cluster_name)}`);
}

export function getCronJobByName(cluster_name, namespace, cornjob) {
  return get(`${cornJobWithNsUrl(cluster_name, namespace)}/${cornjob}`);
}

export function deleteCronJob(cluster_name, cornjob) {
  return del(`${cornJobUrl(cluster_name)}/${cornjob}`);
}

export function createCronJob(cluster_name, cornjob) {
  return post(`${cornJobUrl(cluster_name)}/${cornjob}`);
}