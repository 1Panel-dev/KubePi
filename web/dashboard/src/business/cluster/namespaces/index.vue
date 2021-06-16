<template>
  <layout-content header="Namespaces" v-loading="loading">
    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                    >
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
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.cluster.label')" prop="metadata.labels" min-width="200px">
        <template v-slot:default="{row}">
          <el-tag v-for="(value,key,index) in row.metadata.labels" v-bind:key="index" type="info" size="mini">
            {{ key }}={{ value }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="metadata.status" fix>
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase ==='Active'" type="success" size="mini" plain round>{{ row.status.phase }}</el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | datetimeFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
    <ko-page :page.sync="page" @change="listNamespaces('test1')" ></ko-page>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listNamespace} from "@/api/namespace"
import KoPage from "@/components/ko-page"

export default {
  name: "NamespaceList",
  components: { KoPage, ComplexTable, LayoutContent },
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
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          { field: "name", label: this.$t("commons.table.name"), component: "FuComplexInput", defaultOperator: "eq" },
        ],
      },
      data: [],
      selects: [],
      loading: false,
      page: {
        pageSize: 10,
        nextToken:"",
        remainCount: 0,
      }
    }
  },
  methods: {
    onCreate () {
      this.$router.push({ name: "NamespaceCreate" })
    },
    listNamespaces (clusterName) {
      this.loading = true
      listNamespace(clusterName,this.page.pageSize,this.page.nextToken).then((res) => {
        this.data = res.items
        this.page.nextToken =  res.metadata["continue"] ? res.metadata["continue"] : ""
        this.page.remainCount = res.metadata["remainingItemCount"] ? res.metadata["remainingItemCount"] : 0
      }).catch(error => {
        console.log(error)
      }).finally(() => {
        this.loading = false
      })
    },
    openDetail (row) {
      this.$router.push({ name: "NamespaceDetail", params: { name: row.metadata.name, cluster: "test1" } })
    }
  },
  created () {
    this.listNamespaces("test1")
  }
}
</script>

<style scoped>

</style>