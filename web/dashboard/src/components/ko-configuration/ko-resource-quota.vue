<template>
  <div>
    <el-form>
      <el-row>
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="30%" align="left"><label>{{$t('business.workload.resource')}}</label></th>
            <th scope="col" width="30%" align="left"><label>{{$t('business.workload.limit')}}</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in hards" v-bind:key="index">
            <td>
              <ko-form-item itemType="select2" @change="changeType" v-model="row.key" :selections="resource_list" />
            </td>
            <td>
              <ko-form-item itemType="input" v-model="row.value" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
            </td>
          </tr>
        </table>
      </el-row>

      <el-row style="margin-top:20px" v-if="hasSelector">
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="30%" align="left"><label>ScopeName</label></th>
            <th scope="col" width="30%" align="left"><label>Operator</label></th>
            <th scope="col" width="30%" align="left"><label>Values</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in selectors" v-bind:key="index">
            <td>
              <ko-form-item itemType="select" @change="changeScopeName(row)" v-model="row.scopeName" :selections="scope_list" />
            </td>
            <td>
              <ko-form-item itemType="select2" :disabled="isOperatorOnlyExist(row.scopeName)" @change="row.values = ''" v-model="row.operator" :selections="operator_list" />
            </td>
            <td>
              <ko-form-item :disabled="row.operator === 'Exists' || row.operator === 'DoesNotExist'" :placeholder="$t('business.configuration.split_help')" itemType="input" v-model="row.values" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleScopeDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleScopeAdd">{{ $t("commons.button.add") }}</el-button>
            </td>
          </tr>
        </table>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import { parseArryToObj, parseObjToArry } from "@/utils/objArryParse"

export default {
  name: "KoLimitRange",
  components: { KoFormItem },
  props: {
    resourceQuotaObj: Object,
    hasSelector: Boolean,
  },
  watch: {
    resourceQuotaObj: {
      handler(newVal) {
        if (newVal) {
          this.hards = parseObjToArry(newVal.hard)
          if (newVal.scopeSelector?.matchExpressions) {
            for (const item of newVal.scopeSelector.matchExpressions) {
              this.selectors.push({
                scopeName: item.scopeName,
                operator: item.operator,
                values: item.values ? item.values.join(",") : "",
              })
            }
          }
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      hards: [],
      all_resource_list: ["limits.cpu", "limits.memory", "requests.cpu", "requests.memory", "configmaps", "pods", "replicationcontrollers", "resourcequotas", "services", "services.loadbalancers", "services.nodeports", "secrets", "requests.storage", "persistentvolumeclaims"],
      resource_list: [],

      selectors: [],
      scope_list: [
        { label: "Terminating (.spec.activeDeadlineSeconds >= 0)", value: "Terminating" },
        { label: "NotTerminating (.spec.activeDeadlineSeconds is nil)", value: "NotTerminating" },
        { label: "BestEffort (have best effort quality of service)", value: "BestEffort" },
        { label: "NotBestEffort (not have best effort quality of service)", value: "NotBestEffort" },
        { label: "PriorityClass (references the specified priority class)", value: "PriorityClass" },
      ],
      operator_list: ["In", "NotIn", "Exists", "DoesNotExist"],
    }
  },
  methods: {
    handleAdd() {
      var item = {
        key: "",
        value: "",
      }
      this.hards.push(item)
    },
    handleDelete(index) {
      this.hards.splice(index, 1)
    },
    handleScopeAdd() {
      var item = {
        scopeName: "",
        operator: "",
        value: "",
      }
      this.selectors.push(item)
    },
    handleScopeDelete(index) {
      this.selectors.splice(index, 1)
    },
    changeType() {
      let newTypeList = []
      for (const t of this.all_resource_list) {
        let isExist = false
        for (const item of this.hards) {
          if (item.key === t) {
            isExist = true
            break
          }
        }
        if (!isExist) {
          newTypeList.push(t)
        }
      }
      this.resource_list = newTypeList
    },
    changeScopeName(row) {
      if (this.isOperatorOnlyExist(row.scopeName)) {
        row.operator = "Exists"
      }
    },
    isOperatorOnlyExist(scopeName) {
      return scopeName === "Terminating" || scopeName === "NotTerminating" || scopeName === "BestEffort" || scopeName === "NotBestEffort"
    },
    transformation(spec) {
      spec.hard = parseArryToObj(this.hards)
      if (this.selectors.length === 0) {
        delete spec.scopeSelector
        return
      }
      spec.scopeSelector.matchExpressions = []
      for (const scope of this.selectors) {
        spec.scopeSelector.matchExpressions.push({
          scopeName: scope.scopeName,
          operator: scope.operator,
          values: scope.values ? scope.values.split(",") : undefined,
        })
      }
    },
  },
  created() {
    this.changeType()
  },
}
</script>

<style scoped>
</style>
