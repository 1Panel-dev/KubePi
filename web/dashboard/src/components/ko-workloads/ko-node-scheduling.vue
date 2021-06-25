<template>
  <div>
    <el-row :gutter="20">
      <el-col :span="12">
        <el-divider content-position="left">Pod Scheduling</el-divider>
        <div v-for="(item, index) in podSchedulings" :key="index">
          <el-card style="margin-top: 10px">
            <el-button style="float: right; padding: 3px 0" type="text" @click="handlePodRulesDelete(index)">删 除</el-button>
            <el-row style="margin-top: 20px">
              <el-col :span="12">
                <ko-form-item labelName="Type" clearable itemType="radio" v-model="item.type" :radios="type_list" />
              </el-col>
              <el-col :span="12">
                <ko-form-item labelName="Priority" clearable itemType="radio" v-model="item.priority" :radios="priority_list" />
              </el-col>
            </el-row>
            <el-row style="margin-top: 20px">
              <el-col :span="12">
                <ko-form-item labelName="Namespace" clearable itemType="select" v-model="item.namespaceOperation" :selections="namespace_operation_list" />
              </el-col>
              <el-col :span="12" v-if="item.namespaceOperation === 'selectNamespace'">
                <ko-form-item labelName="Namespace" clearable itemType="select" v-model="item.namespaces" multiple :selections="namespace_list" />
              </el-col>
            </el-row>

            <el-button style="margin-top: 20px" @click="handleMatchAdd(item)">Add</el-button>
            <el-table :data="item.rules">
              <el-table-column min-width="40" label="Key">
                <template v-slot:default="{row}">
                  <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.key" />
                </template>
              </el-table-column>
              <el-table-column min-width="30" label="Operator">
                <template v-slot:default="{row}">
                  <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.operator" :selections="operator_list" />
                </template>
              </el-table-column>
              <el-table-column min-width="40" label="Value">
                <template v-slot:default="{row}">
                  <ko-form-item :withoutLabel="true" v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                  <ko-form-item :withoutLabel="true" v-else clearable itemType="input" v-model="row.value" />
                </template>
              </el-table-column>
              <el-table-column width="60px">
                <template v-slot:default="{row}">
                  <el-button type="text" style="font-size: 10px" @click="handleMatchDelete(item, row)">{{ $t("commons.button.delete") }}</el-button>
                </template>
              </el-table-column>
            </el-table>
            <el-row style="margin-top: 20px">
              <el-col :span=24>
                <ko-form-item labelName="Topology Key" clearable itemType="input" placeholder="e.g. failure-domain.beta.kubernetes.io/zone" v-model="item.topologyKey" />
              </el-col>
            </el-row>
          </el-card>
        </div>
        <el-button style="margin-top: 20px" @click="handlePodRulesAdd()">Add Node Selector</el-button>
      </el-col>

      <el-col :span="12">
        <el-divider content-position="left">Node Scheduling</el-divider>
        <el-row style="margin-top: 20px">
          <el-col :span="24">
            <ko-form-item labelName="Scheduling Type" radioLayout="vertical" clearable itemType="radio" v-model="scheduling_type" :radios="scheduling_type_list" />
          </el-col>
        </el-row>
        <el-row style="margin-top: 20px" v-if="scheduling_type === 'specific_node'">
          <el-col :span="24">
            <ko-form-item labelName="Node Name" clearable itemType="select" v-model="nodeName" :selections="node_list" />
          </el-col>
        </el-row>
        <div v-if="scheduling_type === 'matching_rules'">
          <div v-for="(item, index) in nodeSchedulings" :key="index">
            <el-card style="margin-top: 10px">
              <el-button style="float: right; padding: 3px 0" type="text" @click="handleNodeRulesDelete(index)">删 除</el-button>
              <el-row style="margin-top: 20px">
                <el-col :span="12">
                  <ko-form-item labelName="Priority" clearable itemType="radio" v-model="item.priority" :radios="priority_list" />
                </el-col>
                <el-col :span="12" v-if="item.weight && item.priority === 'Preferred'">
                  <ko-form-item labelName="Weight" clearable itemType="number" v-model="item.weight" />
                </el-col>
              </el-row>

              <el-button style="margin-top: 20px" @click="handleMatchAdd(item)">Add Rule</el-button>
              <el-table :data="item.rules">
                <el-table-column min-width="40" label="Key">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.key" />
                  </template>
                </el-table-column>
                <el-table-column min-width="30" label="Operator">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.operator" :selections="operator_list" />
                  </template>
                </el-table-column>
                <el-table-column min-width="40" label="Value">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                    <ko-form-item :withoutLabel="true" v-else clearable itemType="input" v-model="row.value" />
                  </template>
                </el-table-column>
                <el-table-column width="60px">
                  <template v-slot:default="{row}">
                    <el-button type="text" style="font-size: 10px" @click="handleMatchDelete(item, row)">{{ $t("commons.button.delete") }}</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-card>
          </div>
        </div>
        <el-button v-if="scheduling_type === 'matching_rules'" style="margin-top: 20px" @click="handleNodeRulesAdd()">Add Node Selector</el-button>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-divider content-position="left">Tolerations</el-divider>
      <el-button @click="handleTolerationsAdd()">Add Toleration</el-button>
      <el-table v-if="tolerations.length !== 0" :data="tolerations">
        <el-table-column min-width="40" label="Label Key">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.key" />
          </template>
        </el-table-column>
        <el-table-column min-width="15" label="Operator">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.operator" :selections="tolerations_operator_list" />
          </template>
        </el-table-column>
        <el-table-column min-width="40" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.value" />
          </template>
        </el-table-column>
        <el-table-column min-width="25" label="Effect">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.effect" :selections="effect_list" />
          </template>
        </el-table-column>
        <el-table-column min-width="20" label="Toleration Seconds(s)">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" :disabled="row.effect !== 'NoExecute'" clearable itemType="number" v-model.number="row.tolerationSeconds" />
          </template>
        </el-table-column>
        <el-table-column width="60px">
          <template v-slot:default="{row}">
            <el-button type="text" style="font-size: 10px" @click="handleTolerationsDelete(row)">{{ $t("commons.button.delete") }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoNodeScheduling",
  components: { KoFormItem },
  props: {
    schedulingParentObj: Object,
  },
  data() {
    return {
      scheduling_type_list: [
        { label: "Run pods on any avaliable node", value: "any_node" },
        { label: "Run pods on specific node(s)", value: "specific_node" },
        { label: "Run pods on node(s) matching scheduling rules", value: "matching_rules" },
      ],
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
      tolerations_operator_list: [
        { label: "Exists", value: "Exists" },
        { label: "=", value: "Equal" },
      ],
      effect_list: [
        { label: "All", value: "All" },
        { label: "NoSchedule", value: "NoSchedule" },
        { label: "PreferNoSchedule", value: "PreferNoSchedule" },
        { label: "NoExecute", value: "NoExecute" },
      ],
      namespace_list: [
        { label: "kube-system", value: "kube-system" },
        { label: "kube-public", value: "kube-public" },
        { label: "kube-operator", value: "kube-operator" },
        { label: "default", value: "default" },
      ],
      podSchedulings: [],
      scheduling_type: "any_node",
      nodeName: "",
      node_list: [],
      nodeSchedulings: [],
      tolerations: [],
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

    handleNodeRulesAdd() {
      var item = {
        priority: "Preferred",
        rules: [],
      }
      this.nodeSchedulings.push(item)
    },
    handleNodeRulesDelete(index) {
      this.nodeSchedulings.splice(index, 1)
    },

    handleTolerationsAdd() {
      var item = {
        key: "",
        operator: "",
        value: "",
        effect: "",
        tolerationSeconds: "",
      }
      this.tolerations.unshift(item)
    },
    handleTolerationsDelete(row) {
      for (let i = 0; i < this.tolerations.length; i++) {
        if (this.tolerations[i] === row) {
          this.tolerations.splice(i, 1)
        }
      }
    },

    handleMatchAdd(schedulingItem) {
      var item = {
        key: "",
        operator: "",
        value: "",
      }
      schedulingItem.rules.push(item)
    },
    handleMatchDelete(schedulingItem, row) {
      for (let i = 0; i < schedulingItem.rules.length; i++) {
        if (schedulingItem.rules[i] === row) {
          schedulingItem.rules.splice(i, 1)
        }
      }
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
      if (this.tolerations.length !== 0) {
        parentFrom.tolerations = this.tolerations
      }
      if (this.scheduling_type === "specific_node") {
        parentFrom.nodeName = this.nodeName
      }
      if (this.scheduling_type === "matching_rules") {
        if (this.nodeSchedulings.length !== 0) {
          for (const nS of this.podSchedulings) {
            const matchs = this.getMatchExpress(nS)
            if (matchs.length !== 0) {
              let itemAdd = {}
              switch (nS.priority) {
                case "Preferred":
                  if (!parentFrom.nodeAffinity) {
                    parentFrom.nodeAffinity = {}
                  }
                  if (!parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
                    parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
                  }
                  itemAdd.weight = 1
                  itemAdd.preference = {}
                  itemAdd.preference.matchExpressions = matchs
                  parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                  break
                case "Required":
                  if (!parentFrom.nodeAffinity) {
                    parentFrom.nodeAffinity = {}
                  }
                  if (!parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                    parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution = {}
                  }
                  if (!parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms) {
                    parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms = []
                  }
                  parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.push(itemAdd)
                  break
              }
            }
          }
        }
      }
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
    this.nodeSchedulings = []
    this.podSchedulings = []
    if (this.schedulingParentObj) {
      if (this.schedulingParentObj.affinity) {
        if (this.schedulingParentObj.affinity.nodeAffinity) {
          this.scheduling_type = "matching_rules"
          if (this.schedulingParentObj.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            if (this.schedulingParentObj.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms) {
              const schedulings = this.schedulingParentObj.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms
              for (const s of schedulings) {
                let rules = []
                for (const express of s.matchExpressions) {
                  rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
                this.nodeSchedulings.push({ type: "Node", priority: "Required", rules: rules })
              }
            }
          }
          if (this.schedulingParentObj.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.schedulingParentObj.affinity.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              let rules = []
              for (const express of s.preference.matchExpressions) {
                rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
              }
              console.log(s)
              const weight = s.weight ? s.weight : 1
              this.nodeSchedulings.push({ type: "Node", priority: "Preferred", weight: weight, rules: rules })
            }
          }
        }

        if (this.schedulingParentObj.affinity.podAffinity) {
          if (this.schedulingParentObj.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.schedulingParentObj.affinity.podAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Affinity", "Required", s)
            }
          }
          if (this.schedulingParentObj.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.schedulingParentObj.affinity.podAffinity.preferredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              if (s.podAffinityTerm) {
                this.valueTrans("Affinity", "Preferred", s.podAffinityTerm)
              }
            }
          }
        }
        if (this.schedulingParentObj.affinity.podAntiAffinity) {
          if (this.schedulingParentObj.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.schedulingParentObj.affinity.podAntiAffinity.requiredDuringSchedulingIgnoredDuringExecution
            for (const s of schedulings) {
              this.valueTrans("Anti-Affinity", "Required", s)
            }
          }
          if (this.schedulingParentObj.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
            const schedulings = this.schedulingParentObj.affinity.podAntiAffinity.preferredDuringSchedulingIgnoredDuringExecution
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