<template>
  <layout-content header="CronJobs">
    <div style="float: left">
      <el-button v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'cronjobs',verb:'create'}"
                  type="primary" size="small"  @click="yamlCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'cronjobs',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'cronjobs',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig"
                   :search-config="searchConfig" @search="search">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="120" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="60" prop="metadata.namespace"/>
      <el-table-column :label="$t('business.workload.schedule')" min-width="40" prop="spec.schedule"/>
      <el-table-column :label="$t('business.workload.lastScheduleTime')" min-width="70"
                       prop="status.lastScheduleTime">
        <template v-slot:default="{row}">
          {{ row.status.lastScheduleTime | age }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.suspend')" min-width="40" prop="spec.suspend">
        <template v-slot:default="{row}">
          {{ row.spec.suspend }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="50" prop="metadata.creationTimestamp" fix>
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
import {listWorkLoads, deleteWorkLoad} from "@/api/workloads"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"
import {checkPermissions} from "@/utils/permission"
import {getCronJobByName} from "@/api/cronjobs"

export default {
  name: "CronJobs",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "CronJobEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "cronjobs", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "CronJobEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "cronjobs", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getCronJobByName(this.clusterName,row.metadata.namespace,row.metadata.name))
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "cronjobs", verb: "delete" })
          },
        },
      ],
      loading: false,
      data: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: "",
      },
      selects: [],
      clusterName: "",
    }
  },
  methods: {
    onCreate () {
      this.$router.push({ name: "CronJobCreate", params: { operation: "create" }, query: { yamlShow: false } })
    },
    openDetail (row) {
      this.$router.push({
        name: "CronJobDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
    yamlCreate () {
      this.$router.push({ name: "CronJobCreateYaml", params: { operation: "create" }, query: { type: "cronjobs" } })
    },
    onDelete (row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteWorkLoad(this.clusterName, "cronjobs", row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteWorkLoad(this.clusterName, "cronjobs", select.metadata.namespace, select.metadata.name))
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
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "cronjobs", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
        .then((res) => {
          this.data = res.items
          this.paginationConfig.total = res.total
        })
        .catch((error) => {
          console.log(error)
        })
        .finally(() => {
          this.loading = false
        })
    },
  },
  mounted () {
    this.clusterName = this.$route.query.cluster
    this.search()
  },
}
</script>
