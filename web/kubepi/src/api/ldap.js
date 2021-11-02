import {get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/ldap"

export function getLdap () {
  return get(`${baseUrl}`)
}

export function createLdap (data) {
  return post(`${baseUrl}`, data)
}

export function updateLdap(data){
  return put(`${baseUrl}`, data)
}

export function syncLdap(id,data) {
  return post(`${baseUrl}/sync/${id}`, data)

}
