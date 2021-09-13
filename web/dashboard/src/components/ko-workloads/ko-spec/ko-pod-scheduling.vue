<template>
  <div>
    <el-form label-position="top" :disabled="isReadOnly">
      <div v-for="(item, index) in podSchedulings" :key="index">
        <el-card style="margin-top: 10px">
          <el-row>
            <el-button style="float: right;" type="text" @click="handlePodRulesDelete(index)">{{$t("commons.button.delete")}}</el-button>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item :label="$t('business.workload.type')">
                <ko-form-item itemType="radio" v-model="item.type" :radios="type_list" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('business.workload.priority')">
                <ko-form-item itemType="radio" v-model="item.priority" :radios="priority_list" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="Namespace">
                <ko-form-item itemType="select2" v-model="item.namespaces" multiple :selections="namespace_list" />
              </el-form-item>
            </el-col>
            <el-col :span="12" v-if="item.priority === 'Preferred'">
              <el-form-item :label="$t('business.workload.weight')">
                <ko-form-item itemType="number" v-model="item.weight" />
              </el-form-item>
            </el-col>
          </el-row>

          <div style="margin-top: 20px"><span>matchExpressions</span></div>
          <ko-match-table :matchTables.sync="item.rules" :isReadOnly="isReadOnly" />

          <div style="margin-top: 20px"><span>matchLabels</span></div>
          <table style="width: 100%; margin-top: 5px" class="tab-table">
            <tr v-for="(row, index) in item.labelRules" v-bind:key="'l'+index">
              <td width="48%">
                <ko-form-item itemType="input" v-model="row.key" />
              </td>
              <td width="48%">
                <ko-form-item itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleMatchLabelDelete(item, index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleMatchLabelAdd(item)">{{ $t("commons.button.add") }}</el-button>
              </td>
            </tr>
          </table>

          <el-row style="margin-top: 10px">
            <el-col :span=24>
              <el-form-item :label="$t('business.workload.topology_key')" required>
                <ko-form-item itemType="input" placeholder="e.g. failure-domain.beta.kubernetes.io/zone" v-model="item.topologyKey" />
              </el-form-item>
            </el-col>
          </el-row>
        </el-card>
      </div>
      <el-button @click="handlePodRulesAdd()">{{ $t("commons.button.add") }} Pod {{$t('business.workload.rules')}}</el-button>

      <div v-for="(item, index) in nodeSchedulings" :key="'n'+index">
        <el-card style="margin-top: 10px">
          <el-row>
            <el-button style="float: right; padding: 3px 0" type="text" @click="handleNodeRulesDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </el-row>
          <el-row>
            <el-col :span="12">
              <el-form-item :label="$t('business.workload.priority')">
                <ko-form-item itemType="radio" v-model="item.priority" :radios="priority_list" />
              </el-form-item>
            </el-col>
            <el-col :span="12" v-if="item.priority === 'Preferred'">
              <el-form-item :label="$t('business.workload.weight')">
                <ko-form-item itemType="number" v-model="item.weight" />
              </el-form-item>
            </el-col>
          </el-row>

          <div style="margin-top: 20px"><span>matchLabels</span></div>
          <ko-match-table :matchTables.sync="item.rules" :isReadOnly="isReadOnly" />
          <div style="margin-top: 20px"><span>matchFields</span></div>
          <ko-match-table :matchTables.sync="item.fields" :isReadOnly="isReadOnly" />
        </el-card>
      </div>
      <el-button @click="handleNodeRulesAdd()">{{ $t("commons.button.add") }} Node {{$t('business.workload.rules')}}</el-button>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoMatchTable from "@/components/ko-workloads/ko-match-table"

export default {
  name: "KoPodScheduling",
  components: { KoFormItem, KoMatchTable },
  props: {
    podSchedulingParentObj: Object,
    namespaceList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    namespaceList: {
      handler(newName) {
        this.namespace_list = newName
      },
      immediate: true,
    },
  },
  data() {
    return {
      type_list: [
        { label: this.$t("business.workload.affinity"), value: "Affinity" },
        { label: this.$t("business.workload.anti_affinity"), value: "Anti-Affinity" },
      ],
      priority_list: [
        { label: this.$t("business.workload.preferred"), value: "Preferred" },
        { label: this.$t("business.workload.required"), value: "Required" },
      ],
      namespace_operation_list: [
        { label: this.$t("business.workload.this_namespace"), value: "podNamespace" },
        { label: this.$t("business.workload.these_namespace"), value: "selectNamespace" },
      ],
      operator_list: [
        { label: "<", value: "Lt" },
        { label: ">", value: "Gt" },
        { label: "Exists", value: "Exists" },
        { label: "Not Exist", value: "DoesNotExist" },
        { label: "=", value: "In" },
        { label: "!=", value: "NotIn" },
      ],
      namespace_list: [],
      podSchedulings: [],
      nodeSchedulings: [],
      priority: "Preferred",
    }
  },
  methods: {
    handlePodRulesAdd() {
      var item = {
        type: "Affinity",
        priority: "Preferred",
        namespaces: "",
        weight: 1,
        rules: [],
        labelRules: [],
        topologyKey: "",
      }
      this.podSchedulings.push(item)
    },
    handlePodRulesDelete(index) {
      this.podSchedulings.splice(index, 1)
    },
    handleMatchLabelAdd(schedulingItem) {
      var item = {
        key: "",
        value: "",
      }
      schedulingItem.labelRules.push(item)
    },
    handleMatchLabelDelete(schedulingItem, index) {
      schedulingItem.labelRules.splice(index, 1)
    },
    handleNodeRulesAdd() {
      var item = {
        priority: "Preferred",
        weight: 1,
        rules: [],
        fields: [],
      }
      this.nodeSchedulings.push(item)
    },
    handleNodeRulesDelete(index) {
      this.nodeSchedulings.splice(index, 1)
    },

    getMatchExpress(val) {
      let matchs = []
      if (val.length === 0) {
        return matchs
      }
      for (const rule of val) {
        if (rule.value) {
          matchs.push({
            key: rule.key,
            operator: rule.operator,
            values: rule.value.split(","),
          })
        } else {
          matchs.push({
            key: rule.key,
            operator: rule.operator,
          })
        }
      }
      return matchs
    },

    valueTrans(type, priority, s) {
      let namespaces = s.namespaces || ""
      let rules = []
      if (s.labelSelector.matchExpressions) {
        for (const express of s.labelSelector.matchExpressions) {
          rules.push({ key: express.key, operator: express.operator, value: express.values ? express.values.join(",") : null })
        }
      }
      let labelRules = []
      if (s.labelSelector.matchLabels) {
        for (const key in s.labelSelector.matchLabels) {
          if (Object.prototype.hasOwnProperty.call(s.labelSelector.matchLabels, key)) {
            labelRules.push({
              key: key,
              value: s.labelSelector.matchLabels[key],
            })
          }
        }
      }
      const topologyKey = s.topologyKey ? s.topologyKey : ""
      this.podSchedulings.push({
        type: type,
        priority: priority,
        weight: s.weight || null,
        namespaces: namespaces || "",
        rules: rules,
        labelRules: labelRules,
        topologyKey: topologyKey,
      })
    },
    getMatchLabels(val) {
      let obj = {}
      if (val.length === 0) {
        return obj
      }
      for (let i = 0; i < val.length; i++) {
        if (val[i].key !== "") {
          obj[val[i].key] = val[i].value
        }
      }
      return obj
    },
    checkIsValid() {
      let isValid = true
      for (const po of this.podSchedulings) {
        if (po.topologyKey === "") {
          isValid = false
        }
      }
      return isValid
    },
    transformation(parentFrom) {
      parentFrom.affinity = {}
      if (this.podSchedulings.length !== 0) {
        for (const pS of this.podSchedulings) {
          let itemAdd = {}
          itemAdd.namespaces = (pS.namespaces && pS.namespaces.length !== 0) ? pS.namespaces : undefined
          itemAdd.topologyKey = pS.topologyKey || undefined
          const matchs = this.getMatchExpress(pS.rules)
          const labelMatchs = this.getMatchLabels(pS.labelRules)
          switch (pS.type + "+" + pS.priority) {
            case "Affinity+Required":
              if (!parentFrom.affinity.podAffinity) {
                parentFrom.affinity.podAffinity = {}
              }
              parentFrom.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
              itemAdd.labelSelector = { matchExpressions: matchs, matchLabels: labelMatchs }
              parentFrom.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Affinity+Preferred":
              if (!parentFrom.affinity.podAffinity) {
                parentFrom.affinity.podAffinity = {}
              }
              parentFrom.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
              itemAdd.podAffinityTerm = { labelSelector: { matchExpressions: matchs, matchLabels: labelMatchs } }
              itemAdd.weight = pS.weight
              parentFrom.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Anti-Affinity+Required":
              if (!parentFrom.affinity.podAntiAffinity) {
                parentFrom.affinity.podAntiAffinity = {}
              }
              parentFrom.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
              itemAdd.labelSelector = { matchExpressions: matchs, matchLabels: labelMatchs }
              parentFrom.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Anti-Affinity+Preferred":
              if (!parentFrom.affinity.podAntiAffinity) {
                parentFrom.affinity.podAntiAffinity = {}
              }
              parentFrom.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
              itemAdd.podAffinityTerm = { labelSelector: { matchExpressions: matchs, matchLabels: labelMatchs } }
              itemAdd.weight = pS.weight
              parentFrom.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
          }
        }
      }
      parentFrom.nodeAffinity = {}
      if (this.nodeSchedulings.length !== 0) {
        for (const nS of this.nodeSchedulings) {
          const matchs = this.getMatchExpress(nS.rules)
          const fields = this.getMatchExpress(nS.fields)
          let itemAdd = {}
          switch (nS.priority) {
            case "Preferred":
              parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
              itemAdd.weight = nS.weight
              itemAdd.preference = { matchExpressions: matchs, matchFields: fields }
              parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Required":
              if (!parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution = {}
              }
              parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms = []
              itemAdd.matchExpressions = matchs
              itemAdd.matchFields = fields
              parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.push(itemAdd)
              break
          }
        }
      }
    },
  },
  mounted() {
    this.namespace_list = this.namespaceList
    this.podSchedulings = []
    if (this.podSchedulingParentObj) {
      if (this.podSchedulingParentObj.affinity) {
        if (this.podSchedulingParentObj.affinity.podAffinity) {
          if (this.podSchedulingParentObj.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podSchedulingParentObj.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Affinity", "Required", s)
            }
          }
          if (this.podSchedulingParentObj.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podSchedulingParentObj.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              if (s.podAffinityTerm) {
                this.valueTrans("Affinity", "Preferred", s.podAffinityTerm)
              }
            }
          }
        }
        if (this.podSchedulingParentObj.affinity.podAntiAffinity) {
          if (this.podSchedulingParentObj.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podSchedulingParentObj.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Anti-Affinity", "Required", s)
            }
          }
          if (this.podSchedulingParentObj.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podSchedulingParentObj.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              if (s.podAffinityTerm) {
                this.valueTrans("Anti-Affinity", "Preferred", s.podAffinityTerm)
              }
            }
          }
        }
      }

      if (this.podSchedulingParentObj.nodeAffinity) {
        if (this.podSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
          if (this.podSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms) {
            const schedulings = this.podSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms
            for (const s of schedulings) {
              let rules = []
              if (s.matchExpressions) {
                for (const express of s.matchExpressions) {
                  rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
              let fields = []
              if (s.matchFields) {
                for (const express of s.matchFields) {
                  fields.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
              this.nodeSchedulings.push({ priority: "Required", rules: rules, fields: fields })
            }
          }
        }
        if (this.podSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
          const schedulings = this.podSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution
          for (const s of schedulings) {
            let rules = []
            let fields = []
            if (s.preference) {
              if (s.preference.matchExpressions) {
                for (const express of s.preference.matchExpressions) {
                  rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
              if (s.preference.matchFields) {
                for (const express of s.preference.matchFields) {
                  fields.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
            }
            this.nodeSchedulings.push({ priority: "Preferred", rules: rules, fields: fields, weight: s.weight || null })
          }
        }
      }
    }
  },
}
</script>