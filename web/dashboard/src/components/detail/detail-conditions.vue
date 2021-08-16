<template>
  <div>
    <complex-table :data="conditions">
      <el-table-column :label="$t('business.pod.type')" prop="type" min-width="30" />
      <el-table-column :label="$t('commons.table.status')" prop="status" min-width="30" />
      <el-table-column sortable :label="$t('commons.table.message')" min-width="200">
        <template v-slot:default="{row}">
          <span v-if="row.message">[{{ row.reason }} ]: {{ row.message }}</span>
          <span v-if="!row.message">---</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastUpdateTime')" min-width="40" prop="lastUpdateTime">
        <template v-slot:default="{row}">
          {{ row.lastUpdateTime | age }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.pod.lastTransitionTime')" min-width="40" prop="lastTransitionTime">
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
    conditions: Array,
  },
}
</script>
