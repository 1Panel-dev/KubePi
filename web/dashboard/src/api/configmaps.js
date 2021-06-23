
import {get} from "@/plugins/request"


const configMapUrl = (cluster_name) => {
  return `/proxy/${cluster_name}/api/v1/configmaps`
}

export function listConfigMaps (cluster_name, limit, continueToken, search) {
  let url = configMapUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector = "metadata.namespace="+search
  }
  return get(url,param)
}