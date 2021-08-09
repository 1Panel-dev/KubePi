import {get} from "@/plugins/request"

const nodeUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/api/v1/nodes`
}


export function listNodes(cluster_name, search, keywords, pageNum, pageSize) {
    let url = nodeUrl(cluster_name)
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

export function getNode(cluster_name, name) {
    return get(`${nodeUrl(cluster_name)}/${name}`)
}