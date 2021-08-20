<template>
  <div style="margin-top: 20px">
    <ko-card :title="'Pod ' + $t('business.workload.schedule')">
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
                <el-form-item :label="$t('business.workload.namespace_operation')">
                  <ko-form-item itemType="select" v-model="item.namespaceOperation" :selections="namespace_operation_list" />
                </el-form-item>
              </el-col>
              <el-col :span="12" v-if="item.namespaceOperation === 'selectNamespace'">
                <el-form-item label="Namespace">
                  <ko-form-item itemType="select2" v-model="item.namespaces" multiple :selections="namespace_list" />
                </el-form-item>
              </el-col>
            </el-row>

            <table style="width: 98%" class="tab-table">
              <tr>
                <th scope="col" width="40%" align="left">
                  <label>{{$t('business.workload.key')}}</label>
                </th>
                <th scope="col" width="15%" align="left">
                  <label>{{$t('business.workload.operator')}}</label>
                </th>
                <th scope="col" width="40%" align="left">
                  <label>{{$t('business.workload.value')}}</label>
                </th>
                <th align="left"></th>
              </tr>
              <tr v-for="(row, index) in item.rules" v-bind:key="index">
                <td>
                  <ko-form-item itemType="input" v-model="row.key" />
                </td>
                <td>
                  <ko-form-item itemType="select" v-model="row.operator" :selections="operator_list" />
                </td>
                <td>
                  <ko-form-item v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                  <ko-form-item v-else itemType="input" v-model="row.value" />
                </td>
                <td>
                  <el-button type="text" style="font-size: 10px" @click="handleMatchDelete(item, index)">
                    {{ $t("commons.button.delete") }}
                  </el-button>
                </td>
              </tr>
              <tr>
                <td align="left">
                  <el-button @click="handleMatchAdd(item)">{{ $t("commons.button.add") }}</el-button>
                </td>
              </tr>
            </table>
            <el-row style="margin-top: 10px">
              <el-col :span=24>
                <el-form-item :label="$t('business.workload.topology_key')">
                  <ko-form-item itemType="input" placeholder="e.g. failure-domain.beta.kubernetes.io/zone" v-model="item.topologyKey" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-card>
        </div>
        <el-button @click="handlePodRulesAdd()">{{$t('business.workload.add')}}{{$t('business.workload.node_selector')}}</el-button>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoPodScheduling",
  components: { KoFormItem, KoCard },
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
    }
  },
  methods: {
    handlePodRulesAdd() {
      var item = {
        type: "Affinity",
        priority: "Preferred",
        namespaceOperation: "podNamespace",
        namespaces: "",
        rules: [],
        topologyKey: "",
      }
      this.podSchedulings.push(item)
    },
    handlePodRulesDelete(index) {
      this.podSchedulings.splice(index, 1)
    },

    handleMatchAdd(schedulingItem) {
      var item = {
        key: "",
        operator: "",
        value: "",
      }
      schedulingItem.rules.push(item)
    },
    handleMatchDelete(schedulingItem, index) {
      schedulingItem.rules.splice(index, 1)
    },
    valueTrans(type, priority, s) {
      let namespaceOperation,
        namespaces = ""
      if (s.namespaces) {
        namespaceOperation = "selectNamespace"
        namespaces = s.namespaces
      } else {
        namespaceOperation = "podNamespace"
      }
      let rules = []
      if (s.labelSelector.matchExpressions) {
        for (const express of s.matchExpressions) {
          rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
        }
      }
      const topologyKey = s.topologyKey ? s.topologyKey : ""
      this.podSchedulings.push({
        type: type,
        priority: priority,
        namespaceOperation: namespaceOperation,
        namespaces: namespaces,
        rules: rules,
        topologyKey: topologyKey,
      })
    },
    getMatchExpress(schedule) {
      let matchs = []
      if (schedule.rules.length === 0) {
        return matchs
      }
      for (const rule of schedule.rules) {
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

    transformation(parentFrom) {
      parentFrom.affinity = {}
      if (this.podSchedulings.length !== 0) {
        for (const pS of this.podSchedulings) {
          let itemAdd = {}
          if (pS.namespaceOperation === "selectNamespace" && pS.namespaces.length !== 0) {
            itemAdd.namespaces = pS.namespaces
          }
          if (pS.topologyKey) {
            itemAdd.topologyKey = pS.topologyKey
          }
          const matchs = this.getMatchExpress(pS)
          switch (pS.type + "+" + pS.priority) {
            case "Affinity+Required":
              if (!parentFrom.affinity.podAffinity) {
                parentFrom.affinity.podAffinity = {}
              }
              if (!parentFrom.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                parentFrom.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
              }
              itemAdd.labelSelector = {}
              if (matchs.length !== 0) {
                itemAdd.labelSelector.matchExpressions = matchs
              }
              parentFrom.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Affinity+Preferred":
              if (!parentFrom.affinity.podAffinity) {
                parentFrom.affinity.podAffinity = {}
              }
              if (!parentFrom.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
                parentFrom.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
              }
              itemAdd.podAffinityTerm = {}
              itemAdd.weight = 1
              itemAdd.podAffinityTerm.labelSelector = {}
              if (matchs.length !== 0) {
                itemAdd.podAffinityTerm.labelSelector.matchExpressions = matchs
              }
              parentFrom.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Anti-Affinity+Required":
              if (!parentFrom.affinity.podAntiAffinity) {
                parentFrom.affinity.podAntiAffinity = {}
              }
              if (!parentFrom.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                parentFrom.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
              }
              itemAdd._anti = true
              itemAdd.labelSelector = {}
              if (matchs.length !== 0) {
                itemAdd.labelSelector.matchExpressions = matchs
              }
              parentFrom.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
              break
            case "Anti-Affinity+Preferred":
              if (!parentFrom.affinity.podAntiAffinity) {
                parentFrom.affinity.podAntiAffinity = {}
              }
              if (!parentFrom.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
                parentFrom.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
              }
              itemAdd._anti = true
              itemAdd.podAffinityTerm = {}
              itemAdd.weight = 1
              itemAdd.podAffinityTerm.labelSelector = {}
              if (matchs.length !== 0) {
                itemAdd.podAffinityTerm.labelSelector.matchExpressions = matchs
              }
              parentFrom.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
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
    }
  },
}
</script>