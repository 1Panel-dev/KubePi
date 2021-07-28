<template>
  <div>
    <complex-table :data="conditions">
      <el-table-column :label="$t('business.pod.type')" prop="type">
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="status">
      </el-table-column>
      <el-table-column :label="$t('business.pod.reason')" prop="reason">
      </el-table-column>
      <el-table-column :label="$t('business.pod.message')" prop="message">
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastHeartbeatTime')" prop="lastHeartbeatTime">
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastTransitionTime')" prop="lastTransitionTime">
      </el-table-column>
    </complex-table>
  </div>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailConditions",
  components: { ComplexTable },
  props: {
    conditions: Array
  },
  data () {
    return {}
  },
  methods: {},
  created () {
  }
}
</script>

<style scoped>

</style>
