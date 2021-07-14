<template>
  <layout-content header="Jobs">
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="page" @search="search()">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable :label="$t('commons.table.status')" min-width="40">
        <template v-slot:default="{row}">
          <el-button v-if="row.status.succeeded ===1" type="success" size="mini" plain round>
            Succeeded
          </el-button>
          <el-button v-if="row.status.succeeded ===2" type="warning" size="mini" plain round>
            Failed
          </el-button>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="140">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.namespace.namespace')" min-width="40" prop="metadata.namespace" />
      <el-table-column sortable :label="$t('commons.table.status')" min-width="30">
        <template v-slot:default="{row}">
          {{ row.spec.completions}} / {{ row.spec.parallelism }}
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.workload.duration')" min-width="30">
        <template v-slot:default="{row}">
          {{ getDuration(row) }}S
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
import { listJobs, deleteJob } from "@/api/jobs"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"

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
            this.$router.push({ name: "JobEdit", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({ name: "JobEdit", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: true } })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", row)
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
        },
      ],
      loading: false,
      data: [],
      page: {
        pageSize: 10,
        nextToken: "",
      },
      selects: [],
      clusterName: "",
    }
  },
  methods: {
    onCreate() {
      this.$router.push({ name: "JobCreate", query: { yamlShow: false } })
    },
    openDetail(row) {
      this.$router.push({ name: "JobDetail", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteJob(this.clusterName, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteJob(this.clusterName, select.metadata.name))
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
      let endTime = new Date(row.status.completionTime)
      return Math.floor((endTime - startTime) / 1000)
    },
    search(init) {
      this.loading = true
      this.data = []
      this.page.nextToken = init ? "" : this.page.nextToken
      listJobs(this.clusterName, this.page.pageSize, this.page.nextToken)
        .then((res) => {
          this.data = res.items
          this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
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
