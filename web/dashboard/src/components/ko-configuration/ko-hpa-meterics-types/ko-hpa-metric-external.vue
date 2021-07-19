<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Type" key="external">
          <el-select v-model="row.external.target.type" style="width: 100%" @change="changeExternalType(row)">
            <el-option value="Value" label="Value"></el-option>
            <el-option value="AverageValue" label="Average Value"></el-option>
          </el-select>
        </el-form-item>
      </el-col>
      <el-col :span="12">
        <el-form-item label="Quantity" v-if="row.external.target.type === 'AverageValue'">
          <el-input type="number" v-model="row.external.target.averageValue">
          </el-input>
        </el-form-item>
        <el-form-item label="Quantity" v-if="row.external.target.type === 'Value'">
          <el-input type="number" v-model="row.external.target.value">
          </el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-form-item label="Metric Name">
          <el-input v-model="row.external.metric.name"></el-input>
        </el-form-item>
      </el-col>
    </el-row>
    <el-row>
      <h5>Metric Selector</h5>
      <table style="width: 100%;padding: 0" class="tab-table">
        <tr>
          <th scope="col" width="30%" align="left">
            <label>key</label>
          </th>
          <th scope="Operator" width="30%" align="left">
            <label>key</label>
          </th>
          <th scope="col" width="30%" align="left">
            <label>value</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row,index) in row.external.metric.selector.matchExpressions" v-bind:key="index">
          <td>
            <el-input v-model="row.key"></el-input>
          </td>
          <td>
            <el-select v-model="row.operator" style="width: 100%" @change="changeOperator(row)">
              <el-option label="=" value="In"></el-option>
              <el-option label="≠" value="NotIn"></el-option>
              <el-option label="Exists" value="Exists"></el-option>
              <el-option label="Dose Not Exist" value="DoesNotExist"></el-option>
            </el-select>
          </td>
          <td>
            <el-input v-if="row.operator==='In' || row.operator==='NotIn'" v-model="valueStrings[index]"
                      @change="changeValue(row,index)"></el-input>
            <span v-else>...</span>
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
      </table>
      <el-button style="margin-top: 5px" @click="addRule()">Add Rule</el-button>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "KoHpaMetricExternal",
  components: {},
  props: {
    row: Object
  },
  data () {
    return {
      valueStrings: []
    }
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
    addRule () {
      this.row.external.metric.selector.matchExpressions.push({
        key: "",
        operator: "In",
        values: []
      })
      this.valueStrings.push("")
    },
    changeValue (row, index) {
      if (this.valueStrings[index].indexOf(",") > 0) {
        row.values = this.valueStrings[index].split(",")
      } else {
        row.values.push(this.valueStrings[index])
      }
    },
    handleDelete (index) {
      this.valueStrings.splice(index, 1)
      this.row.external.metric.selector.matchExpressions.splice(index, 1)
    },
    changeOperator (row) {
      if (row.operator === "In" || row.operator === "NotIn") {
        row.values = []
      } else {
        delete row.values
      }
    }
  },
  mounted () {
    if (this.row && this.row.external.metric.selector.matchExpressions.length > 0) {
      for (let i = 0; i < this.row.external.metric.selector.matchExpressions.length; i++) {
        const mach = this.row.external.metric.selector.matchExpressions[i]
        if (mach.values && mach.values instanceof Array) {
          let valueString = ""
          for (let j = 0; j < mach.values.length; j++) {
            if (j + 1 !== mach.values.length) {
              valueString = valueString + mach.values[j] + ","
            } else {
              valueString = valueString + mach.values[j]
            }
          }
          this.valueStrings.push(valueString)
        } else {
          this.valueStrings.push("")
        }
      }
    }
  }
}
</script>

<style scoped>

</style>
