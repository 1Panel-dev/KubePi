<template>
  <div>
    <complex-table :data="data">
      <el-table-column sortable :label="$t('commons.table.name')" prop="name">
      </el-table-column>
      <el-table-column sortable :label="$t('business.network.port')" prop="port">
      </el-table-column>
      <el-table-column sortable :label="$t('business.network.protocol')" prop="protocol">
      </el-table-column>
      <el-table-column sortable label="Target" prop="targetPort">
      </el-table-column>
      <el-table-column sortable label="Node Port " prop="nodePort">
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "ResourcePorts",
  components: { ComplexTable },
  props: {
    data: Array
  },
}
</script>

<style scoped>

</style>
