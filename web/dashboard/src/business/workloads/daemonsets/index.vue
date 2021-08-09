<template>
  <layout-content header="DaemonSets">
    <complex-table :selects.sync="selects" :data="data" v-loading="loading" :pagination-config="paginationConfig" :search-config="searchConfig" @search="search">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate" v-has-permissions="{apiGroup:'',resource:'daemonsets',verb:'create'}">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()" v-has-permissions="{apiGroup:'',resource:'daemonsets',verb:'delete'}">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="120">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.namespace.namespace')" min-width="80" prop="metadata.namespace" />
      <el-table-column sortable :label="$t('business.pod.ready')" min-width="40" prop="status.numberReady" />
      <el-table-column sortable :label="$t('business.workload.current')" min-width="40" prop="status.currentNumberScheduled" />
      <el-table-column sortable :label="$t('business.workload.desired')" min-width="40" prop="status.desiredNumberScheduled" />
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
import { listWorkLoads, deleteWorkLoad } from "@/api/workloads"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "DaemonSets",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "DaemonSetEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false },
            })
          },
          disabled: () => {
            return !checkPermissions({ apiGroup: "", resource: "daemonsets", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "DaemonSetEdit",
              params: { operation: "edit", namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ apiGroup: "", resource: "daemonsets", verb: "update" })
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
          disabled: () => {
            return !checkPermissions({ apiGroup: "", resource: "daemonsets", verb: "delete" })
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
      this.$router.push({ name: "DaemonSetCreate", params: { operation: "create" }, query: { yamlShow: false } })
    },
    openDetail(row) {
      this.$router.push({ name: "DaemonSetDetail", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteWorkLoad(this.clusterName, "daemonsets", row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteWorkLoad(this.clusterName, "daemonsets", row.metadata.namespace, select.metadata.name))
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
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listWorkLoads(this.clusterName, "daemonsets", true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize)
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
