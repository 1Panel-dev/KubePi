<template>
  <div>
    <complex-table :data="pods" v-loading="loading" @search="search">
      <el-table-column :label="$t('commons.table.version')" prop="name" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span>{{ row.metadata.generation }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.image')" prop="dolphin" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-for="(k, index) in row.spec.template.spec.containers" :key="index">
                <span class="label-custom wd" type="info">{{ k.image }}</span>
          </span>
        </template>
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import {listNsReplicaSetsWorkload} from "@/api/replicasets"
import {checkPermissions} from "@/utils/permission"
import {listPodMetrics} from "@/api/apis"

export default {
  name: "KoDetailReplicasets",
  components: { ComplexTable },
  props: {
    cluster: String,
    namespace: String,
    selector: String,
    fieldSelector: String,
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
      listNsReplicaSetsWorkload(this.cluster, this.namespace, this.selector).then((res) => {
        this.pods = res.items
        this.loading = false
        listPodMetrics(this.cluster, this.namespace, this.selector).then(res => {
          this.podUsage = res.items
        })
      })
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
