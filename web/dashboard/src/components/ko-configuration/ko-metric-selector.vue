<template>
  <div style="margin-top: 20px">
    <h5>{{$t('business.configuration.metrics_selector')}}</h5>
    <table style="width: 100%;padding: 0" class="tab-table" v-if="matchExpressions && matchExpressions.length > 0">
      <tr>
        <th scope="col" width="30%" align="left">
          <label>{{$t('business.workload.key')}}</label>
        </th>
        <th scope="col" width="30%" align="left">
          <label>{{$t('business.workload.operator')}}</label>
        </th>
        <th scope="col" width="30%" align="left">
          <label>{{$t('business.workload.value')}}</label>
        </th>
        <th align="left"></th>
      </tr>
      <tr v-for="(row,index) in matchExpressions" v-bind:key="index">
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
    <el-button style="margin-top: 5px" @click="addRule()">{{ $t("commons.button.add") }}{{$t('business.workload.rule')}}</el-button>
  </div>
</template>

<script>
export default {
  name: "KoMetricSelector",
  components: {},
  props: {
    matchExpressions: Array
  },
  data () {
    return {
      valueStrings: []
    }
  },
  methods: {
    changeValue (row, index) {
      if (this.valueStrings[index].indexOf(",") > 0) {
        row.values = this.valueStrings[index].split(",")
      } else {
        row.values.push(this.valueStrings[index])
      }
    },
    handleDelete (index) {
      this.valueStrings.splice(index, 1)
      this.matchExpressions.splice(index, 1)
    },
    changeOperator (row) {
      if (row.operator === "In" || row.operator === "NotIn") {
        row.values = []
      } else {
        delete row.values
      }
    },
    addRule () {
      this.matchExpressions.push({
        key: "",
        operator: "In",
        values: []
      })
      this.valueStrings.push("")
    },
  },
  created () {
    if (this.matchExpressions && this.matchExpressions.length > 0) {
      for (let i = 0; i < this.matchExpressions.length; i++) {
        const mach = this.matchExpressions[i]
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
