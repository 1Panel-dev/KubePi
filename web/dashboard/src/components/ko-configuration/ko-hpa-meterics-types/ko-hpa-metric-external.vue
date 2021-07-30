<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.type')" key="external">
          <el-select v-model="row.external.target.type" style="width: 100%" @change="changeExternalType(row)">
            <el-option value="Value" label="Value"></el-option>
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.external.target.type === 'AverageValue'">
          <el-input type="number" v-model.number="row.external.target.averageValue" required>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.external.target.type === 'Value'">
          <el-input type="number" v-model.number="row.external.target.value" required>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.metrics_name')">
          <el-input v-model="row.external.metric.name"></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row>
      <ko-metric-selector :match-expressions="row.external.metric.selector.matchExpressions"></ko-metric-selector>
    </el-row>
  </div>
</template>

<script>
import KoMetricSelector from "@/components/ko-configuration/ko-metric-selector"
export default {
  name: "KoHpaMetricExternal",
  components: { KoMetricSelector },
  props: {
    row: Object
  },
  methods: {
    changeExternalType (row) {
      if (row.external.target.type === "AverageValue") {
        delete row.external.target.value
        row.external.target.averageValue = 80
      }
      if (row.external.target.type === "Value") {
        delete row.external.target.averageValue
        row.external.target.value = 80
      }
    },
  },
}
</script>

<style scoped>

</style>
