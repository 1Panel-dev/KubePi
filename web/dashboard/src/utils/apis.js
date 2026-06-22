import {listApis, listResourceApi} from "@/api/apis"

export async function checkApi (cluster_name, api_group, version, resource) {
  let apiList = await listApis(cluster_name)
  const group = (apiList.groups || []).find(item => item.name === api_group)
  const versions = group ? (group.versions || []) : []
  if (!group || !versions.some(item => item.version === version || item.groupVersion === `${api_group}/${version}`)) {
    return false
  }
  let result = false
  let res
  try {
    res = await listResourceApi(cluster_name, api_group, version)
  } catch (e) {
    return false
  }
  res.resources.forEach(function (item) {
    if (item.kind === resource) {
      result = true
    }
  })
  return result
}
