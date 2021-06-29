<template>
  <div style="margin-top: 20px">
    <ko-card title="Pod Scheduling">
      <el-form label-position="top">
        <div v-for="(item, index) in podSchedulings" :key="index">
          <el-card style="margin-top: 10px">
            <el-row>
              <el-button style="float: right;" type="text" @click="handlePodRulesDelete(index)">删 除</el-button>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Type">
                  <ko-form-item itemType="radio" v-model="item.type" :radios="type_list" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Priority">
                  <ko-form-item itemType="radio" v-model="item.priority" :radios="priority_list" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Namespace Operation">
                  <ko-form-item itemType="select" v-model="item.namespaceOperation" :selections="namespace_operation_list" />
                </el-form-item>
              </el-col>
              <el-col :span="12" v-if="item.namespaceOperation === 'selectNamespace'">
                <el-form-item label="Namespace">
                  <ko-form-item itemType="select" v-model="item.namespaces" multiple :selections="namespace_list" />
                </el-form-item>
              </el-col>
            </el-row>

            <table style="width: 98%" class="tab-table">
              <tr>
                <th scope="col" width="40%" align="left">
                  <label>key</label>
                </th>
                <th scope="col" width="15%" align="left">
                  <label>operator</label>
                </th>
                <th scope="col" width="40%" align="left">
                  <label>value</label>
                </th>
                <th align="left"></th>
              </tr>
              <tr v-for="(row, index) in item.rules" v-bind:key="index">
                <td>
                  <ko-form-item :withoutLabel="true" itemType="input" v-model="row.key" />
                </td>
                <td>
                  <ko-form-item :withoutLabel="true" itemType="select" v-model="row.operator" :selections="operator_list" />
                </td>
                <td>
                  <ko-form-item :withoutLabel="true" v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                  <ko-form-item :withoutLabel="true" v-else itemType="input" v-model="row.value" />
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
                <el-form-item label="Topology Key">
                  <ko-form-item itemType="input" placeholder="e.g. failure-domain.beta.kubernetes.io/zone" v-model="item.topologyKey" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-card>
        </div>
      </el-form>
      <el-button @click="handlePodRulesAdd()">Add Node Selector</el-button>
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
    podScheduling: Object,
  },
  data() {
    return {
      type_list: [
        { label: "Affinity", value: "Affinity" },
        { label: "Anti-Affinity", value: "Anti-Affinity" },
      ],
      priority_list: [
        { label: "Preferred", value: "Preferred" },
        { label: "Required", value: "Required" },
      ],
      namespace_operation_list: [
        { label: "This pod's namespace", value: "podNamespace" },
        { label: "Pods in these namespaces", value: "selectNamespace" },
      ],
      operator_list: [
        { label: "<", value: "Lt" },
        { label: ">", value: "Gt" },
        { label: "Exists", value: "Exists" },
        { label: "Not Exist", value: "DoesNotExist" },
        { label: "=", value: "In" },
        { label: "!=", value: "NotIn" },
      ],
      namespace_list: [
        { label: "kube-system", value: "kube-system" },
        { label: "kube-public", value: "kube-public" },
        { label: "kube-operator", value: "kube-operator" },
        { label: "default", value: "default" },
      ],
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
      for (const express of s.matchExpressions) {
        rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
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
      if (this.podSchedulings.length !== 0) {
        for (const pS of this.podSchedulings) {
          const matchs = this.getMatchExpress(pS)
          if (matchs.length !== 0) {
            let itemAdd = {}
            if (pS.namespaceOperation === "selectNamespace" && pS.namespaces.length !== 0) {
              itemAdd.namespaces = pS.namespaces
            }
            itemAdd.topologyKey = pS.topologyKey
            switch (pS.type + "+" + pS.priority) {
              case "Affinity+Required":
                if (!parentFrom.podAffinity) {
                  parentFrom.podAffinity = {}
                }
                if (!parentFrom.podAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                  parentFrom.podAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
                }
                itemAdd.labelSelector = {}
                itemAdd.labelSelector.matchExpressions = matchs
                parentFrom.podAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                break
              case "Affinity+Preferred":
                if (!parentFrom.podAffinity) {
                  parentFrom.podAffinity = {}
                }
                if (!parentFrom.podAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
                  parentFrom.podAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
                }
                itemAdd.podAffinityTerm = {}
                itemAdd.weight = 1
                itemAdd.podAffinityTerm.labelSelector = {}
                itemAdd.podAffinityTerm.labelSelector.matchExpressions = matchs
                parentFrom.podAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                break
              case "Anti-Affinity+Required":
                if (!parentFrom.podAntiAffinity) {
                  parentFrom.podAntiAffinity = {}
                }
                if (!parentFrom.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                  parentFrom.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution = []
                }
                itemAdd._anti = true
                itemAdd.labelSelector = {}
                itemAdd.labelSelector.matchExpressions = matchs
                parentFrom.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                break
              case "Anti-Affinity+Preferred":
                if (!parentFrom.podAntiAffinity) {
                  parentFrom.podAntiAffinity = {}
                }
                if (!parentFrom.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
                  parentFrom.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
                }
                itemAdd._anti = true
                itemAdd.podAffinityTerm = {}
                itemAdd.weight = 1
                itemAdd.podAffinityTerm.labelSelector = {}
                itemAdd.podAffinityTerm.labelSelector.matchExpressions = matchs
                parentFrom.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                break
            }
          }
        }
      }
    },
  },
  mounted() {
    this.podSchedulings = []
    if (this.podScheduling) {
      if (this.podScheduling.affinity) {
        if (this.podScheduling.affinity.podAffinity) {
          if (this.podScheduling.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podScheduling.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Affinity", "Required", s)
            }
          }
          if (this.podScheduling.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podScheduling.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              if (s.podAffinityTerm) {
                this.valueTrans("Affinity", "Preferred", s.podAffinityTerm)
              }
            }
          }
        }
        if (this.podScheduling.affinity.podAntiAffinity) {
          if (this.podScheduling.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podScheduling.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Anti-Affinity", "Required", s)
            }
          }
          if (this.podScheduling.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.podScheduling.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
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