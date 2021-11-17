import {post} from "@/plugins/request"

const baseUrl = "/api/v1/imagerepos"

export function searchRepos (pageNum, pageSize, conditions) {
  let url = `${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`
  return post(url, { conditions: conditions })
}

export function createRepo (data) {
  return post(baseUrl, data)
}

export function listInternalRepos(data) {
  return post(baseUrl+"/repositories", data)
}
