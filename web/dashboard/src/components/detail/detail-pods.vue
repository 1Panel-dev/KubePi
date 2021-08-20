<template>
  <div v-if="!loading">
    <complex-table :data="pods" v-loading="loading" @search="search()">
      <el-table-column :label="$t('commons.table.status')" min-width="45">
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase === 'Running' || row.status.phase === 'Succeeded'" type="success" size="mini" plain round>
            {{ $t('commons.status.' + row.status.phase) }}
          </el-button>
          <el-button v-if="row.status.phase !== 'Running' && row.status.phase !== 'Succeeded'" type="warning" size="mini" plain round>
            {{ $t('commons.status.' + row.status.phase) }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="40" prop="metadata.namespace" />
      <el-table-column :label="$t('business.cluster.nodes')" min-width="40" prop="spec.nodeName"  show-overflow-tooltip/>
      <el-table-column sortable :label="$t('business.pod.image')" min-width="120" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(item,index) in row.spec.containers" v-bind:key="index" class="myTag">
            <el-tag type="info" size="small">
              {{ item.image }}
            </el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="50" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import { listPodsWithNsSelector } from "@/api/pods"

export default {
  name: "KoDetailPods",
  components: { ComplexTable },
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
      loading: false,
      pods: [],
    }
  },
  methods: {
    search() {
      this.loading = true
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
    getPodStatus(row) {
      if (row.status.containerStatuses) {
        let readyCount = 0
        for (const c of row.status.containerStatuses) {
          if (c.ready) {
            readyCount++
          }
        }
        return readyCount + "/" + row.status.containerStatuses.length
      }
    },
    getRestartTimes(row) {
      if (row.status.containerStatuses) {
        let restartCount = 0
        for (const c of row.status.containerStatuses) {
          restartCount += c.restartCount
        }
        return restartCount
      }
      return 0
    },
  },
}
</script>

<style scoped>
</style>
