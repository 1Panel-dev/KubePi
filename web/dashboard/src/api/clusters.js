import {get} from "@/plugins/request"

const authUrl = "/api/v1/clusters"

export function getCluster(name) {
    return get(`${authUrl}/${name}`)
}
