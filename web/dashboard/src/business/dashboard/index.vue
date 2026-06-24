<template>
  <layout-content :header="$t('business.dashboard.dashboard')">
    <el-row :gutter="24" class="row-box">
      <el-col :span="24">
        <el-card v-if="cluster" class="el-card" shadow="always" style="background-color: #f8fafc;height: 120px">
          <el-row :gutter="24">
            <el-col :span="8">
              <span class="title">{{ $t("commons.table.name") }}</span>
              <div style="text-align: center">
                <el-tooltip class="item" effect="dark" :content="cluster.name" placement="top">
                  <h1>{{ cluster.name | max15letter }}</h1>
                </el-tooltip>
              </div>
            </el-col>
            <el-col :span="8">
              <span class="title">{{ $t("business.cluster.version") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <el-tooltip class="item" effect="dark" :content="cluster.status.version" placement="top">
                  <h1>{{ cluster.status.version | max15letter}}</h1>
                </el-tooltip>
              </div>
            </el-col>
            <el-col :span="8">
              <span class="title">{{ $t("commons.table.created_time") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <el-tooltip class="item" effect="dark" :content="cluster.createAt | age" placement="top">
                  <h1>{{ cluster.createAt | age }} </h1>
                </el-tooltip>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <br>

    <el-row :gutter="20" v-if="showMetric && hasMetric === 'true'">
      <el-col :span="12">
        <el-card style="background-color: #f8fafc" class="n-card el-card">
          <span>CPU(core) {{clusterInfo.metricCpu}} / {{clusterInfo.allocatableCpu}}</span>
          <el-progress style="margin-top: 20px" :stroke-width="20" :percentage="clusterInfo.cpuPercent"></el-progress>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card style="background-color: #f8fafc" class="n-card el-card">
          <span>Memory(Gi) {{clusterInfo.metricMemory}} / {{clusterInfo.allocatableMemory}}</span>
          <el-progress style="margin-top: 20px" :stroke-width="20" :percentage="clusterInfo.memoryPercent"></el-progress>
        </el-card>
      </el-col>
    </el-row>
    <el-row v-if="showMetric && hasMetric === 'false'" class="metric-notice-row">
      <button class="metric-notice" type="button" @click="dialogMetricVisible = true">
        <i class="el-icon-info metric-notice__icon"></i>
        <span class="metric-notice__text">{{ $t("business.dashboard.metric_server_help") }}</span>
      </button>
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
                     :search-config="searchConfig"
                     :showFullTextSwitch="true" @update:isFullTextSearch="OnIsFullTextSearchChange">
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
      <metric-server @changeVisible="changeVisible" :clusterName="clusterName" :visible="dialogMetricVisible" />
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
import {listEventsWithNs} from "@/api/events"
import {getCluster} from "@/api/clusters"
import {checkPermissions} from "@/utils/permission"
import {mixin} from "@/utils/resourceRoutes"
import {listConfigMaps} from "@/api/configmaps"
import {listSecrets} from "@/api/secrets"
import {listNodeMetrics} from "@/api/apis"
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"
import MetricServer from '@/components/ko-plugin/metric-server'
import { searchFullTextItems } from "@/api/fulltextsearch/fulltextsearch"
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
      showMetric: true,
      hasMetric: "",
      clusterInfo: {
        allocatableCpu: 0,
        allocatableMemory: 0,
        cpuPercent: 0,
        metricCpu: 0,
        metricMemory: 0,
        metricPercent: 0,
      },
      dialogMetricVisible: false,
      loading: false,
      isFullTextSearch: false,
      intervalId:null,//定时器
    }
  },
  methods: {
    formatCpu(percentage) {
      return `${percentage}% \nCPU`
    },
    formatMemory(percentage) {
      return `${percentage}% \nMemory`
    },
    changeVisible(val) {
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
        this.resources=[]
        listNodes(this.clusterName).then(res => {
          const nodes = {
            name: "Nodes",
            count: res.items? res.items.length : 0,
            data: [{
              value: res.items? res.items.length : 0
            }]
          }
          this.clusterInfo.allocatableCpu = 0
          this.clusterInfo.allocatableMemory = 0
          for(const n of res.items) {
            if (n.status?.allocatable?.cpu) {
              this.clusterInfo.allocatableCpu += cpuUnitConvert(n.status.allocatable.cpu)
            }
            if (n.status?.allocatable?.memory) {
              this.clusterInfo.allocatableMemory += memoryUnitConvert(n.status.allocatable.memory)
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
            this.clusterInfo.allocatableCpu = Number((this.clusterInfo.allocatableCpu / 1000).toFixed(2))
            this.clusterInfo.allocatableMemory = Number((this.clusterInfo.allocatableMemory / 1024).toFixed(2))
            this.clusterInfo.metricCpu = Number((this.clusterInfo.metricCpu / 1000).toFixed(2))
            this.clusterInfo.metricMemory = Number((this.clusterInfo.metricMemory / 1024).toFixed(2))
            this.clusterInfo.cpuPercent = Math.round((this.clusterInfo.metricCpu / this.clusterInfo.allocatableCpu).toFixed(2) * 100)
            this.clusterInfo.memoryPercent = Math.round((this.clusterInfo.metricMemory / this.clusterInfo.allocatableMemory).toFixed(2) * 100)
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
      
    },
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if (checkPermissions({scope: "namespace", apiGroup: "", resource: "events", verb: "list"})) {
        const ns =sessionStorage.getItem("namespace")
        if(!this.isFullTextSearch){
          listEventsWithNs(this.clusterName, ns, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
            this.events = res.items
            this.paginationConfig.total = res.total
          }).catch(() => {
            this.events = []
            this.paginationConfig.total = 0
          }).finally(() => {
            this.loading = false
          })
        } else {
          listEventsWithNs(this.clusterName, ns,false)
          .then((res) => {
            const results = searchFullTextItems(res.items,this.searchConfig.keywords);
            this.events =results.slice(this.paginationConfig.currentPage*this.paginationConfig.pageSize-this.paginationConfig.pageSize,this.paginationConfig.currentPage*this.paginationConfig.pageSize)
            this.paginationConfig.total = results.length
          }).catch(() => {
            this.events = []
            this.paginationConfig.total = 0
          }).finally(() => {
            this.loading = false
          }) 
        }
      } else {
        this.loading = false
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
    },
    //改变选项"是否全文搜索"
    OnIsFullTextSearchChange(val){
      this.isFullTextSearch=val
    },
    stopTimeTick(){
      if (!this.intervalId) {
        return;
      }
      clearInterval(this.intervalId); //清除计时器
      this.intervalId = null; //设置为null
    },
    startTimeTick(tick){
      // 计时器为空，操作
      this.intervalId = setInterval(() => {
          this.stopTimeTick()
          //刷新图表
          this.listResources();
          console.log(tick)
          this.startTimeTick(tick)
      },tick*1000);
    }
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.showMetric = checkPermissions({scope: 'cluster', apiGroup: "", resource: "nodes", verb: "list"}) && checkPermissions({scope: 'cluster', apiGroup: "metrics.k8s.io", resource: "nodes", verb: "list"})
    this.listResources()
    this.search()
  },
  watch: {
    /*监控自动刷新变量*/
   "$store.state.app.autorefresh":{
    handler:function(newVal,oldVal){
      if(!newVal || newVal=='-1' || newVal=='undefined' || newVal==''){
        this.stopTimeTick();
      } else {
        this.stopTimeTick();
        // 计时器为空，操作
        this.startTimeTick(parseInt(newVal));
      }
    }
   }
  },
  destroyed(){
    // 在页面销毁后，清除计时器
    this.stopTimeTick();
  }
}
</script>

<style scoped>
.line {
  float: left;
  margin-top: 10px;
  width: 1px;
  height: 80px;
  background: #2563eb;
}

.title {
  margin-left: 10px;
  color: #64748b;
}

.d-card {
  height: 90px;
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  margin-top: 10px;
}

.n-card {
  height: 100px;
  background-color: #f8fafc;
  margin-top: 10px;
  border: none;
}

.metric-notice-row {
  margin-top: 12px;
}

.metric-notice {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  min-height: 52px;
  padding: 12px 16px;
  color: #2563eb;
  text-align: left;
  line-height: 20px;
  background-color: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 4px;
  cursor: pointer;
}

.metric-notice:hover,
.metric-notice:focus {
  background-color: #dbeafe;
  border-color: #93c5fd;
  outline: none;
}

.metric-notice__icon {
  flex: 0 0 auto;
  font-size: 16px;
}

.metric-notice__text {
  flex: 1;
  min-width: 0;
  white-space: normal;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.card-content {
  margin-top: 10px;
  margin-right: 10px;
  text-align: center;
  float: right;
}

.card-content > span:first-child {
  color: #64748b;
}
</style>
