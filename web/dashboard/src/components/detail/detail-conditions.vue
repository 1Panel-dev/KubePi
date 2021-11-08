<template>
  <div>
    <complex-table :data="conditions">
      <el-table-column :label="$t('business.pod.type')" prop="type" min-width="50" show-overflow-tooltip />
      <el-table-column :label="$t('commons.table.status')" prop="status" min-width="30" />
      <el-table-column :label="$t('commons.table.message')" min-width="100" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-if="row.message">[{{ row.reason }} ]: {{ row.message }}</span>
          <span v-if="!row.message">---</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastUpdateTime')" min-width="50" prop="lastUpdateTime">
        <template v-slot:default="{row}">
          {{ row.lastUpdateTime | age }}
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
    conditions: Array,
  },
}
</script>
