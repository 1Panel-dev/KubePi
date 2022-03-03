<template>
  <layout-content header="Jobs">
    <div style="float: left">
      <el-button type="primary" size="small" @click="yamlCreate" v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'jobs',verb:'create'}">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate" v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'jobs',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()" v-has-permissions="{scope:'namespace',apiGroup:'batch',resource:'jobs',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig" @search="search">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="140" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" min-width="40" prop="metadata.namespace" />
      <el-table-column :label="$t('business.workload.completions')" min-width="30">
        <template v-slot:default="{row}">
          {{ row.spec.completions }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" min-width="60">
        <template v-slot:default="{row}">
          <el-tag style="margin-left: 5px" v-if="row.status.active" type="info">active: *{{row.status.active}}</el-tag>
          <el-tag style="margin-left: 5px" v-if="row.status.succeeded" type="success">succeeded: {{row.status.succeeded}}</el-tag>
          <el-tag style="margin-left: 5px" v-if="row.status.failed" type="danger">failed: {{row.status.failed}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.workload.duration')" min-width="30">
        <template v-slot:default="{row}">
          {{ getDuration(row) }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="60" prop="metadata.creationTimestamp" fix>
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
import { listWorkLoads, deleteWorkLoad, getWorkLoadByName } from "@/api/workloads"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "Jobs",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.view_form"),
          icon: "el-icon-view",
          click: (row) => {
            this.$router.push({
              name: "JobEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false, readOnly: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "jobs", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.view_yaml"),
          icon: "el-icon-view",
          click: (row) => {
            this.$router.push({
              name: "JobEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: true, readOnly: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "jobs", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getWorkLoadByName(this.clusterName, "jobs", row.metadata.namespace, row.metadata.name))
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "batch", resource: "jobs", verb: "delete" })
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
    onCreate() {
      this.$router.push({ name: "JobCreate", params: { operation: "create" }, query: { yamlShow: false } })
    },
    openDetail(row) {
      this.$router.push({
        name: "JobDetail",
        params: { namespace: row.metadata.namespace, name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
    yamlCreate() {
      this.$router.push({ name: "JobCreateYaml", params: { operation: "create" }, query: { type: "jobs" } })
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteWorkLoad(this.clusterName, "jobs", row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteWorkLoad(this.clusterName, "jobs", select.metadata.namespace, select.metadata.name))
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
    getDuration(row) {
      let startTime = new Date(row.status.startTime)
      if (!row.status.completionTime) {
        return "/"
      } else {
        let endTime = new Date(row.status.completionTime)
        let t = Math.floor((endTime - startTime) / 1000)
        if (t % 60 !== 0) {
          return (t % 60) + " mins"
        }
        if (t % 3600 !== 0) {
          return (t % 60) + " hours"
        }
        return Math.floor((endTime - startTime) / 1000) + "S"
      }
    },
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "jobs", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
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
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.search()
  },
}
</script>
