<template>
  <layout-content header="Deployments">
    <el-button @click="handleYamlCreate">Yaml创建</el-button>
    <el-button @click="handleFormCreate">Form创建</el-button>
    <el-table :data="data" v-loading="loading">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable label="Name" prop="name" />
      <el-table-column sortable label="Status" prop="status" />
      <el-table-column sortable label="Namespace" prop="namespace" />
      <el-table-column sortable label="Endpoints" prop="ip" />
      <el-table-column sortable label="Age" prop="age" />
      <el-table-column>
        <template v-slot:default="{row}">
          <el-button icon="el-icon-edit" @click="handleEditYaml(row)">Y</el-button>
          <el-button icon="el-icon-edit" @click="handleEditForm(row)">F</el-button>
        </template>
      </el-table-column>
    </el-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { listDeployments } from "@/api/workloads"

export default {
  name: "Deployments",
  components: { LayoutContent },
  data() {
    return {
      loading: false,
      data: [],
      selections: [],
    }
  },
  methods: {
    handleYamlCreate() {
      this.$router.push({ name: "DeploymentYaml" })
    },
    handleFormCreate() {
      this.$router.push({ name: "DeploymentForm" })
    },
    handleEditYaml(row) {
      this.$router.push({ name: "DeploymentYaml", params: { cluster: "songliucs", namespace: row.namespace, name: row.name } })
    },
    handleEditForm(row) {
      this.$router.push({ name: "DeploymentForm", params: { cluster: "songliucs", namespace: row.namespace, name: row.name } })
    },
    search() {
      this.loading = true
      listDeployments("songliucs")
        .then((res) => {
          for (const item of res.items) {
            this.data.push({
              status: item.status.conditions[0].status ? "Active" : "Pending",
              name: item.metadata.name,
              namespace: item.metadata.namespace,
            })
          }
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
    this.search()
  },
}
</script>
