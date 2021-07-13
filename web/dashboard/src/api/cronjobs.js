import { get, post, del } from "@/plugins/request";

const cornJobUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/cronjobs`;
};
const cornJobWithNsUrl = (cluster_name, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/namespaces/${namespaces}/cronjobs`;
};

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
