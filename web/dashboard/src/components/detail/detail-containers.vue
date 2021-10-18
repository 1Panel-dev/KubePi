<template>
  <div v-if="!loading">
    <complex-table :data="data">
      <el-table-column :label="$t('commons.table.status')" min-width="35">
        <template v-slot:default="{row}">
          <el-button v-if="row.state.running" type="success" size="mini" plain round>
            Running
          </el-button>
          <el-button v-if="row.state.waiting" type="primary" size="mini" plain round>
            Waiting
          </el-button>
          <el-button v-if="row.state.pending" type="warning" size="mini" plain round>
            Pending
          </el-button>
          <el-button v-if="row.state.terminated" type="warning" size="mini" plain round>
            Terminated
          </el-button>
          <el-button v-if="row.state.failed" type="danger" size="mini" plain round>
            Failed
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.ready')" prop="ready" min-width="40">
        <template v-slot:default="{row}">
          <i class="el-icon-check" v-if="row.ready" />
          <i class="el-icon-close" v-if="!row.ready" />
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="50" show-overflow-tooltip />
      <el-table-column :label="$t('business.pod.image')" min-width="170" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div class="myTag">
            <el-tag type="info" size="small">
              {{ row.image }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.restarts')" prop="restartCount" min-width="30" />
      <el-table-column :label="$t('commons.table.created_time')" min-width="40">
        <template v-slot:default="{row}">
          <span v-if="row.started">{{ row.state.running.startedAt | age }}</span>
          <span v-if="!row.started">-</span>
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import { getWorkLoadByName } from "@/api/workloads"

export default {
  name: "KoDetailContainers",
  components: { ComplexTable, KoTableOperations },
  props: {
    cluster: String,
    namespace: String,
    name: String,
  },
  watch: {
    name: {
      handler(newVal) {
        if (newVal) {
          this.search()
        }
      },
      immediate: true,
    },
  },
  data() {
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
      data: [],
      loading: false,
    }
  },
  methods: {
    search() {
      this.loading = true
      getWorkLoadByName(this.cluster, "pods", this.namespace, this.name).then((res) => {
        this.data = []
        if (res.status.containerStatuses) {
          for (const c of res.status.containerStatuses) {
            this.data.push(c)
          }
        }
        if (res.status.initContainerStatuses) {
          for (const c of res.status.initContainerStatuses) {
            this.data.push(c)
          }
        }
        this.loading = false
      })
    },
    openTerminal(row) {
      let routeUrl = this.$router.resolve({ path: "/terminal", query: { cluster: this.cluster, namespace: this.namespace, pod: this.name, container: row.name, type: "terminal" } })
      window.open(routeUrl.href, "_blank")
    },
    openTerminalLogs(row) {
      let routeUrl = this.$router.resolve({ path: "/terminal", query: { cluster: this.cluster, namespace: this.namespace, pod: this.name, container: row.name, type: "log" } })
      window.open(routeUrl.href, "_blank")
    },
  },
}
</script>

<style scoped>
</style>
