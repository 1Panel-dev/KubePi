import {post, get} from "@/plugins/request"


const mfa = "/api/v1/mfa"

export function getOtp() {
  return get(`${mfa}/`)
}

export function bind(data, config){
  return post(`${mfa}/bind`, data, undefined, config)
}

export function valid(data, config) {
  return post(`${mfa}/valid`, data, undefined, config)
}
