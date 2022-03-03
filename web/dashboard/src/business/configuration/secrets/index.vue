<template>
  <layout-content header="Secrets">
    <div style="float: left">
      <el-button v-has-permissions="{scope:'namespace',apiGroup:'',resource:'secrets',verb:'create'}"
                  type="primary" size="small"
                  @click="yamlCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'secrets',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'secrets',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search"
                   v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.configuration.type')" prop="type" show-overflow-tooltip>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace" min-width="50px">
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix min-width="50px">
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
import {deleteSecrets, getSecret, listSecrets} from "@/api/secrets"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {checkPermissions} from "@/utils/permission"


const types = ['Opaque','kubernetes.io/dockerconfigjson','kubernetes.io/basic-auth','kubernetes.io/ssh-auth','kubernetes.io/tls']

export default {
  name: "Secrets",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
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
              name: "SecretEdit",
              params: { namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false }
            })
          },
          disabled: (row) => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "",
              resource: "secrets",
              verb: "update"
            }) || types.indexOf(row.type) === -1
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "SecretEdit",
              params: { namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: true }
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "",
              resource: "secrets",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getSecret(this.cluster,row.metadata.namespace,row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "secrets", verb: "delete" })
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
      listSecrets(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.paginationConfig.total = res.total
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({ name: "SecretCreate", query: { yamlShow: false } })
    },
    yamlCreate () {
      this.$router.push({ name: "SecretCreateYaml", query: { type: "secrets" } })
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
          this.ps.push(deleteSecrets(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteSecrets(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "SecretDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
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
