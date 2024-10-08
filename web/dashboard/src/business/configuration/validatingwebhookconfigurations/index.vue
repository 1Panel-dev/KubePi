<template>
  <layout-content header="Validating Webhook Configurations">
    <div style="float: left">
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'admissionregistration.k8s.io',resource:'validatingwebhookconfigurations',verb:'create'}"
                  @click="onCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'admissionregistration.k8s.io',resource:'validatingwebhookconfigurations',verb:'delete'}"
                  :disabled="selects.length===0" @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig"
                   :showFullTextSwitch="true" @update:isFullTextSearch="OnIsFullTextSearchChange">
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
import {changeValidatingwebhookconfiguration, deleteValidatingwebhookconfigurations, getValidatingwebhookconfiguration, listValidatingwebhookconfigurations} from "@/api/validatingwebhookconfiguration"
import {checkPermissions} from "@/utils/permission"
import { searchFullTextItems } from "@/api/fulltextsearch/fulltextsearch"
export default {
  name: "Validatingwebhookconfigurations",
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
              name: "ValidatingwebhookconfigurationEdit",
              params: { name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "admissionregistration.k8s.io",
              resource: "validatingwebhookconfigurations",
              verb: "update"
            })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getValidatingwebhookconfiguration(this.cluster, row.metadata.name))
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
              resource: "validatingwebhookconfigurations",
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
      isFullTextSearch: false
    }
  },
  methods: {
    search () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      if(!this.isFullTextSearch){
        listValidatingwebhookconfigurations(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
       })
      } else {
        listValidatingwebhookconfigurations(this.cluster, false)
        .then((res) => {
          const results = searchFullTextItems(res.items,this.searchConfig.keywords);
          this.data =results.slice(this.paginationConfig.currentPage*this.paginationConfig.pageSize-this.paginationConfig.pageSize,this.paginationConfig.currentPage*this.paginationConfig.pageSize)
          this.paginationConfig.total = results.length
        }).finally(() => {
          this.loading = false
        }) 
      }
    },
    onCreate () {
      this.$router.push({
        name: "ValidatingwebhookconfigurationCreateYaml",
        query: { type: "validatingwebhookconfigurations" },
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
          this.ps.push(deleteValidatingwebhookconfigurations(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteValidatingwebhookconfigurations(this.cluster, select.metadata.name))
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
        changeValidatingwebhookconfiguration(this.cluster, row.metadata.name, { metadata: { annotations: { "validatingwebhookconfiguration.kubernetes.io/is-default-class": "false" } } }).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.search(true)
        })
      } else {
        listValidatingwebhookconfigurations(this.cluster).then(res => {
          this.ps = []
          for (const item of res.items) {
            if (this.checkDefault(item)) {
              this.ps.push(changeValidatingwebhookconfiguration(this.cluster, item.metadata.name, { metadata: { annotations: { "validatingwebhookconfiguration.kubernetes.io/is-default-class": "false" } } }))
            }
          }
          this.ps.push(changeValidatingwebhookconfiguration(this.cluster, row.metadata.name, { metadata: { annotations: { "validatingwebhookconfiguration.kubernetes.io/is-default-class": "true" } } }))
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
      return row.metadata.annotations && row.metadata.annotations["validatingwebhookconfiguration.kubernetes.io/is-default-class"] && row.metadata.annotations["validatingwebhookconfiguration.kubernetes.io/is-default-class"] === "true"
    },
    getWebhooksCount(row) {
      return row["webhooks"] ? row["webhooks"].length : 0;
    },
    //改变选项"是否全文搜索"
    OnIsFullTextSearchChange(val){
      this.isFullTextSearch=val
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

