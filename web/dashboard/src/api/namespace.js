import {get} from "@/plugins/request"

const namespaceUrl = (cluster_name) => {
  return  `/proxy/${cluster_name}/api/v1/namespaces`
}


export function listNamespace(cluster_name) {
    return get(`${namespaceUrl(cluster_name)}`)
}
