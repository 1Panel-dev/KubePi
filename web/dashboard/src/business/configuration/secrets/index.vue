<template>
  <layout-content header="Secrets">
    <complex-table :pagination-config="page" :data="data" :selects.sync="selects" @search="search" v-loading="loading">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.configuration.type')" prop="type">
        <template v-slot:default="{row}">
          {{ row.type }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <!--      <el-table-column :label="$t('business.cluster.label')" prop="metadata.labels" min-width="200px">-->
      <!--        <template v-slot:default="{row}">-->
      <!--          <el-tag v-for="(value,key,index) in row.metadata.labels" v-bind:key="index" type="info" size="mini">-->
      <!--            {{ key }}={{ value }}-->
      <!--          </el-tag>-->
      <!--        </template>-->
      <!--      </el-table-column>-->
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | datetimeFormat }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listSecrets} from "@/api/secrets"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"

export default {
  name: "Secrets",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
      page: {
        pageSize: 10,
        nextToken: ""
      },
      selects: [],
      cluster: "",
      loading: false,
      namespaces: [],
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: "/" + row.metadata.namespace + "/secrets/edit/" + row.metadata.name,
              query: { yamlShow: false }
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: "/" + row.metadata.namespace + "/secrets/edit/" + row.metadata.name,
              query: { yamlShow: true }
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", row)
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          }
        },
      ]
    }
  },
  methods: {
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: ""
        }
      }
      listSecrets(this.cluster, this.page.pageSize, this.page.nextToken).then(res => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({ name: "SecretCreate" })
    },
    onDelete () {
    },
    openDetail (row) {
      this.$router.push({
        path: "/" + row.metadata.namespace + "/secrets/detail/" + row.metadata.name,
        query: { yamlShow: false }
      })
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