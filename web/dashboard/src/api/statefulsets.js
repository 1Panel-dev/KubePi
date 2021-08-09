import { get, del, post, put } from "@/plugins/request";

const statefulsetUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/statefulsets`;
};
const statefulsetUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/apps/v1/namespaces/${namespace}/statefulsets`;
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

export function getStatefulSetByName(cluster_name, namespace, statefulset) {
  return get(`${statefulsetUrlWithNs(cluster_name, namespace)}/${statefulset}`);
}

export function deleteStatefulSet(cluster_name, statefulset) {
  return del(`${statefulsetUrl(cluster_name)}/${statefulset}`);
}

export function createStatefulSet(cluster_name, statefulset) {
  return post(`${statefulsetUrl(cluster_name)}/${statefulset}`);
}

export function updateStatefulSet(cluster_name, statefulset) {
  return put(`${statefulsetUrlWithNs(cluster_name)}/${statefulset}`);
}
