<template>
  <div v-if="!loading">
    <complex-table @row-click="openDetail" :data="data">
      <el-table-column :label="$t('commons.table.status')" min-width="35">
        <template v-slot:default="{row}">
          <el-button v-if="row.container.state.running" type="success" size="mini" plain round>
            Running
          </el-button>
          <el-button v-if="row.container.state.waiting" type="primary" size="mini" plain round>
            Waiting
          </el-button>
          <el-button v-if="row.container.state.pending" type="warning" size="mini" plain round>
            Pending
          </el-button>
          <el-button v-if="row.container.state.terminated" type="warning" size="mini" plain round>
            Terminated
          </el-button>
          <el-button v-if="row.container.state.failed" type="danger" size="mini" plain round>
            Failed
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.ready')" prop="ready" min-width="40">
        <template v-slot:default="{row}">
          <i class="el-icon-check" v-if="row.container.ready" />
          <i class="el-icon-close" v-if="!row.container.ready" />
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="container.name" min-width="50" show-overflow-tooltip />
      <el-table-column :label="$t('business.pod.image')" min-width="150" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div class="myTag">
            <el-tag type="info" size="small">
              {{ row.container.image }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.restarts')" prop="container.restartCount" min-width="30" />
      <el-table-column :label="$t('commons.table.created_time')" min-width="40">
        <template v-slot:default="{row}">
          <span v-if="row.container.started">{{ row.container.state.running.startedAt | age }}</span>
          <span v-if="!row.container.started">-</span>
        </template>
      </el-table-column>

      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>

    <el-row v-if="data.length !== 0" style="margin-top: 20px">
      <ko-detail-containers-info :yamlInfo="form" :containerInfo="containerInfo" />
    </el-row>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import KoDetailContainersInfo from "@/components/detail/detail-containers-info"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "KoDetailContainers",
  components: { ComplexTable, KoTableOperations, KoDetailContainersInfo },
  props: {
    cluster: String,
    yamlInfo: Object,
  },
  watch: {
    yamlInfo: {
      handler(yamlInfo) {
        if (yamlInfo.spec.containers) {
          this.form = yamlInfo
          this.namespace = this.form.metadata.namespace
          this.name = this.form.metadata.name
          this.data = []
          if (yamlInfo.status.containerStatuses) {
            for (const c of yamlInfo.status.containerStatuses) {
              this.data.push({ type: "container", container: c })
            }
          }
          if (yamlInfo.status.initContainerStatuses) {
            for (const c of yamlInfo.status.initContainerStatuses) {
              this.data.push({ type: "init-container", container: c })
            }
          }
          if (this.data.length !== 0) {
            this.containerInfo = {
              name: this.data[0].container.name,
              type: this.data[0].type,
            }
          }
          this.loading = false
        }
      },
      deep: true,
      immediate: true,
    },
  },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.form.detail"),
          icon: "el-icon-view",
          click: (row) => {
            this.openDetail(row)
          }
        },
        {
          label: this.$t("commons.button.logs"),
          icon: "el-icon-tickets",
          click: (row) => {
            this.openTerminalLogs(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods/log", verb: "*" })
          },
        },
        {
          label: this.$t("commons.button.terminal"),
          icon: "iconfont iconline-terminalzhongduan",
          click: (row) => {
            this.openTerminalLogs(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods/exec", verb: "*" })
          },
        },
      ],
      form: {},
      namespace: "",
      name: "",
      data: [],
      containerInfo: {
        type: "",
        name: "",
      },
      loading: true,
    }
  },
  methods: {
    openDetail(row) {
      this.containerInfo = { name: row.container.name, type: row.type }
    },
    openTerminal(row) {
      let routeUrl = this.$router.resolve({ path: "/terminal", query: { cluster: this.cluster, namespace: this.namespace, pod: this.name, container: row.container.name, type: "terminal" } })
      window.open(routeUrl.href, "_blank")
    },
    openTerminalLogs(row) {
      let routeUrl = this.$router.resolve({ path: "/terminal", query: { cluster: this.cluster, namespace: this.namespace, pod: this.name, container: row.container.name, type: "log" } })
      window.open(routeUrl.href, "_blank")
    },
  },
}
</script>

<style lang="scss" scoped>
.complex-table {
  /deep/tbody tr:hover > td {
    cursor: pointer;
  }
}
</style>
