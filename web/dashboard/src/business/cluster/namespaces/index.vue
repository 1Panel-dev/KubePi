<template>
  <layout-content header="Namespaces">
    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" fix>
        <template v-slot:default="{row}">
          <el-link>{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.cluster.label')" prop="label" fix>
        <template v-slot:default="{row}">
          {{ row.label }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="status" fix>
        <template v-slot:default="{row}">
          {{ row.status }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="age" fix>
        <template v-slot:default="{row}">
          {{ row.age }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listNamespace} from "@/api/namespace"

export default {
  name: "NamespaceList",
  components: { ComplexTable, LayoutContent },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: () => {
          }
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
        ],
      },
      data: [],
      selects: []
    }
  },
  methods: {
    onCreate () {
      this.$router.push({ name: "NamespaceCreate" })
    },
    listNamespaces (clusterName) {

      listNamespace(clusterName).then((res) =>{
        console.log(res)
      }).catch(error => {
        console.log(error)
      })
    }
  },
  created () {
    this.listNamespaces("test1")
  }
}
</script>

<style scoped>

</style>