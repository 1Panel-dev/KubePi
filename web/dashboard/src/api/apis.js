import {get} from "@/plugins/request"

const apiUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis`
}

export function listApis(cluster_name) {
  return get(`${apiUrl(cluster_name)}`)
}


export function listResourceApi(cluster_name,apiGroup,version) {
  return get(`${apiUrl(cluster_name)}/${apiGroup}/${version}`)
}
