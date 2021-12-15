import {post} from "@/plugins/request"


const apiUrl = (cluster_name, group, version, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/${group}/${version}/${kind}`
}

const apisUrl = (cluster_name, group, version, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/${group}/${version}/${kind}`
}

const apiKinds = ["namespaces"]

export function postYaml (cluster_name, group, version, kind, data) {
  let url = apisUrl(cluster_name, group, version, kind)
  if (apiKinds.indexOf(kind) > -1) {
    url = apiUrl(cluster_name, group, version, kind)
  }
  return post(url, data)
}
