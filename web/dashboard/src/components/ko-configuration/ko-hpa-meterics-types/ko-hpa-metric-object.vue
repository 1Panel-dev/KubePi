<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.type')" key="external">
          <el-select v-model="row.object.target.type" style="width: 100%" @change="changeExternalType(row)">
            <el-option value="Value" label="Value"></el-option>
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.object.target.type === 'AverageValue'">
          <el-input type="number" v-model.number="row.object.target.averageValue" required>
          </el-input>
        </el-form-item>
        <el-form-item :label="$t('business.configuration.quantity')" v-if="row.object.target.type === 'Value'">
          <el-input type="number" v-model.number="row.object.target.value" required>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.api_version')">
          <el-input v-model="row.object.describedObject.apiVersion" required></el-input>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.kind')">
          <el-input v-model="row.object.describedObject.kind" required></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.name')">
          <el-input v-model="row.object.describedObject.name" required></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item :label="$t('business.configuration.metrics_name')">
          <el-input v-model="row.object.metric.name" required></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row>
      <ko-metric-selector :match-expressions="row.object.metric.selector.matchExpressions"></ko-metric-selector>
    </el-row>
  </div>
</template>

<script>
import KoMetricSelector from "@/components/ko-configuration/ko-metric-selector"

export default {
  name: "KoHpaMetricObject",
  components: { KoMetricSelector },
  props: {
    row: Object
  },
  methods: {
    changeExternalType (row) {
      if (row.object.target.type === "AverageValue") {
        delete row.object.target.value
        row.object.target.averageValue = 80
      }
      if (row.object.target.type === "Value") {
        delete row.object.target.averageValue
        row.object.target.value = 80
      }
    },
  },
}
</script>

<style scoped>

</style>
