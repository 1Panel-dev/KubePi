import { get, patch } from "@/plugins/request";

const daemonsetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/daemonsets`;
};
const daemonsetUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/daemonsets`;
};

export function listDaemonSets(cluster_name, search) {
  let url = daemonsetUrl(cluster_name);
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search;
  }
  return get(url);
}

export function patchDaemonset(cluster_name, namespace, daemonset, data) {
  return patch(`${daemonsetUrlWithNs(cluster_name, namespace)}/${daemonset}`, data);
}

