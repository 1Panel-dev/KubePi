import {get} from "@/plugins/request"

const podUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/pods`
}

export function listPods (cluster_name, limit, continueToken, search) {
  let url = podUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector =  "spec.nodeName=" + search
  }
  return get(url,param)
}


