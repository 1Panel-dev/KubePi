import {del, get, put, post} from "@/plugins/request"

const pvcUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/persistentvolumeclaims`
}
const namespacePvcUrl = (cluster_name, namespace) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/persistentvolumeclaims`
}

export function listPvcs (cluster_name, search, keywords, pageNum, pageSize) {
    let url = pvcUrl(cluster_name)
    const params = {}
    if (search) {
        params["search"] = true
        if (keywords) {
            params["keywords"] = keywords
        }
        if (pageNum && pageSize) {
            params["pageNum"] = pageNum
            params["pageSize"] = pageSize
        }
    }
    return get(url, params)
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

export function updatePvc (cluster_name, namespace, data, name) {
    return put(`${namespacePvcUrl(cluster_name, namespace)}/${name}`, data)
}
