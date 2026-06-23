<template>
  <div>
    <layout-content :header="$t('business.workload.top_pod')">
      <el-alert
        v-if="showText"
        class="top-pod-alert"
        :title="$t('business.pod.metric_server_tip')"
        type="warning"
      />
      <div class="top-pod-sort-toolbar">
        <el-button
          class="top-pod-sort-button"
          :class="{'is-active': namespaceDisabled}"
          size="mini"
          :disabled="namespaceDisabled"
          icon="el-icon-sort-down"
          @click="sortBy('namespace')"
        >
          {{ $t('business.namespace.namespace') }} / {{ $t('business.workload.name') }}
        </el-button>
        <el-button
          class="top-pod-sort-button"
          :class="{'is-active': cpuDisabled}"
          size="mini"
          :disabled="cpuDisabled"
          icon="el-icon-sort-down"
          @click="sortBy('cpu')"
        >
          {{ $t("business.common.cpu") }}
        </el-button>
        <el-button
          class="top-pod-sort-button"
          :class="{'is-active': memoryDisabled}"
          size="mini"
          :disabled="memoryDisabled"
          icon="el-icon-sort-down"
          @click="sortBy('memory')"
        >
          {{ $t('business.workload.memory') }}
        </el-button>
        <el-button
          class="top-pod-refresh-button"
          size="mini"
          icon="el-icon-refresh"
          @click="listPodMetric"
        >
          {{ $t('commons.button.refresh') }}
        </el-button>
      </div>

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
        <el-table-column :label="$t('business.common.cpu')" min-width="35">
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
        <el-table-column :label="$t('business.pod.time_stamp')" min-width="35" show-overflow-tooltip prop="metadata.timestamp" fix>
          <template v-slot:default="{row}">
            {{ row.timestamp | age }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.pod.windows')" show-overflow-tooltip prop="window" min-width="35" fix>
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
        .catch(() => {
          this.showText = true
          this.loading = false
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

<style scoped>
.top-pod-sort-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  margin: 12px 0 10px;
}

.top-pod-sort-toolbar .el-button + .el-button {
  margin-left: 0;
}

.top-pod-sort-button.el-button,
.top-pod-refresh-button.el-button {
  padding: 7px 10px;
  color: var(--kp-text-secondary);
  background-color: var(--kp-surface);
  border-color: var(--kp-border);
  box-shadow: none;
}

.top-pod-sort-button.el-button:hover,
.top-pod-sort-button.el-button:focus,
.top-pod-refresh-button.el-button:hover,
.top-pod-refresh-button.el-button:focus {
  color: var(--kp-primary);
  background-color: var(--kp-primary-soft);
  border-color: #bfdbfe;
}

.top-pod-sort-button.el-button.is-active,
.top-pod-sort-button.el-button.is-disabled.is-active,
.top-pod-sort-button.el-button.is-disabled.is-active:hover,
.top-pod-sort-button.el-button.is-disabled.is-active:focus {
  color: var(--kp-primary);
  background-color: var(--kp-primary-soft);
  border-color: #bfdbfe;
  opacity: 1;
  cursor: default;
}
</style>
