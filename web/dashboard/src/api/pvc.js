import {del, get, patch, post} from "@/plugins/request"

const serviceUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumeclaims`
}
const namespaceSecretUrl = (cluster_name, namespace) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/persistentvolumeclaims`
}

export function listPvcs (cluster_name, limit, continueToken, search) {
    let url = serviceUrl(cluster_name)
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

export function listPvcsWithNs (cluster_name, namespace) {
    return get(`${namespaceSecretUrl(cluster_name, namespace)}`)
}

export function deletePvcs (cluster_name, namespace, name) {
    return del(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`)
}

export function getSecret (cluster_name, namespace, name) {
    return get(`${namespaceSecretUrl(cluster_name, namespace)}/${name}`)
}

export function createSecret (cluster_name, namespace, data) {
    return post(`${namespaceSecretUrl(cluster_name, namespace)}`, data)
}

export function editSecret (cluster_name, namespace, data) {
    return patch(`${namespaceSecretUrl(cluster_name, namespace)}`, data)
}

