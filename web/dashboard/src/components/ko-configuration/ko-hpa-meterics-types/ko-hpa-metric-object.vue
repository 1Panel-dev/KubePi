<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Type" key="external">
          <el-select v-model="row.object.target.type" style="width: 100%" @change="changeExternalType(row)">
            <el-option value="Value" label="Value"></el-option>
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="Quantity" v-if="row.object.target.type === 'AverageValue'">
          <el-input type="number" v-mode.numberl="row.object.target.averageValue">
          </el-input>
        </el-form-item>
        <el-form-item label="Quantity" v-if="row.object.target.type === 'Value'">
          <el-input type="number" v-model.number="row.object.target.value">
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Referent API Version">
          <el-input v-model="row.object.describedObject.apiVersion"></el-input>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="Referent Kind">
          <el-input v-model="row.object.describedObject.kind"></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Referent Name">
          <el-input v-model="row.object.describedObject.name"></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Metric Name">
          <el-input v-model="row.object.metric.name"></el-input>
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
