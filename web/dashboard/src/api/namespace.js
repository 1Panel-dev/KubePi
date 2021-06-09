import {get} from "@/plugins/request"

const baseUrl = "/proxy/api/v1/namespaces"

export function listNamespace() {
    return get(baseUrl)
}
