<template>
  <layout-content header="Events">
    <complex-table :pagination-config="page" :data="data" @search="search()" v-loading="loading">
      <template #toolbar>
        <!--        <el-input :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchName"-->
        <!--                  @change="search(true)" @clear="search(true)"></el-input>-->
        <el-select v-model="searchName" @change="search(true)">
          <el-option label="All Namespaces" value=""></el-option>
          <el-option v-for="namespace in namespaces"
                     :key="namespace.metadata.name"
                     :label="namespace.metadata.name"
                     :value="namespace.metadata.name">
          </el-option>
        </el-select>
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
          <el-link>{{ row.involvedObject.kind }} / {{ row.involvedObject.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.eventTime | datetimeFormat }}
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listEvents} from "@/api/events"
import {listNamespace} from "@/api/namespaces"

export default {
  name: "Events",
  components: { ComplexTable, LayoutContent },
  data () {
    return {
      page: {
        pageSize: 20,
        nextToken: "",
      },
      data: [],
      loading: false,
      clusterName: "",
      searchName: "",
      namespaces: []
    }
  },
  methods: {
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: "",
        }
      }
      listEvents(this.clusterName, this.page.pageSize, this.page.nextToken, this.searchName).then(res => {
        this.loading = false
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
      })
    },
    listAllNameSpaces () {
      listNamespace(this.clusterName).then(res => {
        this.namespaces = res.items
      })
    }
  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
    this.listAllNameSpaces()
  }
}
</script>

<style scoped>

</style>