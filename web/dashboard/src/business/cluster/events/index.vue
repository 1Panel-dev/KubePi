<template>
  <layout-content header="Events">
    <complex-table :data="data" @search="search" v-loading="loading" :pagination-config="paginationConfig"
                   :search-config="searchConfig">
      <template #toolbar>
      </template>
      <el-table-column :label="$t('business.event.reason')" prop="reason" fix max-width="50px">
        <template v-slot:default="{row}">
          {{ row.reason }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.event.type')" prop="type" fix max-width="50px">
        <template v-slot:default="{row}">
          <el-tag v-if="row.type ==='Normal'">{{ row.type }}</el-tag>
          <el-tag v-else type="danger">{{ row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" fix max-width="50px">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.event.message')" prop="resource" fix min-width="200px"
                       show-overflow-tooltip>
        <template v-slot:default="{row}">
          {{ row.message }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.event.resource')" prop="resource" fix min-width="200px"
                       show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="toResource(row.involvedObject.kind,row.metadata.namespace,row.involvedObject.name)">
            {{ row.involvedObject.kind }} / {{ row.involvedObject.name }}
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.event.count')" prop="count"
                       show-overflow-tooltip>
      </el-table-column>
      <el-table-column :label="$t('commons.table.time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          <!--          <span>{{row.metadata.creationTimestamp | age}}</span>-->
          <!--          <span v-else-if="row.lastTimestamp">{{ row.lastTimestamp | age }}</span>-->
          <span v-if="row.lastTimestamp">{{ row.lastTimestamp | age }}</span>
          <span v-else-if="row.metadata.creationTimestamp">{{ row.metadata.creationTimestamp | age }}</span>
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listEvents} from "@/api/events"
import {mixin} from "@/utils/resourceRoutes"

export default {
  name: "Events",
  components: { ComplexTable, LayoutContent },
  mixins: [mixin],
  data () {
    return {
      data: [],
      loading: false,
      clusterName: "",
      namespaces: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      }
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listEvents(this.clusterName, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.loading = false
        this.data = res.items
        this.paginationConfig.total = res.total
      })
    },
  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
