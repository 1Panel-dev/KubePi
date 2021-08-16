<template>
  <div>
    <complex-table :data="conditions">
      <el-table-column :label="$t('business.pod.type')" prop="type" min-width="30">
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="status" min-width="30">
      </el-table-column>
      <el-table-column :label="$t('business.pod.reason')" prop="reason" min-width="60">
      </el-table-column>
      <el-table-column :label="$t('business.pod.message')" prop="message" min-width="120">
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastHeartbeatTime')" min-width="40" prop="lastHeartbeatTime" fix>
        <template v-slot:default="{row}">
          {{ row.lastHeartbeatTime | age }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastTransitionTime')" min-width="40" prop="lastTransitionTime" fix>
        <template v-slot:default="{row}">
          {{ row.lastTransitionTime | age }}
        </template>
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
