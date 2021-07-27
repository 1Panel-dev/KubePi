<template>
  <layout-content header="Pods">
    <complex-table :selects.sync="selects" :data="data" v-loading="loading"  @search="search()">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="100">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.namespace.namespace')" min-width="45" prop="metadata.namespace" />
      <el-table-column sortable :label="$t('commons.table.status')" min-width="30">
        <template v-slot:default="{row}">
          {{getPodStatus(row)}}
        </template>
      </el-table-column>
      <el-table-column sortable :label="$t('business.workload.restarts')" min-width="30">
        <template v-slot:default="{row}">
          {{getRestartTimes(row)}}
        </template>
      </el-table-column>
      <el-table-column sortable label="IP" min-width="45" prop="status.podIP" />
      <el-table-column sortable :label="$t('business.cluster.nodes')" min-width="40" prop="spec.nodeName" />
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
import { listPods, deletePod } from "@/api/pods"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"

export default {
  name: "Pods",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data() {
    return {
      buttons: [
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
      selects: [],
      clusterName: "",
    }
  },
  methods: {
    openDetail(row) {
      this.$router.push({ name: "PodDetail", params: { namespace: row.metadata.namespace, name: row.metadata.name }, query: { yamlShow: false } })
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deletePod(this.clusterName, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deletePod(this.clusterName, select.metadata.name))
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
    getPodStatus(row) {
      if (row.status.containerStatuses) {
        let readyCount = 0
        for (const c of row.status.containerStatuses) {
          if (c.ready) {
            readyCount++
          }
        }
        return readyCount + "/" + row.status.containerStatuses.length
      }
      return
    },
    getRestartTimes(row) {
      if (row.status.containerStatuses) {
        let restartCount = 0
        for (const c of row.status.containerStatuses) {
          restartCount += c.restartCount
        }
        return restartCount
      }
      return 0
    },
    search() {
      this.loading = true
      this.data = []
      listPods(this.clusterName)
        .then((res) => {
          this.data = res.items
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
