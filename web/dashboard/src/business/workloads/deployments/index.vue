<template>
  <layout-content header="Deployments">
    <el-button @click="create">创 建</el-button>
    <el-table :data="data" v-loading="loading">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column sortable label="State" prop="ip" />
      <el-table-column sortable label="Name" prop="ip" />
      <el-table-column sortable label="Namespace" prop="ip" />
      <el-table-column sortable label="Image" prop="ip" />
      <el-table-column sortable label="Endpoints" prop="ip" />
      <el-table-column sortable label="Ready" prop="ip" />
      <el-table-column sortable label="Up-to-data" prop="ip" />
      <el-table-column sortable label="Available" prop="ip" />
      <el-table-column sortable label="Age" prop="ip" />
    </el-table>

  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"

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
    create() {
      this.$router.push({ name: "DeploymentCreate" })
    },
    search() {
      // this.loading = true
      // const { currentPage, pageSize } = this.paginationConfig
      // searchHosts(currentPage, pageSize, condition)
      //   .then((data) => {
      //     this.loading = false
      //     this.data = data.items
      //     this.paginationConfig.total = data.total
      //   })
      //   .catch(() => {
      //     this.loading = false
      //   })
    },
  },
  mounted() {
    this.search()
  },
}
</script>
