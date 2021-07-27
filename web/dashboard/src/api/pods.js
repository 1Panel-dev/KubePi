import { get, del } from "@/plugins/request";

const podUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/pods`;
};
const podUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/pods`;
};

export function listPods(cluster_name,  search) {
  let url = podUrl(cluster_name);
  const param = {};
  if (search && search !== "") {
    param.fieldSelector = "spec.nodeName=" + search;
  }
  return get(url, param);
}

export function listPodsWithNsSelector(cluster_name, namespace, selectors) {
  let url = podUrlWithNs(cluster_name, namespace);
  const param = {};
  if (selectors && selectors !== "") {
    param.labelSelector = selectors;
  }
  return get(url, param);
}

export function getPodByName(cluster_name, namespace, pod) {
  return get(`${podUrlWithNs(cluster_name, namespace)}/${pod}`);
}

export function deletePod(cluster_name, pod) {
  return del(`${podUrl(cluster_name)}/${pod}`);
}

export function getPodLogsByName(cluster_name, namespace, pod, paramInfo) {
  return get(`${podUrlWithNs(cluster_name, namespace)}/${pod}/log` + paramInfo);
}
