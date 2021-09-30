import {post} from "@/plugins/request";

const baseUrl = "/api/v1/webkubectl/session"

export function getTerminalSession(clusterName) {
    return post(`${baseUrl}`, {cluster: clusterName})
}