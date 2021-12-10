<template>
  <layout-content :header="$t('business.system.login_log')">
    <complex-table :search-config="searchConfig" :data="data" :pagination-config="paginationConfig" @search="search">
      <el-table-column :label="$t('business.system.username')" prop="userName" fix />
      <el-table-column :label="$t('business.system.ip')" prop="ip" fix />
      <el-table-column :label="$t('business.system.city')" prop="city" fix />
      <el-table-column :label="$t('commons.table.created_time')" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | datetimeFormat }}
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import { searchLoginLogs } from "@/api/systems"

export default {
  name: "LoginLog",
  components: { ComplexTable, LayoutContent },
  data() {
    return {
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          {
            field: "userName",
            label: this.$t("business.system.username"),
            component: "FuComplexInput",
            defaultOperator: "eq",
          },
          {
            field: "ip",
            label: this.$t("business.system.ip"),
            component: "FuComplexInput",
            defaultOperator: "eq",
          },
          {
            field: "city",
            label: this.$t("business.system.city"),
            component: "FuComplexInput",
            defaultOperator: "eq",
          },
        ],
      },
      data: [],
    }
  },
  methods: {
    search(conditions) {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchLoginLogs(currentPage, pageSize, conditions).then((data) => {
        this.loading = false
        this.data = data.data.items
        this.paginationConfig.total = data.data.total
      })
    },
  },
  created() {
    this.search()
  },
}
</script>

<style scoped>
</style>
