<template>
  <layout-content header="Custom Resource">
    <complex-table :data="data" @search="search" :selects.sync="selects" v-loading="loading"
      :pagination-config="paginationConfig" :search-config="searchConfig"
      :showFullTextSwitch="true" @update:isFullTextSearch="OnIsFullTextSearchChange">
      <template #header>
        <el-button type="primary" size="small" :disabled="selects.length === 0" @click="onDelete()"
          v-has-permissions="{ scope: 'namespace', apiGroup: 'apiextensions.k8s.io', resource: 'customresourcedefinitions', verb: 'delete' }">
          {{ $t("commons.button.delete") }}
        </el-button>
      </template>
      <el-table-column type="selection" fix></el-table-column>

      <el-table-column label="ApiVersion" show-overflow-tooltip prop="apiVersion">
      </el-table-column>
      <el-table-column label="Kind" show-overflow-tooltip prop="kind">
      </el-table-column>
      <el-table-column :label="$t('commons.table.name')" show-overflow-tooltip prop="metadata.name">
        <template v-slot:default="{ row }">
          <el-link @click="openEditPage(row, '0')" class="iconfont iconhuaban88"><span>
              {{ row.metadata.name }}
            </span></el-link>

        </template>
      </el-table-column>
      <el-table-column v-if="show" :label="$t('business.namespace.namespace')" prop="metadata.namespace">
        <template v-slot:default="{ row }">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{ row }">
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
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {
  deleteResource, getResource,
  listResourceByGroup
} from "@/api/customresourcedefinitions"
import { checkPermissions } from "@/utils/permission"
import { searchFullTextItems } from "@/api/fulltextsearch/fulltextsearch"


export default {
  name: "CRList",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  props: {
    names: String,
    version: String,
    group: String,
    scope: {
      type: String,
      default: "Cluster"
    }
  },
  data() {
    return {
      data: [],
      selects: [],
      loading: false,
      cluster: "",
      show: false,
      buttons: [
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.openEditPage(row, '1')
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "apiextensions.k8s.io",
              resource: "customresourcedefinitions",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.view_yaml"),
          icon: "el-icon-view",
          click: (row) => {
            this.openEditPage(row, '0')
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getResource(this.cluster, this.version, this.group, this.names, row.metadata.namespace, row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "apiextensions.k8s.io",
              resource: "customresourcedefinitions",
              verb: "delete"
            })
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
      },
      isFullTextSearch: false
    }
  },
  methods: {
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if(!this.isFullTextSearch){
        listResourceByGroup(this.cluster, this.version, this.group, this.names, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
       })
      } else {
        listResourceByGroup(this.cluster, this.version, this.group, this.names, false)
        .then((res) => {
          const results = searchFullTextItems(res.items,this.searchConfig.keywords);
          this.data =results.slice(this.paginationConfig.currentPage*this.paginationConfig.pageSize-this.paginationConfig.pageSize,this.paginationConfig.currentPage*this.paginationConfig.pageSize)
          this.paginationConfig.total = results.length
        }).finally(() => {
          this.loading = false
        }) 
      }
    },
    onDelete(row) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteResource(this.cluster, this.version, this.group, this.names, row.metadata.namespace.row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteResource(this.cluster, this.version, this.group, this.names, select.metadata.namespace, select.metadata.name))
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
    openEditPage(row, editable) {
      this.$router.push({
        name: "CustomResourceEdit",
        params: {
          name: row.metadata.name,
          cluster: this.cluster,
          version: this.version,
          group: this.group,
          names: this.names,
          editable: editable
        },
        query: {
          namespace: row.metadata.namespace
        }
      })
    },
    //改变选项"是否全文搜索"
    OnIsFullTextSearchChange(val){
      this.isFullTextSearch=val
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.show = this.scope === "Namespaced"
    this.search()
  }
}
</script>

<style scoped></style>
