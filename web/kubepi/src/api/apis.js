import {get} from "@/plugins/request"

const baseUrl = "/api/v1"

export function listApiResource() {
    return get(baseUrl)
}