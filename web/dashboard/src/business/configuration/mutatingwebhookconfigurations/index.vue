<template>
  <layout-content header="Mutatingwebhookconfiguration">
    <div style="float: left">
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'admissionregistration.k8s.io',resource:'mutatingwebhookconfigurations',verb:'create'}"
                  @click="onCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'admissionregistration.k8s.io',resource:'mutatingwebhookconfigurations',verb:'delete'}"
                  :disabled="selects.length===0" @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip></el-table-column>
      <el-table-column label="webhooks" prop="webhooks" fix>
        <template v-slot:default="{row}">
          {{ getWebhooksCount(row) }}
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
import ComplexTable from "@/components/complex-table/index"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {changeMutatingwebhookconfiguration, deleteMutatingwebhookconfigurations, getMutatingwebhookconfiguration, listMutatingwebhookconfigurations} from "@/api/mutatingwebhookconfiguration"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "Mutatingwebhookconfigurations",
  components: { ComplexTable, LayoutContent, KoTableOperations },
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
              name: "MutatingwebhookconfigurationEdit",
              params: { name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "admissionregistration.k8s.io",
              resource: "mutatingwebhookconfigurations",
              verb: "update"
            })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getMutatingwebhookconfiguration(this.cluster, row.metadata.name))
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "admissionregistration.k8s.io",
              resource: "mutatingwebhookconfigurations",
              verb: "delete"
            })
          },
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: "",
      },
    }
  },
  methods: {
    search () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      listMutatingwebhookconfigurations(this.cluster, true, this.searchConfig.keywords, currentPage, pageSize).then((res) => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "MutatingWebhookConfigurationCreateYaml",
        query: { type: "mutatingwebhookconfigurations" },
      })
    },
    onDelete (row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteMutatingwebhookconfigurations(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteMutatingwebhookconfigurations(this.cluster, select.metadata.name))
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
    changeDefault (row) {
      if (this.checkDefault(row)) {
        changeMutatingwebhookconfiguration(this.cluster, row.metadata.name, { metadata: { annotations: { "mutatingwebhookconfiguration.kubernetes.io/is-default-class": "false" } } }).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.search(true)
        })
      } else {
        listMutatingwebhookconfigurations(this.cluster).then(res => {
          this.ps = []
          for (const item of res.items) {
            if (this.checkDefault(item)) {
              this.ps.push(changeMutatingwebhookconfiguration(this.cluster, item.metadata.name, { metadata: { annotations: { "mutatingwebhookconfiguration.kubernetes.io/is-default-class": "false" } } }))
            }
          }
          this.ps.push(changeMutatingwebhookconfiguration(this.cluster, row.metadata.name, { metadata: { annotations: { "mutatingwebhookconfiguration.kubernetes.io/is-default-class": "true" } } }))
          if (this.ps.length !== 0) {
            Promise.all(this.ps)
              .then(() => {
                this.$message({
                  type: "success",
                  message: this.$t("commons.msg.update_success"),
                })
                this.search(true)
              })
              .catch(() => {
                this.search(true)
              })
          }
        })
      }
    },
    checkDefault (row) {
      return row.metadata.annotations && row.metadata.annotations["mutatingwebhookconfiguration.kubernetes.io/is-default-class"] && row.metadata.annotations["mutatingwebhookconfiguration.kubernetes.io/is-default-class"] === "true"
    },
    getWebhooksCount(row) {
      return row["webhooks"] ? row["webhooks"].length : 0;
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  },
}
</script>

<style scoped>
</style>

