<template>
  <layout-content header="Custom Resource Definitions">
    <complex-table :data="data" @search="search" v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate" v-has-permissions="{apiGroup:'apiextensions.k8s.io',resource:'customresourcedefinitions',verb:'create'}">
            YAML
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()" v-has-permissions="{apiGroup:'apiextensions.k8s.io',resource:'customresourcedefinitions',verb:'delete'}">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.spec.names.kind }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="Group" prop="spec.group">
        <template v-slot:default="{row}">
          {{ row.spec.group }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.custom_resource.full_name')" prop="metadata.name">
        <template v-slot:default="{row}">
          {{ row.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.custom_resource.namespaced')" prop="metadata.name">
        <template v-slot:default="{row}">
          {{ row.spec.scope === "Namespaced" ? "True" : "False" }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteCustomResource, listCustomResources} from "@/api/customresourcedefinitions"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "CustomResourceDefinitions",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
      selects: [],
      loading: false,
      cluster: "",
      buttons: [
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "CustomResourceDefinitionEdit",
              params: { name: row.metadata.name }
            })
          },
          disabled:()=>{
            return !checkPermissions({scope:'namespace',apiGroup:"apiextensions.k8s.io",resource:"customresourcedefinitions",verb:"update"})
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
          },
          disabled:()=>{
            return !checkPermissions({scope:'namespace',apiGroup:"apiextensions.k8s.io",resource:"customresourcedefinitions",verb:"delete"})
          }
        },
      ],
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
      listCustomResources(this.cluster,true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "CustomResourceDefinitionCreate",
      })
    },
    onDelete (row) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteCustomResource(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteCustomResource(this.cluster, select.metadata.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
            .then(() => {
              this.search()
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.search()
            })
        }
      })
    },
    openDetail (row) {
      this.$router.push({
        name: "CustomResourceDefinitionDetail",
        params: { name: row.metadata.name },
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
