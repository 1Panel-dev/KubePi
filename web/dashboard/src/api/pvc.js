import {del, get, patch, post} from "@/plugins/request"

const pvcUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumeclaims`
}
const namespacePvcUrl = (cluster_name, namespace) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/persistentvolumeclaims`
}

export function listPvcs (cluster_name, continueToken, search) {
    let url = pvcUrl(cluster_name)
    const param = {}
    if (search && search !== "") {
        param.fieldSelector = "metadata.namespace=" + search
    }
    return get(url, param)
}

export function listPvcsWithNs (cluster_name, namespace) {
    return get(`${namespacePvcUrl(cluster_name, namespace)}`)
}

export function deletePvcs (cluster_name, namespace, name) {
    return del(`${namespacePvcUrl(cluster_name, namespace)}/${name}`)
}

export function getPvc (cluster_name, namespace, name) {
    return get(`${namespacePvcUrl(cluster_name, namespace)}/${name}`)
}

export function createPvc (cluster_name, namespace, data) {
    return post(`${namespacePvcUrl(cluster_name, namespace)}`, data)
}

export function editPvc (cluster_name, namespace, data) {
    return patch(`${namespacePvcUrl(cluster_name, namespace)}`, data)
}