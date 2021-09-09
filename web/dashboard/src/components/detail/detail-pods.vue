<template>
  <div>
    <complex-table :data="pods" v-loading="loading" @search="search()">
      <el-table-column :label="$t('commons.table.status')" min-width="45">
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase === 'Running' || row.status.phase === 'Succeeded'" type="success" size="mini" plain round>
            {{row.status.phase}}
          </el-button>
          <el-button v-if="row.status.phase !== 'Running' && row.status.phase !== 'Succeeded'" type="warning" size="mini" plain round>
            {{row.status.phase}}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="40" prop="metadata.namespace" show-overflow-tooltip />
      <el-table-column :label="$t('business.cluster.nodes')" min-width="40" prop="spec.nodeName" show-overflow-tooltip />
      <el-table-column :label="$t('business.pod.image')" min-width="120" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(item,index) in row.spec.containers" v-bind:key="index" class="myTag">
            <el-tag type="info" size="small">
              {{ item.image }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="40" prop="metadata.creationTimestamp" fix>
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
import { listPodsWithNsSelector } from "@/api/pods"
import { randomNum } from "@/utils/randomNum"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "KoDetailPods",
  components: { ComplexTable, KoTableOperations },
  props: {
    cluster: String,
    namespace: String,
    selector: String,
    fieldSelector: String,
  },
  watch: {
    selector: {
      handler(newSelector) {
        if (newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
    fieldSelector: {
      handler(newSelector) {
        if (newSelector) {
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
      pods: [],
    }
  },
  methods: {
    search() {
      this.loading = true
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "list" })) {
        return
      }
      listPodsWithNsSelector(this.cluster, this.namespace, this.selector, this.fieldSelector).then((res) => {
        this.pods = res.items
        this.loading = false
      })
    },
    openDetail(row) {
      this.$router.push({
        name: "PodDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
    openTerminal(row) {
      let containers = []
      for (const c of row.spec.containers) {
        containers.push(c.name)
      }
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "terminal",
        key: randomNum(8),
        name: row.metadata.name,
        cluster: this.cluster,
        namespace: row.metadata.namespace,
        pod: row.metadata.name,
        container: containers[0],
        containers: containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
    openTerminalLogs(row) {
      let containers = []
      for (const c of row.spec.containers) {
        containers.push(c.name)
      }
      let existTerminals = this.$store.getters.terminals
      const item = {
        type: "logs",
        key: randomNum(8),
        name: row.metadata.name,
        cluster: this.cluster,
        namespace: row.metadata.namespace,
        pod: row.metadata.name,
        container: containers[0],
        containers: containers,
      }
      existTerminals.push(item)
      this.$store.commit("terminal/TERMINALS", existTerminals)
    },
  },
}
</script>

<style scoped>
.btnSize {
  width: 28px;
  height: 28px;
}
</style>
