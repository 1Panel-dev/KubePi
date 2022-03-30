import {post, get} from "@/plugins/request"


const mfa = "/api/v1/mfa"

export function getOtp() {
  return get(`${mfa}/`)
}

export function bind(data){
  return post( `${mfa}/bind`,data)
}

export function valid(data) {
  return post(`${mfa}/valid`,data)
}
