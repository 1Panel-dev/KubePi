<template>
  <layout-content header="HorizontalPodAutoscalers">
    <complex-table :pagination-config="page" :data="data" @sarch="search" v-loading="loading">
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
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column label="Workload" prop="spec.scaleTargetRef.kind">
        <template v-slot:default="{row}">
          {{ row.spec.scaleTargetRef.kind }}/{{ row.spec.scaleTargetRef.name }}
        </template>
      </el-table-column>
      <el-table-column label="Min Replicas" prop="spec.scaleTargetRef.minReplicas ">
        <template v-slot:default="{row}">
          {{ row.spec.minReplicas }}
        </template>
      </el-table-column>
      <el-table-column label="Max Replicas" prop="spec.scaleTargetRef.maxReplicas">
        <template v-slot:default="{row}">
          {{ row.spec.maxReplicas }}
        </template>
      </el-table-column>
      <el-table-column label="Current Replicas" prop="status.currentReplicas">
        <template v-slot:default="{row}">
          {{ row.status.currentReplicas }}
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
import ComplexTable from "@/components/complex-table"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteHpa, listHpas} from "@/api/hpa"

export default {
  name: "HPA",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
      page: {
        pageSize: 10,
        nextToken: ""
      },
      selects: [],
      loading: false,
      cluster: "",
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "HPAEdit",
              params: { namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false }
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "HPAEdit",
              params: { name: row.metadata.name, namespace: row.metadata.namespace },
              query: { yamlShow: true }
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", row)
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          }
        },
      ],
    }
  },
  methods: {
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: ""
        }
      }
      listHpas(this.cluster, this.page.pageSize, this.page.nextToken).then(res => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({
        name: "HPACreate",
      })
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
          this.ps.push(deleteHpa(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteHpa(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "HPADetail",
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
