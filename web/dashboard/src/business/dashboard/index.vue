<template>
  <layout-content :header="$t('business.dashboard.dashboard')">
    <el-row :gutter="24">
      <el-col :span="24">
        <el-card shadow="always" style="background-color: #243441;height: 120px">
          <el-row :gutter="24">
            <el-col :span="6">
              <span class="title">{{ $t("commons.table.name") }}</span>
              <div style="text-align: center">
                <h1>测试</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("business.cluster.version") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>v1.18.0</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("business.cluster.nodes") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>2</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("commons.table.created_time") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>16 days ago</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <br>
    <el-row :gutter="24" class="resources">
      <el-col :span="4" v-for="resource in resources" v-bind:key="resource.name">
        <el-card :body-style="{padding: '0px'}" class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              <div>
                <ko-charts :chart-data="resource" :key="resource.name"></ko-charts>
              </div>
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>{{ resource.name }}</span>
                <h1>{{ resource.count }}</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="24">
      <h3>Events</h3>
      <complex-table :pagination-config="page" :data="events" @search="search()" v-loading="loading">
        <el-table-column label="Reason" prop="reason" fix max-width="50px">
          <template v-slot:default="{row}">
            {{ row.reason }}
          </template>
        </el-table-column>
        <el-table-column label="Namespace" prop="namespace" fix max-width="50px">
          <template v-slot:default="{row}">
            {{ row.metadata.namespace }}
          </template>
        </el-table-column>
        <el-table-column label="Message" prop="resource" fix min-width="200px"  show-overflow-tooltip>
          <template v-slot:default="{row}">
            {{ row.message }}
          </template>
        </el-table-column>
        <el-table-column label="Resource" prop="resource" fix min-width="200px"  show-overflow-tooltip>
          <template v-slot:default="{row}">
            <el-link>{{ row.involvedObject.kind }} /  {{ row.involvedObject.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column :label="$t('commons.table.time')" prop="metadata.creationTimestamp" fix>
          <template v-slot:default="{row}">
            {{ row.eventTime | datetimeFormat }}
          </template>
        </el-table-column>
      </complex-table>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoCharts from "@/components/ko-charts"
import {listNamespace} from "@/api/namespaces"
import {listIngresses} from "@/api/ingress"
import {listPvs} from "@/api/pv"
import {listDeployments} from "@/api/workloads"
import {listStatefulSets} from "@/api/statefulsets"
import {listJobs} from "@/api/jobs"
import {listDaemonSets} from "@/api/daemonsets"
import {listServices} from "@/api/services"
import ComplexTable from "@/components/complex-table"
import {listEvents} from "@/api/events"

export default {
  name: "Dashboard",
  components: { ComplexTable, KoCharts, LayoutContent },
  data () {
    return {
      clusterName: "test1",
      resources: [],
      page: {
        pageSize: 10,
        nextToken: "",
        remainCount: 0,
        items: 0,
        currentPage: 1
      },
      events: [],
      loading: false
    }
  },
  methods: {
    listResources () {
      listNamespace(this.clusterName).then(res => {
        const namespaces = {
          name: "Namespaces",
          count: res.items.length,
          data: this.getData(res.items, "status.phase")
        }
        this.resources.push(namespaces)
      })
      listIngresses(this.clusterName).then(res => {
        const ingresses = {
          name: "Ingresses",
          count: res.items.length,
          data: [{
            value: res.items.length,
            name: ""
          }]
        }
        this.resources.push(ingresses)
      })
      listPvs(this.clusterName).then(res => {
        const persistentVolumes = {
          name: "PersistentVolumes",
          count: res.items.length,
          data: this.getData(res.items, "status.phase")
        }
        this.resources.push(persistentVolumes)
      })
      listDeployments(this.clusterName).then(res => {
        const deployments = {
          name: "Deployments",
          count: res.items.length,
          data: this.getData(res.items, "status.conditions.type")
        }
        this.resources.push(deployments)
      })
      listStatefulSets(this.clusterName).then(res => {
        const statefulSets = {
          name: "StatefulSets",
          count: res.items.length,
          data: this.getData(res.items, "status.replicas")
        }
        this.resources.push(statefulSets)
      })
      listJobs(this.clusterName).then(res => {
        const jobs = {
          name: "Jobs",
          count: res.items.length,
          data: this.getData(res.items, "metadata.status.active")
        }
        this.resources.push(jobs)
      })
      listDaemonSets(this.clusterName).then(res => {
        const daemonSets = {
          name: "DaemonSets",
          count: res.items.length,
          data: this.getData(res.items, "status.numberUnavailable")
        }
        this.resources.push(daemonSets)
      })
      listServices(this.clusterName).then(res => {
        const services = {
          name: "Services",
          count: res.items.length,
          data: this.getData(res.items, "status.loadBalancer")
        }
        this.resources.push(services)
      })
      this.search()
    },
    search() {
      this.loading = true
      listEvents(this.clusterName,this.page.pageSize, this.page.nextToken).then(res => {
        this.loading = false
        this.events = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.page.remainCount = res.metadata["remainingItemCount"] ?  res.metadata["remainingItemCount"]: 0
        if (!res.metadata["remainingItemCount"]) {
          this.page.items = res.items.length
        }else {
          this.page.items = 0
        }
      })
    },
    getData (items, keys) {
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
      return result
    },
    traverse (obj, keys) {
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
  created () {
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