import {post,get} from "@/plugins/request"

const authUrl = "/api/v1/clusters"

export function create(data) {
  return post(authUrl, data)
}

export function listAll(){
  return get(authUrl)
}
