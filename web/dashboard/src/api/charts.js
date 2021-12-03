import {del, get, post, put} from "@/plugins/request"


const chartUrl = (cluster) => {
  return `/api/v1/charts/${cluster}`
}

const repoUrl = (cluster) => {
  return `/api/v1/charts/${cluster}/repos`
}

const appUrl = (cluster) => {
  return `/api/v1/apps/${cluster}`
}

export function createChart (cluster, data) {
  return post(chartUrl(cluster), data)
}

export function deleteChart (cluster, name) {
  return del(`${chartUrl(cluster)}/${name}`)
}

export function getChart (cluster, repo, name) {
  const params = {}
  params["repo"] = repo
  return get(`${chartUrl(cluster)}/${name}`, params)
}

export function updateChart (cluster, name, chart) {
  return put(`${chartUrl(cluster)}/${name}`, chart)
}

export function listRepos (cluster) {
  return get(repoUrl(cluster))
}

export function createRepo (cluster, data) {
  return post(`${repoUrl(cluster)}`, data)
}

export function deleteRepo (cluster, name) {
  return del(`${repoUrl(cluster)}/${name}`)
}

export function updateRepo(cluster, name,data) {
  return put(`${repoUrl(cluster)}/${name}`, data)
}

export function syncRepo(cluster, name) {
  return post(`${repoUrl(cluster)}/sync/${name}`)
}

export function getRepo (cluster, name) {
  return get(`${repoUrl(cluster)}/${name}`)
}

export function searchCharts (cluster, repo, page, size, keywords) {

  const url = chartUrl(cluster)
  const params = {}
  params["repo"] = repo

  if (keywords) {
    params["pattern"] = keywords
  }
  if (page && size) {
    params["pageNum"] = page
    params["pageSize"] = size
  }
  return get(url + "/search", params)
}

export function getChartByVersion (cluster, repo, name, version) {
  const params = {}
  params["repo"] = repo
  params["version"] = version
  return get(`${chartUrl(cluster)}/detail/${name}`, params)
}
export function installChart (cluster, data) {
  return post(`${chartUrl(cluster)}/install`, data)
}

export function searchInstalled (cluster, page, size, keywords) {

  const params = {}
  if (keywords) {
    params["pattern"] = keywords
  }
  if (page && size) {
    params["pageNum"] = page
    params["pageSize"] = size
  }
  return get(`${appUrl(cluster)}/search`, params)
}

export function deleteApp (cluster,namespace, name) {
  return del(`${appUrl(cluster)}/${namespace}/${name}`)
}

export function getApp (cluster, name) {
  return get(`${appUrl(cluster)}/${name}`)
}

export function getChartUpdate (cluster, chart, name) {
  const params = {}
  params["chart"] = chart
  return get(`${appUrl(cluster)}/update/${name}`,params)
}

export function upgradeChart (cluster, name, data) {
  return put(`${appUrl(cluster)}/${name}`, data)
}





