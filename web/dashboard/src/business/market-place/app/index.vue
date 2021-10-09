<template>
  <layout-content :header="$t('business.chart.app_installed')">
    <complex-table v-loading="loading" :data="data" :search-config="searchConfig" :selects.sync="selects" :pagination-config="paginationConfig"
                   @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'charts',verb:'delete'}" type="primary" size="small"
                     @click="onDelete">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" fix>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" fix>
      </el-table-column>
      <el-table-column :label="'Chart'" prop="chart" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.chart.version')" prop="version" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.version }}
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {searchInstalled} from "@/api/charts"
import ComplexTable from "@/components/complex-table"

export default {
  name: "Apps",
  components: { LayoutContent, ComplexTable },
  props: {},
  data () {
    return {
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      apps: [],
      cluster: "",
      data: [],
      loading: false,
      selects: []
    }
  },
  methods: {
    search (change) {
      if (change) {
        this.paginationConfig = {
          currentPage: 1,
          pageSize: 10,
          total: 0,
        }
      }
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchInstalled(this.cluster, currentPage, pageSize, this.searchConfig.keywords).then(res => {
        this.data = res.data.items
        this.paginationConfig.total = res.data.total
        this.loading = false
      })
    },
    onDelete () {

    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
