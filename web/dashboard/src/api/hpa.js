import {get,} from "@/plugins/request"

const hpaUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/autoscaling/v1/horizontalpodautoscalers`
}
//
// const namespaceResourceQuotaUrl = (cluster_name, namespace) => {
//   return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/resourcequotas`
// }

export function listHpas (cluster_name, limit, continueToken, search) {
  let url = hpaUrl(cluster_name)
  const param = {}
  if (limit && limit !== 0) {
    param.limit = limit
  }
  if (continueToken && continueToken !== "") {
    param.continue = continueToken
  }
  if (search && search !== "") {
    param.fieldSelector = "metadata.namespace=" + search
  }
  return get(url, param)
}
