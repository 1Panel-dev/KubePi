<template>
  <div>
    <layout-content header="Top Pod">
      <el-alert v-if="showText" :title="$t('business.pod.metric_server_tip')" type="warning" />
      <br>
      <el-button style="margin-left: 20px" :disabled="namespaceDisabled" icon="el-icon-sort-down" type="text" @click="sortBy('namespace')">{{ $t('business.namespace.namespace') }} / {{ $t('business.workload.name') }}</el-button>
      <el-button icon="el-icon-sort-down" :disabled="cpuDisabled" type="text" @click="sortBy('cpu')">CPU</el-button>
      <el-button icon="el-icon-sort-down" :disabled="memoryDisabled" type="text" @click="sortBy('memory')">{{ $t('business.workload.memory') }}</el-button>

      <complex-table :data="data" v-loading="loading">
        <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace" min-width="60" show-overflow-tooltip fix>
          <template v-slot:default="{row}">
            {{ row.metadata.namespace }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip min-width="80">
          <template v-slot:default="{row}">
            <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.workload.container_name')" show-overflow-tooltip min-width="80">
          <template v-slot:default="{row}">
            <div v-for="(c, index) in row.containers" :key="index">
              <div>
                <span>{{ c.name }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="CPU" min-width="35">
          <template v-slot:default="{row}">
            <div v-for="(c, index) in row.containers" :key="index">
              <div>
                <span>{{ c.usage.cpu | cpu }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.workload.memory')" min-width="35">
          <template v-slot:default="{row}">
            <div v-for="(c, index) in row.containers" :key="index">
              <div>
                <span>{{ c.usage.memory | mi-memory }} Mi</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('commons.table.created_time')" min-width="35" show-overflow-tooltip prop="metadata.timestamp" fix>
          <template v-slot:default="{row}">
            {{ row.timestamp | age }}
          </template>
        </el-table-column>
        <el-table-column label="window" show-overflow-tooltip prop="window" min-width="35" fix>
          <template v-slot:default="{row}">
            {{ row.window }}
          </template>
        </el-table-column>
      </complex-table>
    </layout-content>
  </div>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import { listPodMetrics } from "@/api/apis"

export default {
  name: "PodTop",
  components: { LayoutContent, ComplexTable },
  data() {
    return {
      loading: false,
      showText: false,
      data: [],
      clusterName: "",
      namespaceDisabled: false,
      cpuDisabled: true,
      memoryDisabled: false,
    }
  },
  methods: {
    listPodMetric() {
      this.loading = true
      const namespace = sessionStorage.getItem("namespace")
      listPodMetrics(this.clusterName, namespace)
        .then((res) => {
          this.data = res.items
          for (const item of this.data) {
            let cpu = 0
            let memory = 0
            for (const c of item.containers) {
              cpu += Number(c.usage.cpu.replace("n", ""))
              memory += Number(c.usage.memory.replace("Ki", ""))
            }
            item.cpu = cpu / 100000
            item.memory = memory
          }
          this.data.sort(function (a, b) {
            return b.cpu - a.cpu
          })
          this.loading = false
        })
        .catch((error) => {
          this.showText = true
          this.loading = false
          console.log(error.message)
        })
    },
    sortBy(val) {
      this.loading = true
      switch (val) {
        case "namespace":
          this.data.sort(function (a, b) {
            var namespaceA = a.metadata.namespace.toUpperCase()
            var namespaceB = b.metadata.namespace.toUpperCase()
            var nameA = a.metadata.name.toUpperCase()
            var nameB = b.metadata.name.toUpperCase()
            if (namespaceA < namespaceB) {
              return -1
            }
            if (namespaceA > namespaceB) {
              return 1
            }
            if (nameA < nameB) {
              return -1
            }
            if (nameA > nameB) {
              return 1
            }
            return 0
          })
          this.namespaceDisabled = true
          this.cpuDisabled = false
          this.memoryDisabled = false
          this.loading = false
          break
        case "cpu":
          this.data.sort(function (a, b) {
            return b.cpu - a.cpu
          })
          this.namespaceDisabled = false
          this.cpuDisabled = true
          this.memoryDisabled = false
          this.loading = false
          break
        case "memory":
          this.data.sort(function (a, b) {
            return b.memory - a.memory
          })
          this.namespaceDisabled = false
          this.cpuDisabled = false
          this.memoryDisabled = true
          this.loading = false
          break
      }
    },
    openDetail(row) {
      this.$router.push({
        name: "PodDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.listPodMetric()
  },
}
</script>

