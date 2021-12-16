import {del, get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/imagerepos"

export function searchRepos (pageNum, pageSize, conditions) {
  let url = `${baseUrl}/search?pageNum=${pageNum}&&pageSize=${pageSize}`
  return post(url, { conditions: conditions })
}

export function createRepo (data) {
  return post(baseUrl, data)
}

export function listInternalRepos (data) {
  return post(baseUrl + "/repositories", data)
}

export function deleteRepo (name) {
  return del(`${baseUrl}/${name}`)
}

export function getRepo (name) {
  return get(`${baseUrl}/${name}`)
}

export function updateRepo (name, data) {
  return put(`${baseUrl}/${name}`, data)
}

export function listRepoByCluster(cluster) {
  return get(`${baseUrl}/cluster/${cluster}`)
}

export function listImages(cluster,repo) {
  return get(`${baseUrl}/images/${cluster}/${repo}`)
}

export function listImagesByRepo(repo){
  return get(`${baseUrl}/images/${repo}`)
}
