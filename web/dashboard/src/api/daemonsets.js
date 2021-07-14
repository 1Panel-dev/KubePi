import { get, del, post, put } from "@/plugins/request";

const daemonsetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/daemonsets`;
};
const daemonsetUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/daemonsets`;
};

export function listDaemonSets(cluster_name, limit, continueToken, search) {
  let url = daemonsetUrl(cluster_name);
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

export function getDaemonSetByName(cluster_name, namespace, daemonset) {
  return get(`${daemonsetUrlWithNs(cluster_name, namespace)}/${daemonset}`);
}

export function deleteDaemonSet(cluster_name, daemonset) {
  return del(`${daemonsetUrl(cluster_name)}/${daemonset}`);
}

export function createDaemonSet(cluster_name, daemonset) {
  return post(`${daemonsetUrl(cluster_name)}/${daemonset}`);
}

export function updateDaemonSet(cluster_name, daemonset) {
  return put(`${daemonsetUrlWithNs(cluster_name)}/${daemonset}`);
}
