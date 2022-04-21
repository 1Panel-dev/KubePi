<template>
  <div>
    <complex-table :data="pods" v-loading="loading" @search="search">
      <el-table-column :label="$t('commons.table.version')" prop="name" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span>{{ row.metadata.annotations["deployment.kubernetes.io/revision"] }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.image')" prop="image" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-for="(k, index) in row.spec.template.spec.containers" :key="index">
                <span class="label-custom wd" type="info">{{ k.image }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.time')" prop="time" min-width="80" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span>{{ row.metadata.creationTimestamp | age }}</span>
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import {checkPermissions} from "@/utils/permission"
import {listNsReplicaSetsWorkload} from "@/api/replicasets"
import { patchDeployment } from "@/api/deployments"


export default {
  name: "KoDetailReplicasets",
  components: { ComplexTable, KoTableOperations },
  props: {
    cluster: String,
    namespace: String,
    selector: String,
    fieldSelector: String,
    name: String,
  },
  watch: {
    selector: {
      handler (newSelector) {
        if (newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
    fieldSelector: {
      handler (newSelector) {
        if (newSelector) {
          this.search()
        }
      },
      immediate: true,
    },
  },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.rollback"),
          icon: "el-icon-back",
          click: (row) => {
            this.OptionRollback(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods/log", verb: "*" })
          },
        },
      ],
      loading: false,
      pods: [],
      podUsage: [],
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "list" })) {
        return
      }
      listNsReplicaSetsWorkload(this.cluster, this.namespace, this.selector, this.fieldSelector).then((res) => {
        this.loading = false
        res.items.sort((a, b) => b.metadata.annotations["deployment.kubernetes.io/revision"] - a.metadata.annotations["deployment.kubernetes.io/revision"])
        for (var i = 0; i < res.items.length; i++) {
          this.pods.push(res.items[i])
        }
      })

    },
    OptionRollback (row) {
      this.$confirm(
          this.$t("commons.confirm_message.rollback"),
          this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "warning",
          }).then(() => {
        patchDeployment(this.cluster, row.metadata.namespace, this.name, {"spec": {"template": row.spec.template}})
            .then(() => {
              this.dialogModifyVersionVisible = false
              this.loading = true
              this.$message({
                type: "success",
                message: this.$t("commons.msg.operation_success"),
              })
            })
            .finally(() => {
              this.loading = false
            })
      })
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
  },
}
</script>

<style scoped>
.btnSize {
  width: 28px;
  height: 28px;
}
</style>