<template>
  <div v-if="!loading">
    <complex-table :data="form.status.containerStatuses">
      <el-table-column sortable :label="$t('commons.table.status')" min-width="35">
        <template v-slot:default="{row}">
          <el-button v-if="row.state.running" type="success" size="mini" plain round>
            {{ $t('commons.status.Running') }}
          </el-button>
          <el-button v-if="!row.state.running" type="warning" size="mini" plain round>
            {{ $t('commons.status.Failed') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.pod.ready')" prop="ready" min-width="40">
        <template v-slot:default="{row}">
          <i class="el-icon-check" v-if="row.ready" />
          <i class="el-icon-close" v-if="!row.ready" />
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="50" show-overflow-tooltip />
      <el-table-column sortable :label="$t('business.pod.image')" min-width="170" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div class="myTag">
            <el-tag type="info" size="small">
              {{ row.image }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.workload.restarts')" prop="restartCount" min-width="30" />
      <el-table-column sortable :label="$t('commons.table.created_time')" min-width="40">
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
import { randomNum } from "@/utils/randomNum"

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
        if(newVal) {
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
      loading: false,
      form: {
        metadata: {},
        spec: {
          nodeName: "",
        },
        status: {
          containerStatuses: [],
        },
      },
    }
  },
  methods: {
    search() {
      this.loading = true
      getWorkLoadByName(this.cluster, "pods", this.namespace, this.name).then((res) => {
        this.form = res
        this.loading = false
      })
    },
    openTerminal(row) {
      let containers = [row.name]
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "terminal",
        key: randomNum(8),
        name: this.name,
        cluster: this.cluster,
        namespace: this.namespace,
        pod: this.name,
        container: row.name,
        containers: containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
    openTerminalLogs(row) {
      let containers = [row.name]
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "logs",
        key: randomNum(8),
        name: this.name,
        cluster: this.cluster,
        namespace: this.namespace,
        pod: this.name,
        container: row.name,
        containers: containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
  },
}
</script>

<style scoped>
</style>
