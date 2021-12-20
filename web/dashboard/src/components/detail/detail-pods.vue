<template>
  <div>
    <complex-table :data="pods" v-loading="loading" @search="search">
      <el-table-column :label="$t('commons.table.status')" min-width="45">
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase === 'Running' || row.status.phase === 'Succeeded'" type="success"
                     size="mini" plain round>
            {{ row.status.phase }}
          </el-button>
          <el-button v-if="row.status.phase !== 'Running' && row.status.phase !== 'Succeeded'" type="warning"
                     size="mini" plain round>
            {{ row.status.phase }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="40" prop="metadata.namespace"
                       show-overflow-tooltip/>
      <el-table-column :label="'CPU ' + $t('business.workload.reservation')" min-width="45">
        <template v-slot:default="{row}">
          {{ row.cpuRequest }}
        </template>
      </el-table-column>
      <el-table-column :label="'CPU ' + $t('business.workload.limit')" min-width="45">
        <template v-slot:default="{row}">
          {{ row.cpuLimit }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.memory') + $t('business.workload.reservation')" min-width="50">
        <template v-slot:default="{row}">
          {{ row.memoryRequest }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.memory') + $t('business.workload.limit')" min-width="45">
        <template v-slot:default="{row}">
          {{ row.memoryLimit }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="40" prop="metadata.creationTimestamp" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import {listPodsWithNsSelector} from "@/api/pods"
import {checkPermissions} from "@/utils/permission"
import { cpuUnitConvert, memeryUnitConvert } from "@/utils/unitConvert"

export default {
  name: "KoDetailPods",
  components: { ComplexTable, KoTableOperations },
  props: {
    cluster: String,
    namespace: String,
    selector: String,
    fieldSelector: String,
    allocatable: Object,
  },
  watch: {
    selector: {
      handler (newSelector) {
        if (newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
    fieldSelector: {
      handler (newSelector) {
        if (newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
  },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.terminal"),
          icon: "iconfont iconline-terminalzhongduan",
          click: (row) => {
            this.openTerminal(row)
          },
        },
        {
          label: this.$t("commons.button.logs"),
          icon: "el-icon-tickets",
          click: (row) => {
            this.openTerminalLogs(row)
          },
        },
      ],
      loading: false,
      pods: [],
      podUsage: [],
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "list" })) {
        return
      }
      listPodsWithNsSelector(this.cluster, this.namespace, this.selector, this.fieldSelector).then((res) => {
        this.pods = res.items
        for (const item of this.pods) {
          let cpuLimit = 0
          let memoryLimit = 0
          let cpuRequest = 0
          let memoryRequest = 0
          for (const c of item.spec.containers) {
            if(c.resources?.limits?.cpu) {
              cpuLimit += cpuUnitConvert(c.resources.limits.cpu)
            }
            if(c.resources?.limits?.memory) {
              memoryLimit += memeryUnitConvert(c.resources.limits.memory)
            }
            if(c.resources?.requests?.cpu) {
              cpuRequest += cpuUnitConvert(c.resources.requests.cpu)
            }
            if(c.resources?.requests?.memory) {
              memoryRequest += memeryUnitConvert(c.resources.requests.memory)
            }
          }
          item.cpuLimit = cpuLimit !== 0 ? (cpuLimit + "m (" + (Math.floor(cpuLimit / cpuUnitConvert(this.allocatable.cpu) * 100)) + "%)") : 0
          item.memoryLimit = memoryLimit !== 0 ? (memoryLimit + "Mi (" + (Math.floor(memoryLimit / memeryUnitConvert(this.allocatable.memory) * 100)) + "%)") : 0
          item.cpuRequest = cpuRequest !== 0 ? (cpuRequest  + "m (" + (Math.floor(cpuRequest / cpuUnitConvert(this.allocatable.cpu) * 100)) + "%)") : 0
          item.memoryRequest = memoryRequest !== 0 ? (memoryRequest  + "Mi (" + (Math.floor(memoryRequest / memeryUnitConvert(this.allocatable.memory) * 100)) + "%)") : 0
        }
        this.loading = false
      })
    },
    openDetail (row) {
      this.$router.push({
        name: "PodDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
    openTerminal (row) {
      let c = row.spec.containers[0].name
      let routeUrl = this.$router.resolve({
        path: "/terminal",
        query: {
          cluster: this.cluster,
          namespace: row.metadata.namespace,
          pod: row.metadata.name,
          container: c,
          type: "terminal"
        }
      })
      window.open(routeUrl.href, "_blank")
    },
    openTerminalLogs (row) {
      let c = row.spec.containers[0].name
      let routeUrl = this.$router.resolve({
        path: "/terminal",
        query: {
          cluster: this.cluster,
          namespace: row.metadata.namespace,
          pod: row.metadata.name,
          container: c,
          type: "log"
        }
      })
      window.open(routeUrl.href, "_blank")
    },
    getPodUsage (name, type) {
      let result = "0 m"
      if (this.podUsage.length > 0) {
        for (let item of this.podUsage) {
          if (item.metadata.name === name) {
            let usage = 0
            for (let container of item.containers) {
              if (type === "cpu") {
                if (container.usage.cpu.indexOf("n") > -1) {
                  usage = usage + parseInt(container.usage.cpu)
                }
                if (container.usage.cpu.indexOf("m") > -1) {
                  usage = usage + parseInt(container.usage.cpu) * 1000 * 1000
                }
              }
              if (type === "memory") {
                if (container.usage.memory.indexOf("Ki") > -1) {
                  usage = usage + parseInt(container.usage.memory)
                }
                if (container.usage.memory.indexOf("Mi") > -1) {
                  usage = usage + parseInt(container.usage.memory) * 1000
                }
              }
            }
            const unit = type === "cpu" ? "m" : "Mi"
            if (type === "cpu") {
              result = (usage / 1000000).toFixed(2)
            } else {
              result = (usage / 1000).toFixed(2)
            }
            result = result + unit
          }
        }
      }
      return result
    }
  },
}
</script>

<style scoped>
    .btnSize {
        width: 28px;
        height: 28px;
    }
</style>
