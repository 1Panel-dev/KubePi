import {post} from "@/plugins/request"


const apiUrl = (cluster_name, version, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/${version}/${kind}`
}

const apiNsUrl = (cluster_name, version, namespace, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/${version}/namespaces/${namespace}/${kind}`
}


const apisUrl = (cluster_name, group, version, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/${group}/${version}/${kind}`
}

const apisNsUrl = (cluster_name, version, group, namespace, kind) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/${group}/${version}/namespaces/${namespace}/${kind}`
}

const apiKinds = ["Namespace", "Pod","Service","Endpoints"]

export function postYaml (cluster_name, kind, data) {
  return post(getUrl(cluster_name, kind, data), data)
}

export function getUrl (cluster_name, kind, data) {
  let url = ""
  if (apiKinds.indexOf(data.kind) > -1) {
    if (data.metadata.namespace !== "") {
      url = apiNsUrl(cluster_name, data.apiVersion, data.metadata.namespace, kind)
    } else {
      url = apiUrl(cluster_name, data.apiVersion, kind)
    }
  } else {
    const apiVersions = data.apiVersion.split("/")
    const group = apiVersions[0]
    const version = apiVersions[1]
    if (data.metadata.namespace !== "") {
      url = apisNsUrl(cluster_name, version, group, data.metadata.namespace, kind)
    } else {
      url = apisUrl(cluster_name, version, group, kind)
    }
  }
  return url
}
