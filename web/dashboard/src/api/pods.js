import { get, del } from "@/plugins/request";

const podUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/pods`;
};
const podUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/pods`;
};

export function listPods(cluster_name, limit, continueToken, search) {
  let url = podUrl(cluster_name);
  const param = {};
  if (limit && limit !== 0) {
    param.limit = limit;
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken;
  }
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

export function getPodByName(cluster_name, namespace, cornjob) {
  return get(`${podUrlWithNs(cluster_name, namespace)}/${cornjob}`);
}

export function deletePod(cluster_name, cornjob) {
  return del(`${podUrl(cluster_name)}/${cornjob}`);
}
