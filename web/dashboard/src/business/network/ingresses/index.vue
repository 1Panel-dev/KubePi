<template>
  <layout-content header="Ingresses">
    <div style="float: left">
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
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
      </el-table-column>
      <el-table-column :label="$t('business.network.target')" prop="metadata.rules" min-width="200px"
                       show-overflow-tooltip>
        <template v-slot:default="{row}">
          <div v-for="(rule,index) in row.spec.rules" :key="index">
            <div v-for="(path,index) in rule.http.paths" :key="index">
              <span class="span-link" :href="'http://' + rule.host + (path.path ? path.path : '')" target="_blank">
                {{ "http://" + rule.host + (path.path ? path.path : "") }}
              </span>
              --->
              <span class="span-link" @click="toResource('Service',row.metadata.namespace,getService(path.backend))">
                {{ getService(path.backend) }}:{{ getPort(path.backend) }}
              </span>
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
import {checkApi} from "@/utils/apis"
import {get} from "@/utils/object"

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
      newVersion: false,
      servicePath: "",
      portPath: "",
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: `/ingresses/${row.metadata.namespace}/${row.metadata.name}/edit`,
              query: { yamlShow: false, mode: "edit" }
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
              path: `/ingresses/${row.metadata.namespace}/${row.metadata.name}/edit`,
              query: { yamlShow: true, mode: "edit" }
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
            downloadYaml(row.metadata.name + ".yml", getIngress(this.cluster, row.metadata.namespace, row.metadata.name))
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
        path: `/ingresses/create`,
        query: { yamlShow: false, mode: "create" }
      })
    },
    yamlCreate () {
      this.$router.push({
        path: `/ingresses/create/yaml`,
        query: { type: "ingresses"}
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
    list () {
      this.loading = true
      checkApi(this.cluster, "networking.k8s.io", "v1", "Ingress").then(res => {
        this.newVersion = res
        this.servicePath = this.serviceNamePath()
        this.portPath = this.servicePortPath()
        this.search()
      })
    },
    getService (backend) {
      return get(backend, this.servicePath)
    },
    getPort (backend) {
      return get(backend, this.portPath)
    },
    serviceNamePath () {
      const nestedPath = "service.name"
      const flatPath = "serviceName"
      return this.newVersion ? nestedPath : flatPath
    },
    servicePortPath () {
      const nestedPath = "service.port.number"
      const flatPath = "servicePort"
      return this.newVersion ? nestedPath : flatPath
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
    this.list()
  }
}
</script>

<style scoped>

</style>
