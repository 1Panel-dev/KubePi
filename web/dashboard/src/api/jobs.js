import {get} from "@/plugins/request"


const jobUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/jobs`
}
const jobUrlWithNs = (cluster_name, namespace) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/namespaces/${namespace}/jobs`
}

export function listJobs (cluster_name, limit, continueToken, search) {
  let url = jobUrl(cluster_name)
  if (limit) {
    url += "?limit=" + limit
  }
  if (continueToken) {
    url += "&continue=" + continueToken
  }
  if (search && search !== "") {
    url += "&fieldSelector=metadata.name=" + search
  }
  return get(url)
}

//https://127.0.0.1:8443/apis/batch/v1/namespaces/kube-operator/jobs?labelSelector=apprepositories.kubeapps.com%2Frepo-name%3Dkubeoperator%2Capprepositories.kubeapps.com%2Frepo-namespace%3Dkube-operator&limit=500
export function listJobsWithNsSelector (cluster_name, namespace, selectors) {
  let url = jobUrlWithNs(cluster_name, namespace)
  const param = {}
  if (selectors && selectors !== "") {
    param.labelSelector = selectors
  }
  return get(url,param)
}
