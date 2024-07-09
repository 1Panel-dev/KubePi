import {get} from "@/plugins/request"
import qs from "qs";
const prometheusUrl = (cluster_name) => {
  return `/api/v1/proxy/${cluster_name}/k8s/apis/monitoring.coreos.com/v1/prometheuses`
}
const prometheusMetricsUrl = (cluster_name,namespace,prometheus_service_name,prometheus_prefix,prometheus_service_port) => {
  return `/api/v1/proxy/${cluster_name}/k8s/api/v1/namespaces/${namespace}/services/${prometheus_service_name}:${prometheus_service_port}/proxy${prometheus_prefix}api/v1/query_range`
}
/*获取已部署的prometheus服务器*/
export function listPrometheuses (cluster_name) {
    let url = prometheusUrl(cluster_name)
    const param = {}
    return get(url, param)
}

/*获取容器指标（mem and fs）*/
export function getContainerMetricsMemAndFs (cluster_name,namespace,prometheus_service_name,prometheus_prefix,prometheus_service_port,pod_namespace,pod_name,container_name,secs) {
  let url = prometheusMetricsUrl(cluster_name,namespace,prometheus_service_name,prometheus_prefix,prometheus_service_port)
  const startTimestamp = Math.floor(Date.now() / 1000) - secs; // Unix timestamp for 30 minutes ago
  const endTimestamp = Math.floor(Date.now() / 1000); // Current Unix timestamp
  const param = {
    "query": `{__name__=~"container_fs_usage_bytes|container_memory_working_set_bytes",namespace="${pod_namespace}",pod="${pod_name}",container="${container_name}"}`,
    "start": startTimestamp,
    "end": endTimestamp,
    "step": "1m"
  }
  return get(url+"?"+qs.stringify(param) , {})
}
/*获取容器指标（mem and fs）*/
export function getContainerMetricsCpu(cluster_name,namespace,prometheus_service_name,prometheus_prefix,prometheus_service_port,pod_namespace,pod_name,container_name,secs) {
  let url = prometheusMetricsUrl(cluster_name,namespace,prometheus_service_name,prometheus_prefix,prometheus_service_port)
  const startTimestamp = Math.floor(Date.now() / 1000) - secs; // Unix timestamp for 30 minutes ago
  const endTimestamp = Math.floor(Date.now() / 1000); // Current Unix timestamp
  const param = {
    "query": `sum(irate(container_cpu_usage_seconds_total{namespace="${pod_namespace}",pod="${pod_name}",container="${container_name}"}[5m]))`,
    "start": startTimestamp,
    "end": endTimestamp,
    "step": "1m"
  }
  return get(url+"?"+qs.stringify(param) , {})
}
