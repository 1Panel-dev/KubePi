<template>
  <layout-content header="Namespaces" v-loading="loading">
    <div style="float: left">
      <el-button v-has-permissions="{scope:'cluster',apiGroup:'',resource:'namespaces',verb:'create'}"
                   type="primary" size="small"
                   @click="yamlCreate">
          YAML
        </el-button>
        <el-button v-has-permissions="{scope:'cluster',apiGroup:'',resource:'namespaces',verb:'create'}"
                   type="primary" size="small"
                   @click="onCreate">
          {{ $t("commons.button.create") }}
        </el-button>
        <el-button v-has-permissions="{scope:'cluster',apiGroup:'',resource:'namespaces',verb:'delete'}"
                   type="primary" size="small"
                   :disabled="selects.length===0" @click="onDelete()">
          {{ $t("commons.button.delete") }}
        </el-button>
    </div>
    <complex-table :selects.sync="selects" :pagination-config="paginationConfig" :search-config="searchConfig"
                   :data="data" @search="search">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
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
          <el-button v-if="row.status.phase ==='Active'" type="success" size="mini" plain round>
            {{
              row.status.phase
            }}
          </el-button>
          <el-button v-if="row.status.phase ==='Terminating'" type="warning" size="mini" plain round>
            {{
              row.status.phase
            }}
          </el-button>
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
import {listNamespace, deleteNamespace, getNamespace} from "@/api/namespaces"
import KoTableOperations from "@/components/ko-table-operations"
import {downloadYaml} from "@/utils/actions"
import {checkPermissions} from "@/utils/permission"


export default {
  name: "NamespaceList",
  components: { KoTableOperations, ComplexTable, LayoutContent },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: "/namespaces/edit/" + row.metadata.name,
              query: { yamlShow: false }
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "namespaces", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "namespaces", verb: "update" })
          },
          click: (row) => {
            this.$router.push({
              path: "/namespaces/edit/" + row.metadata.name,
              query: { yamlShow: true }
            })

          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "namespaces", verb: "get" })
          },
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getNamespace(this.clusterName,row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "namespaces", verb: "delete" })
          },
          click: (row) => {
            this.onDelete(row)
          }
        },
      ],
      data: [],
      selects: [],
      loading: false,
      clusterName: "",
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
    onCreate () {
      this.$router.push({ name: "NamespaceCreate", query: { yamlShow: false } })
    },
    yamlCreate () {
      this.$router.push({ name: "NamespaceCreateYaml", query: { type: "namespaces" } })
    },
    search (resetPage) {
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      this.loading = true
      listNamespace(this.clusterName, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then((res) => {
        this.data = res.items
        this.paginationConfig.total = res.total
      }).catch(error => {
        console.log(error)
      }).finally(() => {
        this.loading = false
      })
    },
    openDetail (row) {
      this.$router.push({ name: "NamespaceDetail", params: { name: row.metadata.name } })
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
          this.ps.push(deleteNamespace(this.clusterName, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteNamespace(this.clusterName, select.metadata.name))
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
              location.reload()
            })
            .catch(() => {
              this.search()
            })
        }
      })
    }
  },
  created () {
    this.clusterName = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
