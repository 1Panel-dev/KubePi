<template>
  <div>
    <complex-table :data="pods" v-loading="loading" @search="search()">
      <el-table-column  :label="$t('commons.table.status')" min-width="45">
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
      <el-table-column  :label="$t('commons.table.name')" prop="name" min-width="100">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column  :label="$t('business.namespace.namespace')" min-width="45" prop="metadata.namespace"/>
      <el-table-column  :label="$t('business.pod.ready')" min-width="30">
        <template v-slot:default="{row}">
          {{ getPodStatus(row) }}
        </template>
      </el-table-column>
      <el-table-column  :label="$t('business.workload.restarts')" min-width="30">
        <template v-slot:default="{row}">
          {{ getRestartTimes(row) }}
        </template>
      </el-table-column>
      <el-table-column  label="IP" min-width="45" prop="status.podIP"/>
      <el-table-column  :label="$t('business.cluster.nodes')" min-width="40" prop="spec.nodeName"/>
      <el-table-column :label="$t('commons.table.created_time')" min-width="60" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import {listPodsWithNsSelector} from "@/api/pods"

export default {
  name: "ResourcePod",
  components: { ComplexTable },
  props: {
    cluster: String,
    namespace: String,
    selector: String
  },
  data () {
    return {
      loading: false,
      pods: []
    }
  },
  methods: {
    search () {
      this.loading = true
      listPodsWithNsSelector(this.cluster, this.namespace, this.selector).then(res => {
        this.pods = res.items
        this.loading = false
      })
    },
    openDetail (row) {
      this.$router.push({
        name: "PodDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false }
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
  created () {
    if (this.selector !== ""){
      this.search()
    }
  }
}
</script>

<style scoped>

</style>
