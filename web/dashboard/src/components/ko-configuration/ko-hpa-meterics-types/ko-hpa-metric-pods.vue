<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.type')">
          <el-select v-model="row.pods.target.type" style="width: 100%">
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.pods.target.type === 'AverageValue'">
          <el-input type="number" v-model.number="row.pods.target.averageValue" required>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.pods.target.type === 'Value'">
          <el-input type="number" v-model.number="row.pods.target.value" required>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.metrics_name')">
          <el-input v-model="row.pods.metric.name" required></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row>
      <ko-metric-selector :match-expressions="row.pods.metric.selector.matchExpressions"></ko-metric-selector>
    </el-row>
  </div>
</template>

<script>
import KoMetricSelector from "@/components/ko-configuration/ko-metric-selector"
export default {
  name: "KoHpaMetricPods",
  components: { KoMetricSelector },
  props: {
    row: Object
  },
}
</script>

<style scoped>

</style>
