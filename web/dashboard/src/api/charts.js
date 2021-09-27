import {del, get, post, put} from "@/plugins/request"

const baseUrl = "/api/v1/charts"

export function createChart(data) {
  return post(baseUrl, data)
}

export function searchCharts(cluster,page, size, keywords) {
  let url = `${baseUrl}/search?pageNum=${page}&pageSize=${size}&showExtra=true&cluster=${cluster}`
  if (keywords) {
    url = `${url}&keywords=${keywords}`
  }
  return post(url)
}

export function deleteChart(name) {
  return del(`${baseUrl}/${name}`)
}

export function getChart(name) {
  return get(`${baseUrl}/${name}`)
}

export function updateChart(name, user) {
  return put(`${baseUrl}/${name}`, user)
}
