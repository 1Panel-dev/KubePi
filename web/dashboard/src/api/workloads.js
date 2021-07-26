import { get, del, post, put } from "@/plugins/request";

const appsV1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/${type}`;
};
const appsV1UrlWithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespaces}/${type}`;
};
const batchV1beta1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/${type}`;
};
const batchV1beta1WithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1beta1/namespaces/${namespaces}/${type}`;
};
const batchV1Url = (cluster_name, type) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/${type}`;
};
const batchV1WithNsUrl = (cluster_name, type, namespaces) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/namespaces/${namespaces}/${type}`;
};

export function listWorkLoads(
  cluster_name,
  type,
  limit,
  continueToken,
  search
) {
  let url;
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      url = appsV1Url(cluster_name, type);
      break;
    case "cronjobs":
    case "pods":
      url = batchV1beta1Url(cluster_name, type);
      break;
    case "jobs":
      url = batchV1Url(cluster_name, type);
      break;
  }
  if (limit) {
    url += "?limit=" + limit;
  }
  if (continueToken) {
    url += "&continue=" + continueToken;
  }
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search;
  }
  return get(url);
}

export function listWorkLoadsByNs(cluster_name, type, namespace) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return get(`${appsV1UrlWithNsUrl(cluster_name, type, namespace)}`);
    case "cronjobs":
    case "pods":
      return get(`${batchV1beta1WithNsUrl(cluster_name, type, namespace)}`);
    case "jobs":
      return get(`${batchV1WithNsUrl(cluster_name, type, namespace)}`);
  }
}

export function getWorkLoadByName(cluster_name, type, namespace, name) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return get(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`
      );
    case "cronjobs":
    case "pods":
      return get(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`
      );
    case "jobs":
      return get(`${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`);
  }
}

export function deleteWorkLoad(cluster_name, type, namespace, name) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return get(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`
      );
    case "cronjobs":
    case "pods":
      return del(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`
      );
    case "jobs":
      return del(`${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`);
  }
}

export function createWorkLoad(cluster_name, type, namespace, data) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return post(`${appsV1UrlWithNsUrl(cluster_name, type, namespace)}`, data);
    case "cronjobs":
    case "pods":
      return post(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}`,
        data
      );
    case "jobs":
      return post(`${batchV1WithNsUrl(cluster_name, type, namespace)}`, data);
  }
}

export function updateWorkLoad(cluster_name, type, namespace, name, data) {
  switch (type) {
    case "deployments":
    case "statefulsets":
    case "daemonsets":
      return put(
        `${appsV1UrlWithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      );
    case "cronjobs":
    case "pods":
      return put(
        `${batchV1beta1WithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      );
    case "jobs":
      return put(
        `${batchV1WithNsUrl(cluster_name, type, namespace)}/${name}`,
        data
      );
  }
}
