import {post, get, del} from "@/plugins/request"

const authUrl = "/api/v1/clusters"

export function create(data) {
  return post(authUrl, data)
}

export function listAll(){
  return get(authUrl)
}

export function deleteBy(name){
  return del(`${authUrl}/${name}`)
}

export function searchCluster(page,size) {
  return post(`${authUrl}/search?pageNum=${page}&pageSize=${size}`)
}
