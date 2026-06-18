import {get, post, put, del} from "@/plugins/request"

const baseUrl = "/api/v1/monitor"

// 仪表盘
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

// 指标
export function searchMetrics (pageNum, pageSize, conditions) {
    let url = `${baseUrl}/metrics/search?pageNum=${pageNum}&&pageSize=${pageSize}`
    return post(url, { conditions: conditions })
}

export function deleteMetrics (name) {
    return del(`${baseUrl}/metrics/${name}`)
}

export function createMetrics (data) {
    return post(`${baseUrl}/metrics`, data)
}

export function getMetrics (name) {
    return get(`${baseUrl}/metrics/${name}`)
}

export function updateMetrics (name, data) {
    return put(`${baseUrl}/metrics/${name}`, data)
}

export function getExplorer (name) {
    return get(`${baseUrl}/metrics/${name}/explorer`)
}

export function testConnectMetrics(name) {
    return get(`${baseUrl}/metrics/${name}/test/connect`)
}

export function queryMetrics(name, query, time) {
    return get(`${baseUrl}/metrics/${name}/query`, {
      promql: query,
      time: time
    })
}