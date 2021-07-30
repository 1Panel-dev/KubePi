<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.resource')">
      <table style="width: 100%" class="tab-table">
        <tr>
          <th scope="col" width="20%" align="left">
            <label>{{ $t("business.network.verbs") }}</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>{{ $t("business.workload.resource") }}</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>{{ $t("business.network.non_resource_url") }}</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>{{ $t("business.network.api_group") }}</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in rules" v-bind:key="index">
          <td>
            <el-select multiple v-model="row.verbs" style="width: 100%">
              <el-option label="create" :value="'create'"></el-option>
              <el-option label="delete" :value="'delete'"></el-option>
              <el-option label="get" :value="'list'"></el-option>
              <el-option label="patch" :value="'patch'"></el-option>
              <el-option label="update" :value="'update'"></el-option>
              <el-option label="watch" :value="'watch'"></el-option>
            </el-select>
          </td>
          <td>
            <el-input v-model="resources[index]"
                      :placeholder="$t('business.access_control.resource_helper')"
                      @change="changeResource(row,index)"></el-input>
          </td>
          <td>
            <el-input v-model="nonResourceURLs[index]"
                      :placeholder="$t('business.access_control.resource_helper')"
                      @change="changeNonResourceURLs(row,index)"></el-input>
          </td>
          <td>
            <el-select multiple v-model="rules[index].groups" style="width: 100%">
              <el-option v-for="(row,index) in groups" :label="row.name" :key="index" :value="row.name"></el-option>
            </el-select>

<!--            <el-input v-model="rules[index].groups"-->
<!--                      :placeholder="$t('business.access_control.resource_helper')"></el-input>-->
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="removeResource(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
      </table>
      <el-button style="margin-top: 5px" @click="addResource()"> {{
          $t("commons.button.add")
        }}{{ $t("business.workload.resource") }}
      </el-button>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import {listApis} from "@/api/apis"

export default {
  name: "KoGrantResource",
  components: { KoCard },
  props: {
    rulesArray: Array,
    cluster: String
  },
  data () {
    return {
      rules: [],
      resources: [],
      nonResourceURLs: [],
      groups: []
    }
  },
  methods: {
    removeResource (index) {
      this.resources.splice(index, 1)
      this.nonResourceURLs.splice(index, 1)
      this.rules.splice(index, 1)
    },
    addResource () {
      this.rules.push({
        verbs: [],
        resources: [],
        nonResourceURLs: [],
      })
      this.resources.push("")
      this.nonResourceURLs.push("")
    },
    changeResource (row, index) {
      if (this.resources[index].indexOf(",") > 0) {
        row.resources = this.resources[index].split(",")
      } else {
        row.resources = this.resources[index] === "" ? [] : [this.resources[index]]
      }
      this.deleteAttributes(row)
    },
    changeNonResourceURLs (row, index) {
      if (this.nonResourceURLs[index].indexOf(",") > 0) {
        row.nonResourceURLs = this.nonResourceURLs[index].split(",")
      } else {
        row.nonResourceURLs = this.nonResourceURLs[index] === "" ? [] : [this.nonResourceURLs[index]]
      }
      this.deleteAttributes(row)
    },
    transform (array, resource) {
      let valueString = ""
      for (let i = 0; i < array?.length; i++) {
        const valueArray = array[i]
        valueString = valueString + valueArray + ","
      }
      valueString = valueString.substr(0, valueString.length - 1)
      resource.push(valueString)
    },
    deleteAttributes (row) {
      if (row.nonResourceURLs?.length === 0 && row.resources?.length !== 0) {
        delete row.nonResourceURLs
      }
      if (row.nonResourceURLs?.length !== 0 && row.resources?.length === 0) {
        delete row.resources
      }
    }
  },
  created () {
    if (this.rulesArray && this.rulesArray.length > 0) {
      this.rules = this.rulesArray
      for (let i = 0; i < this.rulesArray.length; i++) {
        this.transform(this.rulesArray[i].resources, this.resources)
        this.transform(this.rulesArray[i].nonResourceURLs, this.nonResourceURLs)
      }
    } else {
      this.rules.push({
        verbs: [],
        resources: [],
        nonResourceURLs: [],
      })
      this.resources.push("")
      this.nonResourceURLs.push("")
      this.$emit("update:rulesArray", this.rules)
    }
    if (this.cluster) {
      listApis(this.cluster).then(res => {
        this.groups = res.groups
      })
    }
  }
}
</script>

<style scoped>

</style>
