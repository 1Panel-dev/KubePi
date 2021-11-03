import {listResourceApi} from "@/api/apis"

export async function checkApi (cluster_name, api_group, version, resource) {
  let result = false
  let res = await listResourceApi(cluster_name, api_group, version)
  res.resources.forEach(function (item) {
    if (item.kind === resource) {
      result = true
    }
  })
  return result
}


