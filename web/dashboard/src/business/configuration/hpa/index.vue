<template>
  <layout-content header="Horizontal Pod Autoscalers">
    <div style="float: left">
      <el-button
              v-has-permissions="{scope:'namespace',apiGroup:'autoscaling',resource:'horizontalpodautoscalers',verb:'create'}"
              type="primary" size="small"
              @click="yamlCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'autoscaling',resource:'horizontalpodautoscalers',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'autoscaling',resource:'horizontalpodautoscalers',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" @search="search" v-loading="loading" :selects.sync="selects"
                   :pagination-config="paginationConfig"
                   :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
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
      <el-table-column :label="$t('commons.table.status')" prop="metadata.status" fix>
        <template v-slot:default="{row}">
          <el-button v-if="row.status.currentReplicas ===row.status.desiredReplicas" type="success" size="mini" plain
                     round>
            Active
          </el-button>
          <el-button v-if="row.status.currentReplicas !==row.status.desiredReplicas" type="danger" size="mini" plain
                     round>
            Pending
          </el-button>
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
import {deleteHpa, getHpa, listHpas} from "@/api/hpa"
import {checkPermissions} from "@/utils/permission"

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
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "autoscaling",
              resource: "horizontalpodautoscalers",
              verb: "update"
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
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "autoscaling",
              resource: "horizontalpodautoscalers",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getHpa(this.cluster,row.metadata.namespace,row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "namespace",
              apiGroup: "autoscaling",
              resource: "horizontalpodautoscalers",
              verb: "delete"
            })
          }
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      }
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listHpas(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "HPACreate", query: { yamlShow: false }
      })
    },
    yamlCreate () {
      this.$router.push({
        name: "HPACreateYaml", query: { type: "horizontalpodautoscalers" }
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
