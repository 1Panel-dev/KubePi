import {get} from "@/plugins/request"


const jobUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/batch/v1/jobs`
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