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
            <td v-if="loadSuffix(row) !== 'none'">
              <ko-form-item :deviderName="loadSuffix(row)" itemType="number" v-model="row.value" />
            </td>
            <td v-else>
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

      <div v-if="hasSelector">
        <h4 style="float: left; margin-top:20px">Quota Scopes</h4>
        <el-tooltip class="item" effect="dark" placement="bottom">
          <i style="margin-top: 20px; margin-left: 7px" class="el-icon-question"></i>
          <div slot="content">
            <div><span>{{ $t('business.configuration.best_effort', ['BestEffort']) }}</span></div>
            <ul>pods</ul>
            <div><span>{{ $t('business.configuration.best_effort', ['Terminating, NotTerminating, NotBestEffort and PriorityClass']) }}</span></div>
            <ul>pods</ul>
            <ul>cpu</ul>
            <ul>memory</ul>
            <ul>requests.cpu</ul>
            <ul>requests.memory</ul>
            <ul>limits.cpu</ul>
            <ul>limits.memory</ul>
          </div>
        </el-tooltip>
      </div>

      <el-row v-if="hasSelector">
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
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"

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
          let data = []
          if (newVal.hard) {
            for (const key in newVal.hard) {
              if (Object.prototype.hasOwnProperty.call(newVal.hard, key)) {
                if (key.indexOf(".cpu") !== -1) {
                  data.push({
                    key: key,
                    value: cpuUnitConvert(newVal.hard[key]),
                  })
                } else if (key.indexOf(".memory") !== -1 || key.indexOf(".storage") !== -1) {
                  data.push({
                    key: key,
                    value: memoryUnitConvert(newVal.hard[key]),
                  })
                } else {
                  data.push({
                    key: key,
                    value: newVal.hard[key],
                  })
                }
              }
            }
          }

          this.hards = data
          this.changeType()
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
    loadSuffix(row) {
      if (row.key.indexOf(".memory") !== -1) {
        return "Mi"
      } else if (row.key.indexOf(".cpu") !== -1) {
        return "mCPU"
      } else if (row.key.indexOf(".storage") !== -1) {
        return "Gi"
      } else {
        return "none"
      }
    },
    transformation(spec) {
      let obj = {}
      for (let i = 0; i < this.hards.length; i++) {
        if (this.hards[i].key !== "") {
          switch (this.loadSuffix(this.hards[i])) {
            case "Mi": {
              obj[this.hards[i].key] = this.hards[i].value + "Mi"
              break
            }
            case "mCPU": {
              obj[this.hards[i].key] = this.hards[i].value + "m"
              break
            }
            case "Gi": {
              obj[this.hards[i].key] = this.hards[i].value + "Gi"
              break
            }
            default: {
              obj[this.hards[i].key] = this.hards[i].value
              break
            }
          }
        }
      }

      spec.hard = obj
      if (this.selectors.length === 0) {
        delete spec.scopeSelector
        return
      }
      spec.scopeSelector = {matchExpressions: []}
      for (const scope of this.selectors) {
        spec.scopeSelector.matchExpressions.push({
          scopeName: scope.scopeName,
          operator: scope.operator,
          values: scope.values ? scope.values.split(",") : undefined,
        })
      }
    },
  },
}
</script>

<style scoped>
</style>
