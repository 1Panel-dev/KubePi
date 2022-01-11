<template>
  <layout-content :header="$t('business.system.operation_log')">
    <div v-loading="loading">
      <complex-table :search-config="searchConfig" :data="data" :pagination-config="paginationConfig" @search="search">
        <el-table-column :label="$t('business.system.operator')" prop="operator" fix />
        <el-table-column :label="$t('business.system.operation')" prop="operation" fix>
          <template v-slot:default="{row}">
            {{ translate(row.operation) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.system.operation_domain')" prop="operationDomain" fix>
          <template v-slot:default="{row}">
            {{ translate(row.operationDomain) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('business.system.specific_information')" prop="specificInformation" fix />
        <el-table-column :label="$t('commons.table.created_time')" fix>
          <template v-slot:default="{row}">
            {{ row.createAt | datetimeFormat }}
          </template>
        </el-table-column>
      </complex-table>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import { searchOperationLogs } from "@/api/systems"

export default {
  name: "OperationLog",
  components: { ComplexTable, LayoutContent },
  data() {
    return {
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      loading: false,
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          {
            field: "operator",
            label: this.$t("business.system.operator"),
            component: "FuComplexInput",
            defaultOperator: "eq",
          },
          {
            field: "operation",
            label: this.$t("business.system.operation"),
            component: "FuComplexSelect",
            defaultOperator: "eq",
            options: [
              { label: this.translate("post"), value: "post" },
              { label: this.translate("put"), value: "put" },
              { label: this.translate("delete"), value: "delete" },
            ],
          },
          {
            field: "operationDomain",
            label: this.$t("business.system.operation_domain"),
            component: "FuComplexSelect",
            defaultOperator: "eq",
            options: [
              { label: this.translate("clusters"), value: "clusters" },
              { label: this.translate("users"), value: "users" },
              { label: this.translate("clusters_members"), value: "clusters_members" },
              { label: this.translate("clusters_clusterroles"), value: "clusters_clusterroles" },
              { label: this.translate("ldap"), value: "ldap" },
              { label: this.translate("imagerepos"), value: "imagerepos" },
            ],
          },
          {
            field: "specificInformation",
            label: this.$t("business.system.specific_information"),
            component: "FuComplexInput",
            defaultOperator: "eq",
          },
        ],
      },
      data: [],
    }
  },
  methods: {
    translate(a) {
      return this.$t(a)
    },
    search(conditions) {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchOperationLogs(currentPage, pageSize, conditions).then((data) => {
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
