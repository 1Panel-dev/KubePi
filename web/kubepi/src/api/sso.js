import {get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/sso"

export function getSso () {
    return get(`${baseUrl}`)
}

export function createSso (data) {
    return post(`${baseUrl}`, data)
}

export function updateSso(data){
    return put(`${baseUrl}`, data)
}

export function isSso() {
    return get(`${baseUrl}/status`)
}

export function ssoLogin() {
    return (`${baseUrl}/login`)
}

export function testConnect(data) {
    return post(`${baseUrl}/test/connect`, data)
}