<template>
  <layout-content header="APIServices">
    <div style="float: left">
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'apiregistration.k8s.io',resource:'apiservices',verb:'delete'}"
                  :disabled="selects.length===0" @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip></el-table-column>
      <el-table-column label="Service"  fix>
        <template v-slot:default="{row}">
          {{ getService(row) }}
        </template>
      </el-table-column>
      <el-table-column label="Service"  fix>
        <template v-slot:default="{row}">
          {{ getService(row) }}
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
import {changeApiservice, deleteApiservices, getApiservice, listApiservices} from "@/api/apiservice"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "Apiservices",
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
              name: "ApiserviceEdit",
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
            downloadYaml(row.metadata.name + ".yml", getApiservice(this.cluster, row.metadata.name))
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
              apiGroup: "apiregistration.k8s.io",
              resource: "apiservices",
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
      listApiservices(this.cluster, true, this.searchConfig.keywords, currentPage, pageSize).then((res) => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "ApiserviceCreateYaml",
        query: { type: "apiservices" },
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
          this.ps.push(deleteApiservices(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteApiservices(this.cluster, select.metadata.name))
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
        changeApiservice(this.cluster, row.metadata.name, { metadata: { annotations: { "apiservice.kubernetes.io/is-default-class": "false" } } }).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.search(true)
        })
      } else {
        listApiservices(this.cluster).then(res => {
          this.ps = []
          for (const item of res.items) {
            if (this.checkDefault(item)) {
              this.ps.push(changeApiservice(this.cluster, item.metadata.name, { metadata: { annotations: { "apiservice.kubernetes.io/is-default-class": "false" } } }))
            }
          }
          this.ps.push(changeApiservice(this.cluster, row.metadata.name, { metadata: { annotations: { "apiservice.kubernetes.io/is-default-class": "true" } } }))
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
      return row.metadata.annotations && row.metadata.annotations["apiservice.kubernetes.io/is-default-class"] && row.metadata.annotations["apiservice.kubernetes.io/is-default-class"] === "true"
    },
    getService(data) {
      if (data.spec.service) {
         return data.spec.service.namespace + "/" + data.spec.service.name;
      } else {
         return "Local";
      }
    },
    getAvailable(data){
      let available="";
      if (data.status && data.status.conditions) {
       for (let i = 0, s = data.status.conditions.length; i < s; i++) {
        if (data.status.conditions[i].type == "Available") {
          if (data.status.conditions[i].status == "True") {
            available = "True";
          } else {
            available =
              data.status.conditions[i].reason +
              "(" +
              data.status.conditions[i].message +
              ")";
          }
       }
      }
     }
     return available
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


