<template>
  <complex-table :data="subjects">
    <el-table-column :label="$t('commons.table.name')" prop="name" min-width="30" />
    <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" min-width="30" />
    <el-table-column label="Kind" prop="kind" min-width="30" />
    <el-table-column :label="$t('business.access_control.api_group')" prop="apiGroup" min-width="30" />
  </complex-table>
</template>

<script>
import ComplexTable from "@/components/complex-table"

export default {
  name: "KoDetailRoleBinding",
  components: {ComplexTable},
  props: {
    subjects: Array
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
