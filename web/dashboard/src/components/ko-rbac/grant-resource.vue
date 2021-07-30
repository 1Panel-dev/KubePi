<template>
  <div style="margin-top: 20px">
    <ko-card title="Grant Resources">
      <table style="width: 100%" class="tab-table">
        <tr>
          <th scope="col" width="20%" align="left">
            <label>Verbs</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>Resource</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>Non-Resource URLs</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>API Groups</label>
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
            <el-input v-model="apiGroups[index]" @change="changeApiGroups(row,index)"
                      :placeholder="$t('business.access_control.resource_helper')"></el-input>
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="removeResource(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
      </table>
      <el-button style="margin-top: 5px" @click="addResource()"> {{ $t("commons.button.add") }} Resource</el-button>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"

export default {
  name: "KoGrantResource",
  components: { KoCard },
  props: {
    rulesArray: Array
  },
  data () {
    return {
      rules: [],
      resources: [],
      nonResourceURLs: [],
      apiGroups: []
    }
  },
  methods: {
    removeResource (index) {
      this.resources.splice(index, 1)
      this.nonResourceURLs.splice(index, 1)
      this.rules.splice(index, 1)
      this.apiGroups.splice(index, 1)
    },
    addResource () {
      this.rules.push({
        verbs: [],
        resources: [],
        nonResourceURLs: [],
        apiGroups: []
      })
      this.resources.push("")
      this.nonResourceURLs.push("")
      this.apiGroups.push("")
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
    changeApiGroups (row, index) {
      if (this.apiGroups[index].indexOf(",") > 0) {
        row.apiGroups = this.apiGroups[index].split(",")
      } else {
        row.apiGroups = this.apiGroups[index] === "" ? [] : [this.apiGroups[index]]
      }
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
        this.transform(this.rulesArray[i].apiGroups, this.apiGroups)
      }
    } else {
      this.rules.push({
        verbs: [],
        resources: [],
        nonResourceURLs: [],
        apiGroups: []
      })
      this.resources.push("")
      this.nonResourceURLs.push("")
      this.apiGroups.push("")
      this.$emit("update:rulesArray", this.rules)
    }
  }
}
</script>

<style scoped>

</style>
