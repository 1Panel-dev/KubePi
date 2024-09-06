import { get, patch } from "@/plugins/request";

const statefulsetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/statefulsets`;
};
const statefulsetUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/statefulsets`;
};
const statefulsetUrlWithNsAndName = (cluster_name, namespace ,name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/statefulsets/${name}`;
};

export function listStatefulSets(cluster_name,  search, keywords, pageNum, pageSize) {
  let url = statefulsetUrl(cluster_name);
  const params = {}
  if (search) {
    params["search"] = true
    if (keywords) {
      params["keywords"] = keywords
    }
    if (pageNum && pageSize) {
      params["pageNum"] = pageNum
      params["pageSize"] = pageSize
    }
  }
  return get(url,params);
}

export function getStatefulSet(cluster_name,  namespace ,name) {
  let url = statefulsetUrlWithNsAndName(cluster_name,  namespace ,name);
  const params = {}
  return get(url,params);
}

export function scaleStatefulset(cluster_name, namespace, statefulset, data) {
  return patch(`${statefulsetUrlWithNs(cluster_name, namespace)}/${statefulset}/scale`, data);
}

export function patchStatefulset(cluster_name, namespace, statefulset, data) {
  return patch(`${statefulsetUrlWithNs(cluster_name, namespace)}/${statefulset}`, data);
}