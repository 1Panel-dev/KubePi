<template>
  <layout-content header="Deployments">
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
      <el-table-column sortable label="Name" prop="name" />
      <el-table-column sortable label="Status" prop="status" >
        <template v-slot:default="{row}">
          <el-button v-if="row.status ==='Active'" type="success" size="mini" plain round>
            {{
              row.status
            }}
          </el-button>
          <el-button v-if="row.status ==='Terminating'" type="warning" size="mini" plain round>
            {{
              row.status.phase
            }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column sortable label="Namespace" prop="namespace" />
      <el-table-column sortable label="Endpoints" prop="ip" />
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.created_time | datetimeFormat }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { listDeployments, deleteDeployment } from "@/api/workloads"
import { downloadYaml } from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import ComplexTable from "@/components/complex-table"

export default {
  name: "Deployments",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({ name: "DeploymentEdit", params: { cluster: "songliucs", namespace: row.namespace, name: row.name, yamlShow: false } })
          },
        },
         {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({ name: "DeploymentEdit", params: { cluster: "songliucs", namespace: row.namespace, name: row.name, yamlShow: false } })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({ name: "DeploymentEdit", params: { cluster: "songliucs", namespace: row.namespace, name: row.name, yamlShow: true } })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.name + ".yml", row)
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
    }
  },
  methods: {
    onCreate () {
      this.$router.push({ name: "DeploymentCreate", params: { cluster: "songliucs", yamlShow: false }})
    },
    onDelete(row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteDeployment(this.clusterName, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteDeployment(this.clusterName, select.metadata.name))
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
    search(init) {
      this.loading = true
      this.data = []
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: "",
        }
      }
      listDeployments("songliucs")
        .then((res) => {
          for (const item of res.items) {
            this.data.push({
              status: item.status.conditions[0].status ? "Active" : "Pending",
              name: item.metadata.name,
              namespace: item.metadata.namespace,
              created_time: item.metadata.creationTimestamp,
            })
          }
          this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        })
        .catch((error) => {
          console.log(error)
        })
        .finally(() => {
          this.loading = false
        })
    },
    openDetail (row) {
      this.$router.push({ name: "NamespaceDetail", params: { name: row.metadata.name, cluster: "test1" } })
    },
  },
  mounted() {
    this.search()
  },
}
</script>
