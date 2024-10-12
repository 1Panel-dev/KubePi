import {get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/monitor"

export function getGrafana () {
    return get(`${baseUrl}/grafana`)
}

export function createGrafana (data) {
    return post(`${baseUrl}/grafana`, data)
}

export function updateGrafana(data){
    return put(`${baseUrl}/grafana`, data)
}

export function importGrafanaDashboard(data) {
    return post(`${baseUrl}/grafana/import`, data)
}

export function testConnectGrafana(data) {
    return post(`${baseUrl}/grafana/test/connect`, data)
}