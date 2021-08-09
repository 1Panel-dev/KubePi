import {del, get, patch, post} from "@/plugins/request"

const namespaceUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces`
}


export function listNamespace(cluster_name, search, keywords, pageNum, pageSize) {
    let url = namespaceUrl(cluster_name)
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

export function createNamespace(cluster_name, data) {
    return post(`${namespaceUrl(cluster_name)}`, data)
}

export function getNamespace(cluster_name, namespace) {
    return get(`${namespaceUrl(cluster_name)}/${namespace}`)
}

export function deleteNamespace(cluster_name, namespace) {
    return del(`${namespaceUrl(cluster_name)}/${namespace}`)
}

export function updateNamespace(cluster_name, namespace, data) {
    return patch(`${namespaceUrl(cluster_name)}/${namespace}`, data)
}
