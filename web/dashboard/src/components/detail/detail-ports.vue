<template>
  <div>
    <complex-table :data="data">
      <el-table-column  :label="$t('commons.table.name')" prop="name">
      </el-table-column>
      <el-table-column  :label="$t('business.network.listening_port')" prop="port">
      </el-table-column>
      <el-table-column  :label="$t('business.network.protocol')" prop="protocol">
      </el-table-column>
      <el-table-column   :label="$t('business.network.target_port')" prop="targetPort">
      </el-table-column>
      <el-table-column  label="Node Port" prop="nodePort">
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailPorts",
  components: { ComplexTable },
  props: {
    data: Array
  },
}
</script>

<style scoped>

</style>
