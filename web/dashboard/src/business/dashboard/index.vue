<template>
  <layout-content :header="$t('business.dashboard.dashboard')">
    <el-row :gutter="24" class="row-box">
      <el-col :span="24">
        <el-card v-if="cluster" class="el-card" shadow="always" style="background-color: #243441;height: 120px">
          <el-row :gutter="24">
            <el-col :span="8">
              <span class="title">{{ $t("commons.table.name") }}</span>
              <div style="text-align: center">
                <h1>{{ cluster.name }}</h1>
              </div>
            </el-col>
            <el-col :span="8">
              <span class="title">{{ $t("business.cluster.version") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>{{ cluster.status.version }}</h1>
              </div>
            </el-col>
            <el-col :span="8">
              <span class="title">{{ $t("commons.table.created_time") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>{{ cluster.createAt | age }} </h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <br>

    <el-row :gutter="20" v-if="hasMetric === 'true'">
      <el-col :span="12">
        <el-card style="background-color: #212e38" class="n-card el-card">
          <span>CPU(core) {{clusterInfo.metricCpu}} / {{clusterInfo.allocatCpu}}</span>
          <el-progress style="margin-top: 20px" :stroke-width="20" :percentage="clusterInfo.cpuPercent"></el-progress>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card style="background-color: #212e38" class="n-card el-card">
          <span>Memory(Gi) {{clusterInfo.metricMemory}} / {{clusterInfo.allocatMemory}}</span>
          <el-progress style="margin-top: 20px" :stroke-width="20" :percentage="clusterInfo.memoryPercent"></el-progress>
        </el-card>
      </el-col>
    </el-row>
    <el-row v-if="hasMetric !== 'true'">
      <el-alert type="info" :closable="false">
        <el-button type="text" style="font-size: 15px" @click="dialogMetricVisible = true" icon="el-icon-warning">{{ $t("business.dashboard.metric_server_help") }}</el-button>
      </el-alert> 
    </el-row>
    <br>

    <el-row :gutter="20" class="resources row-box">
      <el-col v-for="(resource,index) in resources" v-bind:key="resource.name" :xs="8" :sm="8" :lg="6">
        <el-card :body-style="{padding: '0px'}" @click.native="jumpTo(resource.name)" class="d-card el-card">
          <el-row :gutter="24">
            <el-col :span="10">
              <div>
                <ko-charts :chart-data="resource" :key="resource.name" :index="index"></ko-charts>
              </div>
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>{{ resource.name }}</span>
                <h1 style="text-align: right"><a>{{ resource.count }}</a></h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="24" v-has-permissions="{apiGroup:'',resource:'events',verb:'list'}">
      <h4 style="margin-left: 10px;float: left">{{$t('business.event.event')}}</h4>
      <complex-table style="margin-top:20px" :data="events" @search="search" v-loading="loading" :pagination-config="paginationConfig"
                     :search-config="searchConfig">
        <el-table-column :label="$t('business.event.reason')" prop="reason" fix max-width="50px">
          <template v-slot:default="{row}">
            {{ row.reason }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" fix max-width="50px">
          <template v-slot:default="{row}">
            {{ row.metadata.namespace }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.event.message')" prop="resource" fix min-width="200px" show-overflow-tooltip>
          <template v-slot:default="{row}">
            {{ row.message }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.event.resource')" prop="resource" fix min-width="200px" show-overflow-tooltip>
          <template v-slot:default="{row}">
            <span class="span-link" @click="toResource(row.involvedObject.kind,row.metadata.namespace,row.involvedObject.name)">
              {{ row.involvedObject.kind }} / {{ row.involvedObject.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('commons.table.time')" prop="metadata.creationTimestamp" fix>
          <template v-slot:default="{row}">
            <span v-if="row.eventTime">{{ row.eventTime | age }}</span>
            <span v-else-if="row.lastTimestamp">{{ row.lastTimestamp | age }}</span>
            <span v-else-if="row.firstTimestamp">{{ row.firstTimestamp | age }}</span>
          </template>
        </el-table-column>
      </complex-table>
    </el-row>
    <div v-if="dialogMetricVisible">
      <metric-server @changeVisble="changeVisble" :clusterName="clusterName" :visible="dialogMetricVisible" />
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoCharts from "@/components/ko-charts"
import {listNamespace} from "@/api/namespaces"
import {listDeployments} from "@/api/deployments"
import {listStatefulSets} from "@/api/statefulsets"
import {listDaemonSets} from "@/api/daemonsets"
import {listServices} from "@/api/services"
import {listNodes} from "@/api/nodes"
import ComplexTable from "@/components/complex-table"
import {listEvents} from "@/api/events"
import {getCluster} from "@/api/clusters"
import {checkPermissions} from "@/utils/permission"
import {mixin} from "@/utils/resourceRoutes"
import {listConfigMaps} from "@/api/configmaps"
import {listSecrets} from "@/api/secrets"
import {listNodeMetrics} from "@/api/apis"
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"
import MetricServer from '@/components/ko-plugin/metric-server'

export default {
  name: "Dashboard",
  components: {ComplexTable, KoCharts, MetricServer, LayoutContent},
  mixins: [mixin],
  data() {
    return {
      cluster: null,
      clusterName: "",
      resources: [],
      events: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      hasMetric: "false",
      clusterInfo: {
        allocatCpu: 0,
        allocatMemory: 0,
        cpuPercent: 0,
        metricCpu: 0,
        metricMemory: 0,
        metricPercent: 0,
      },
      dialogMetricVisible: false,
      loading: false
    }
  },
  methods: {
    formatCpu(percentage) {
      return `${percentage}% \nCPU`
    },
    formatMemory(percentage) {
      return `${percentage}% \nMemory`
    },
    changeVisble(val) {
      this.dialogMetricVisible = val
    },
    jumpTo(val) {
      this.$router.push({name: val})
    },
    listResources() {
      getCluster(this.clusterName).then(res => {
        this.cluster = res.data
      })
      if (checkPermissions({scope: 'cluster', apiGroup: "", resource: "nodes", verb: "list"})) {
        listNodes(this.clusterName).then(res => {
          const nodes = {
            name: "Nodes",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.clusterInfo.allocatCpu = 0
          this.clusterInfo.allocatMemory = 0
          for(const n of res.items) {
            if (n.status?.allocatable?.cpu) {
              this.clusterInfo.allocatCpu += cpuUnitConvert(n.status.allocatable.cpu)
            }
            if (n.status?.allocatable?.memory) {
              this.clusterInfo.allocatMemory += memoryUnitConvert(n.status.allocatable.memory)
            }
          }
          this.resources.push(nodes)
          listNodeMetrics(this.clusterName).then(metricNodes => {
            this.hasMetric = "true"
            this.clusterInfo.metricCpu = 0
            this.clusterInfo.metricMemory = 0
            for(const n of metricNodes.items) {
              if (n.usage?.cpu) {
                this.clusterInfo.metricCpu += cpuUnitConvert(n.usage.cpu)
              }
              if (n.usage?.memory) {
                this.clusterInfo.metricMemory += memoryUnitConvert(n.usage.memory)
              }
            }
            this.clusterInfo.allocatCpu = Number((this.clusterInfo.allocatCpu / 1000).toFixed(2))
            this.clusterInfo.allocatMemory = Number((this.clusterInfo.allocatMemory / 1024).toFixed(2))
            this.clusterInfo.metricCpu = Number((this.clusterInfo.metricCpu / 1000).toFixed(2))
            this.clusterInfo.metricMemory = Number((this.clusterInfo.metricMemory / 1024).toFixed(2))
            this.clusterInfo.cpuPercent = Math.round((this.clusterInfo.metricCpu / this.clusterInfo.allocatCpu).toFixed(2) * 100)
            this.clusterInfo.memoryPercent = Math.round((this.clusterInfo.metricMemory / this.clusterInfo.allocatMemory).toFixed(2) * 100)
          }).catch(() => {
            this.hasMetric = "false"
          })
        })
      }
      if (checkPermissions({scope: 'cluster', apiGroup: "", resource: "namespaces", verb: "list"})) {
        listNamespace(this.clusterName).then(res => {
          const namespaces = {
            name: "Namespaces",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(namespaces)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "apps", resource: "deployments", verb: "list"})) {
        listDeployments(this.clusterName).then(res => {
          const deployments = {
            name: "Deployments",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(deployments)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "apps", resource: "statefulsets", verb: "list"})) {
        listStatefulSets(this.clusterName).then(res => {
          const statefulSets = {
            name: "StatefulSets",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(statefulSets)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "apps", resource: "daemonsets", verb: "list"})) {
        listDaemonSets(this.clusterName).then(res => {
          const daemonSets = {
            name: "DaemonSets",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(daemonSets)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "", resource: "services", verb: "list"})) {
        listServices(this.clusterName).then(res => {
          const services = {
            name: "Services",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(services)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "", resource: "configmaps", verb: "list"})) {
        listConfigMaps(this.clusterName).then(res => {
          const services = {
            name: "ConfigMaps",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(services)
        })
      }
      if (checkPermissions({scope: "namespace", apiGroup: "", resource: "secrets", verb: "list"})) {
        listSecrets(this.clusterName).then(res => {
          const services = {
            name: "Secrets",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.resources.push(services)
        })
      }
      this.search()
    },
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if (checkPermissions({scope: "namespace", apiGroup: "", resource: "events", verb: "list"})) {
        listEvents(this.clusterName, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
          this.loading = false
          this.events = res.items
          this.paginationConfig.total = res.total
        })
      }
    },
    getData(items, keys) {
      let key = []
      let result = []
      for (const item of items) {
        const name = this.traverse(item, keys)
        if (key.indexOf(name) === -1) {
          key.push(name)
          const d = {
            value: 1,
            name: name
          }
          result.push(d)
        } else {
          for (let i = 0; i < result.length; i++) {
            if (result[i].name === name) {
              result[i].value++
              break
            }
          }
        }
      }
      if (items.length === 0) {
        result = [{
          value: 0
        }]
      }
      return result
    },
    traverse(obj, keys) {
      if (keys === "status.conditions.type") {
        if (obj.status.conditions[0].type && obj.status.conditions[0].status === "True") {
          return obj.status.conditions[0].type
        } else {
          return "Error"
        }
      } else if (keys === "status.replicas") {
        return "In Progress"
      } else if (keys === "metadata.status.active" || keys === "status.loadBalancer") {
        return "Active"
      } else if (keys === "status.numberUnavailable") {
        if (obj.status.numberUnavailable > 0) {
          return "Error"
        } else {
          return "Active"
        }
      } else {
        return keys.split(".").reduce(function (cur, key) {
          return cur[key]
        }, obj)
      }
    }
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.listResources()
  }
}
</script>

<style scoped>
.line {
  float: left;
  margin-top: 10px;
  width: 1px;
  height: 80px;
  background: #3884c5;
}

.title {
  margin-left: 10px;
  color: #a1a9ae;
}

.d-card {
  height: 90px;
  background-color: #1d3e4d;
  margin-top: 10px;
}

.n-card {
  height: 100px;
  background-color: #1d3e4d;
  margin-top: 10px;
  border: none;
}

.card-content {
  margin-top: 10px;
  margin-right: 10px;
  text-align: center;
  float: right;
}

.card-content > span:first-child {
  color: #a1a9ae;
}
</style>
