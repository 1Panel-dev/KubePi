import { get, post, del, put } from "@/plugins/request";

const jobUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/jobs`;
};
const jobUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/namespaces/${namespace}/jobs`;
};

export function listJobs(cluster_name,  search) {
  let url = jobUrl(cluster_name);
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search;
  }
  return get(url);
}

export function listJobsWithNsSelector(cluster_name, namespace, selectors) {
  let url = jobUrlWithNs(cluster_name, namespace);
  const param = {};
  if (selectors && selectors !== "") {
    param.labelSelector = selectors;
  }
  return get(url, param);
}

export function getJobByName(cluster_name, namespace, job) {
  return get(`${jobUrlWithNs(cluster_name, namespace)}/${job}`);
}

export function deleteJob(cluster_name, job) {
  return del(`${jobUrl(cluster_name)}/${job}`);
}

export function createJob(cluster_name, job) {
  return post(`${jobUrl(cluster_name)}/${job}`);
}

export function updateJob(cluster_name, job) {
  return put(`${jobUrlWithNs(cluster_name)}/${job}`);
}
