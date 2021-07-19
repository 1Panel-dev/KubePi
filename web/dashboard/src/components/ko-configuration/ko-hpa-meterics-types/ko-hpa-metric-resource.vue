<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Resource Name" :key="1">
          <el-select v-model="row.resource.name" style="width: 100%">
            <el-option value="cpu" label="CPU"></el-option>
            <el-option value="memory" label="Memory"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Type" key="resource" :key="2">
          <el-select v-model="row.resource.target.type" style="width: 100%" @change="changeResourceType(row)">
            <el-option value="Utilization" label="Average Utilization"></el-option>
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="Quantity" v-if="row.resource.target.type === 'AverageValue'" :key="3">
          <el-input type="number" v-model="row.resource.target.averageValue">
            <template slot="append">mCpus</template>
          </el-input>
        </el-form-item>
        <el-form-item label="Quantity" v-if="row.resource.target.type === 'Utilization'" :key="4">
          <el-input type="number" v-model="row.resource.target.averageUtilization">
            <template slot="append">%</template>
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "KoHpaMetricResource",
  props: {
    row: Object
  },
  methods: {
    changeResourceType (row) {
      if (row.resource.target.type === "AverageValue") {
        delete row.resource.target.averageUtilization
        row.resource.target.averageValue = 80
      }
      if (row.resource.target.type === "Utilization") {
        delete row.resource.target.averageValue
        row.resource.target.averageUtilization = 80
      }
    },
  },

}
</script>

<style scoped>

</style>
