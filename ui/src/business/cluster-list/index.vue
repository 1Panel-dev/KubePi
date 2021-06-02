<template>
  <layout-content :header="$t('business.cluster.list')">
    <complex-table  :search-config="searchConfig"  :data="data" :pagination-config="paginationConfig">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" >
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button  type="primary" size="small" @click="onImport">
            {{ $t("commons.button.import") }}
          </el-button>
<!--          <el-button size="small" :disabled="clusterSelection.length < 1 || isDeleteButtonDisable" @click="onDelete()" v-permission="['ADMIN','PROJECT_MANAGER']">-->
<!--            {{ $t("commons.button.delete") }}-->
<!--          </el-button>-->
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" min-width="100" prop="name" fix>
        <template v-slot:default="{row}">
          {{ row.name }}
<!--          <el-link v-if="row.status === 'Running'" type="info" @click="goForDetail(row)">{{ row.name }}</el-link>-->
<!--          <span v-if="row.status !== 'Running'">{{ row.name }}</span>-->
        </template>
      </el-table-column>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"

export default {
  name: "ClusterList",
  components: { ComplexTable, LayoutContent },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          { field: "name", label: this.$t("commons.table.name"), component: "FuComplexInput", defaultOperator: "eq" },
          // { field: "created_at", label: this.$t("commons.table.create_time"), component: "FuComplexDateTime", valueFormat: "yyyy-MM-dd HH:mm:ss" },
        ],
      },
      data: [
        {name: "测试"},
        {name: "测试2"}
      ]
    }
  },
  methods: {
    onImport() {
      this.$router.push({ name: "ClusterImport" })
    },
  }
}
</script>

<style scoped>

</style>