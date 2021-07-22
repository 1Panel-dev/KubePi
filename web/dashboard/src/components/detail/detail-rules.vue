<template>
  <div>
    <complex-table :data="rules">
      <el-table-column :label="$t('business.network.host')" prop="host">
      </el-table-column>
      <el-table-column :label="$t('business.network.path')" prop="path">
      </el-table-column>
      <el-table-column :label="$t('business.network.service_name')" prop="serviceName">
      </el-table-column>
      <el-table-column :label="$t('business.network.service_port')" prop="servicePort">
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailRule",
  components: { ComplexTable },
  props: {
    data: Array
  },
  data () {
    return {
      rules: []
    }
  },
  created () {
    if (this.data) {
      for (const rule of this.data) {
        for (const p of rule.http.paths) {
          let obj = {
            host: rule.host,
            path: p.path ? p.path : "",
            pathType: p.pathType,
            serviceName: p.backend.service ? p.backend.service.name : "",
            servicePort: p.backend.service ? p.backend.service.port?.number : "",
          }
          this.rules.push(obj)
        }
      }
    }
  }
}
</script>

<style scoped>

</style>
