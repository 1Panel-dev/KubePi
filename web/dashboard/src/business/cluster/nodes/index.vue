<template>
  <layout-content header="Nodes">
    <div style="float: left">
      <el-button type="primary" size="small"
                   v-has-permissions="{scope:'cluster',apiGroup:'',resource:'nodes',verb:'update'}"
                   @click="cordon(true)" icon="el-icon-video-pause">
          {{$t('business.node.cordon')}}
        </el-button>
        <el-button type="primary" size="small"
                   v-has-permissions="{scope:'cluster',apiGroup:'',resource:'nodes',verb:'update'}"
                   @click="cordon(false)" icon="el-icon-video-play">
          {{$t('business.node.uncordon')}}
        </el-button>
        <el-button type="primary" size="small"
                   v-has-permissions="{scope:'cluster',apiGroup:'',resource:'nodes',verb:'update'}"
                   @click="drain()" icon="el-icon-refresh-right">
          {{$t('business.node.drain')}}
        </el-button>
    </div>
    <complex-table :data="data" v-loading="loading" :selects.sync="selects" :pagination-config="paginationConfig" @search="search"
                   :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix show-overflow-tooltip max-width="30px">
        <template v-slot:default="{row}">
          <span class="span-link" @click="onDetail(row)"> {{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Internal IP" prop="metadata.name" max-width="30px">
        <template v-slot:default="{row}">
          <div v-for="(address,index) in row.status.addresses" v-bind:key="index">
            <span v-if="address.type === 'InternalIP'">{{ address.address }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')">
        <template v-slot:default="{row}">
          <el-button v-if="row.nodeStatus.indexOf('NotReady') !== -1" type="warning" size="mini" plain round>
            NotReady
          </el-button>
          <el-button v-else type="success" size="mini" plain round>
            Ready
          </el-button>
          <div v-if="row.nodeStatus.indexOf('SchedulingDisabled') !== -1" style="margin-top:3px">
            <el-button type="warning" size="mini" plain round>
              SchedulingDisabled
            </el-button>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.node.role')" prop="metadata.labels" min-width="180px" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(item,name,index) in row.metadata.labels" :key="index" style="display:inline-block">
            <span v-if="item"></span>
            <el-tag v-if="name.indexOf('node-role.kubernetes.io') > -1">{{ name.substring(name.indexOf(".io") + 4) }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="'CPU '+ $t('business.node.usage')" min-width="80px">
        <template v-slot:default="{row}">
          <div v-if="hasMetric === 'true'">
            <div><span>{{row.cpuUsagePersent}}%</span></div>
            <div><span>{{row.cpuUsage | cpu}} / {{ row.status.allocatable.cpu | cpu }} {{ $t('business.node.core') }}</span></div>
          </div>
          <div v-else>
            <div><span> - </span></div>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.memory')+ $t('business.node.usage')" min-width="80px">
        <template v-slot:default="{row}">
          <div v-if="hasMetric === 'true'">
            <div><span>{{row.memoryUsagePersent}}%</span></div>
            <div><span>{{row.memoryUsage | gi-memory}} / {{ row.status.allocatable.memory | gi-memory }} Gi</span></div>
          </div>
          <div v-else>
            <div><span> - </span></div>
          </div>
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {cordonNode, listNodes} from "@/api/nodes"
import KoTableOperations from "@/components/ko-table-operations"
import {evictionPod, listPodsWithNsSelector} from "@/api/pods"
import {checkPermissions} from "@/utils/permission"
import {listNodeMetrics} from "@/api/apis"
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"

export default {
  name: "NodeList",
  components: { KoTableOperations, ComplexTable, LayoutContent },
  data () {
    return {
      data: [],
      loading: false,
      clusterName: "",
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      hasMetric: "false",
      selects: [],
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: "/nodes/edit/" + row.metadata.name,
              query: { yamlShow: false }
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "cluster", apiGroup: "", resource: "nodes", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.view_yaml"),
          icon: "el-icon-view",
          click: (row) => {
            this.$router.push({ path: "/nodes/detail/" + row.metadata.name, query: { yamlShow: "true" } })
          }
        },
      ]
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      const { currentPage, pageSize } = this.paginationConfig
      listNodes(this.clusterName, true, this.searchConfig.keywords, currentPage, pageSize).then(res => {
        let data = res.items
        for (const node of data) {
          node.nodeStatus = "NotReady"
          for(const condition of node.status.conditions) {
            if (condition.type === "Ready") {
              if (condition.status === "True") {
                node.nodeStatus = "Ready"
              }
              break
            }
          }
          if (node.spec.unschedulable) {
            node.nodeStatus += ", SchedulingDisabled"
          }
        }
        this.paginationConfig.total = res.total
        listNodeMetrics(this.clusterName).then(res => {
          this.hasMetric = "true"
          for(const n of data) {
            for (const item of res.items) {
              if (n.metadata.name === item.metadata.name) {
                if (item.usage?.cpu) {
                  n.cpuUsage = cpuUnitConvert(item.usage.cpu) + "m"
                  n.cpuUsagePersent = Math.round((cpuUnitConvert(item.usage.cpu) / cpuUnitConvert(n.status.allocatable.cpu)).toFixed(2) * 100)
                }
                if (item.usage?.memory) {
                  n.memoryUsage = memoryUnitConvert(item.usage.memory).toFixed(2) + "Mi"
                  n.memoryUsagePersent = Math.round((memoryUnitConvert(item.usage.memory) / memoryUnitConvert(n.status.allocatable.memory)).toFixed(2) * 100)
                }
              }
            }
          }
          this.data = data
          this.loading = false
        }).catch(() => {
          this.data = data
          this.loading = false
          this.hasMetric = "false"
        })
      })
    },
    onDetail (row) {
      this.$router.push({ path: "/nodes/detail/" + row.metadata.name, query: { yamlShow: "false" } })
    },
    checkIsValid(isCordon) {
      if (isCordon) {
        for (const item of this.selects) {
          if (item.spec.unschedulable) {
            this.$message({ type: "warning", message: this.$t("business.node.existing_cordoned") })
            return false 
          }
        }
      } else {
        for (const item of this.selects) {
          if (!item.spec.unschedulable) {
            this.$message({ type: "warning", message: this.$t("business.node.existing_actived") })
            return false
          }
        }
      }
      return true
    },
    cordon(isCordon) {
      if (!this.checkIsValid(isCordon)) {
        return
      }
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
          const ps = []
          let data = { spec: { unschedulable: isCordon } }
          for (const item of this.selects) {
            ps.push(cordonNode(this.clusterName, item.metadata.name, data))
          }
          Promise.all(ps)
          .then(() => {
            this.search()
            this.$message({ type: "success", message: this.$t("commons.msg.operation_success") })
          })
        })
    },
    drain() {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        const ps = []
        this.loading = true
        let data = { spec: { unschedulable: true } }
        for (const item of this.selects) {
          ps.push(
            cordonNode(this.clusterName, item.metadata.name, data).then(() => {
              listPodsWithNsSelector(this.clusterName,"","","spec.nodeName="+item.metadata.name).then(res => {
                this.batchDeletePod(res.items, item.metadata.name)
              })
            })
          )
        }
        Promise.all(ps)
        .then(() => {
          this.search()
          this.$message({ type: "success", message: this.$t("commons.msg.operation_success") })
        })
        .catch(() => {
          this.loading = false
        })
      })
    },
    batchDeletePod (items, nodeName) {
      const ps = []
      for (const pod of items) {
        if (pod.spec.nodeName === nodeName) {
          if (pod.metadata.ownerReferences && pod.metadata.ownerReferences[0].kind === "DaemonSet") {
            return
          }
          const rmPod = {
            apiVersion: "policy/v1beta1",
            kind: "Eviction",
            metadata: {
              name: pod.metadata.name,
              namespace: pod.metadata.namespace,
              creationTimestamp: null,
            },
            deleteOptions: {},
          }
          ps.push(evictionPod(this.clusterName, pod.metadata.namespace, pod.metadata.name, rmPod))
        }
      }
      Promise.all(ps)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("business.node.drain_success"),
          })
        })
    }

  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
