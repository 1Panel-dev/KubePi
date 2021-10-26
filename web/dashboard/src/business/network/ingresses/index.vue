<template>
  <layout-content header="Ingresses">
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <template #header>
        <el-button
                v-has-permissions="{scope:'namespace',apiGroup:'networking.k8s.io',resource:'ingresses',verb:'create'}"
                type="primary" size="small"
                @click="yamlCreate">
          YAML
        </el-button>
        <el-button type="primary" size="small" @click="onCreate"
                   v-has-permissions="{scope:'namespace',apiGroup:'networking.k8s.io',resource:'ingresses',verb:'create'}">
          {{ $t("commons.button.create") }}
        </el-button>
        <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                   v-has-permissions="{scope:'namespace',apiGroup:'networking.k8s.io',resource:'ingresses',verb:'delete'}">
          {{ $t("commons.button.delete") }}
        </el-button>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
      </el-table-column>
      <el-table-column :label="$t('business.network.target')" prop="metadata.rules" min-width="150px"
                       show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(rule,index) in row.spec.rules" :key="index">
            <div v-for="(path,index) in rule.http.paths" :key="index">
              <el-link :href="'http://' + rule.host + (path.path ? path.path : '')" target="_blank">
                {{ "http://" + rule.host + (path.path ? path.path : "") }}
              </el-link>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" min-width="40px" fix>
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
import {downloadYaml} from "@/utils/actions"
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteIngress, getIngress, listIngresses} from "@/api/ingress"
import {mixin} from "@/utils/resourceRoutes"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "Ingresses",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  mixins: [mixin],
  data () {
    return {
      data: [],
      selects: [],
      cluster: "",
      loading: false,
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "IngressEdit",
              params: { namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false }
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "networking.k8s.io",
              resource: "ingresses",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "IngressEdit",
              params: { name: row.metadata.name, namespace: row.metadata.namespace },
              query: { yamlShow: true }
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "networking.k8s.io",
              resource: "ingresses",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getIngress(this.cluster,row.metadata.namespace,row.metadata.name))
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
              apiGroup: "networking.k8s.io",
              resource: "ingresses",
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
      }
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listIngresses(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.paginationConfig.total = res.total
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({
        name: "IngressCreate", query: { yamlShow: false }
      })
    },
    yamlCreate () {
      this.$router.push({
        name: "IngressCreate", query: { yamlShow: true }
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
          this.ps.push(deleteIngress(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteIngress(this.cluster, select.metadata.namespace, select.metadata.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
            .then(() => {
              this.search(true)
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.search(true)
            })
        }
      })
    },
    openDetail (row) {
      this.$router.push({
        name: "IngressDetail",
        params: { name: row.metadata.name, namespace: row.metadata.namespace },
        query: { yamlShow: false }
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
