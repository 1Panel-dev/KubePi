<template>
  <div style="margin-top: 20px">
    <ko-card title="Node Scheduling">
      <el-form label-position="top" ref="form" :disabled="isReadOnly">
        <el-row>
          <el-col :span="24">
            <el-form-item label="Scheduling Type" v-if="enableSchedulingList">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="scheduling_type" :radios="scheduling_type_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="scheduling_type === 'specific_node'">
          <el-col :span="24">
            <el-form-item label="Node Name">
              <ko-form-item itemType="select2" v-model="nodeName" :selections="node_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <div v-if="scheduling_type === 'matching_rules'">
          <div v-for="(item, index) in nodeSchedulings" :key="index">
            <el-card style="margin-top: 10px">
              <el-row>
                <el-button style="float: right; padding: 3px 0" type="text" @click="handleNodeRulesDelete(index)">删 除</el-button>
              </el-row>
              <el-row>
                <el-col :span="12">
                  <el-form-item label="Priority" v-if="enablePrioritySelect">
                    <ko-form-item itemType="radio" v-model="item.priority" :radios="priority_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12" v-if="item.weight && item.priority === 'Preferred'">
                  <el-form-item label="Weight">
                    <ko-form-item itemType="number" v-model="item.weight" />
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
                    <el-button @click="handleMatchAdd(item)">Add Rule</el-button>
                  </td>
                </tr>
              </table>
            </el-card>
          </div>
        </div>
        <el-button v-if="scheduling_type === 'matching_rules'" @click="handleNodeRulesAdd()">Add Node Selector</el-button>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoNodeScheduling",
  components: { KoFormItem, KoCard },
  props: {
    nodeSchedulingParentObj: Object,
    nodeList: Array,
    isReadOnly: Boolean,
    nodeSchedulingType: String,
  },
  watch: {
    nodeList: {
      handler(newName) {
        this.node_list = []
        for (const node of newName) {
          this.node_list.push(node.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      scheduling_type_list: [
        { label: "Run pods on any avaliable node", value: "any_node" },
        { label: "Run pods on specific node(s)", value: "specific_node" },
        { label: "Run pods on node(s) matching scheduling rules", value: "matching_rules" },
      ],
      priority_list: [
        { label: "Preferred", value: "Preferred" },
        { label: "Required", value: "Required" },
      ],
      operator_list: [
        { label: "<", value: "Lt" },
        { label: ">", value: "Gt" },
        { label: "Exists", value: "Exists" },
        { label: "Not Exist", value: "DoesNotExist" },
        { label: "=", value: "In" },
        { label: "!=", value: "NotIn" },
      ],
      scheduling_type: "any_node",
      nodeName: "",
      node_list: [],
      nodeSchedulings: [],
      enableSchedulingList: true,
      enablePrioritySelect: true
    }
  },
  methods: {
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
      if (this.scheduling_type === "specific_node") {
        parentFrom.nodeName = this.nodeName
      }
      if (this.scheduling_type === "matching_rules") {
        if (this.nodeSchedulings.length !== 0) {
          for (const nS of this.nodeSchedulings) {
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
                  itemAdd.preference = {}
                  itemAdd.preference.matchExpressions = matchs
                  parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.push(itemAdd)
                  break
                case "None":
                  if (!parentFrom.nodeAffinity) {
                    parentFrom.nodeAffinity = {}
                  }
                  itemAdd.matchExpressions = matchs
                  parentFrom.nodeAffinity.required.nodeSelectorTerms.push(itemAdd)
                  break
              }
            }
          }
        }
      }
    },
  },
  mounted() {
    this.nodeSchedulings = []
    if (this.nodeSchedulingParentObj) {
      if (this.nodeSchedulingParentObj.nodeAffinity) {
        this.scheduling_type = "matching_rules"
        if (this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
          if (this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms) {
            const schedulings = this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms
            for (const s of schedulings) {
              let rules = []
              for (const express of s.preference.matchExpressions) {
                rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
              }
              this.nodeSchedulings.push({ type: "Node", priority: "Required", rules: rules })
            }
          }
        }
        if (this.nodeSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
          const schedulings = this.nodeSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution
          for (const s of schedulings) {
            let rules = []
            for (const express of s.preference.matchExpressions) {
              rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
            }
            const weight = s.weight ? s.weight : 1
            this.nodeSchedulings.push({ type: "Node", priority: "Preferred", weight: weight, rules: rules })
          }
        }
        if (this.nodeSchedulingParentObj.nodeAffinity.required) {
          const scheduling = this.nodeSchedulingParentObj.nodeAffinity.required.nodeSelectorTerms
          for (const s of scheduling) {
            let rules = []
            for (const express of s.matchExpressions) {
              rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
            }
            this.nodeSchedulings.push({ type: "Node", priority: "None", rules: rules })
          }
        }
      }
    }
    if (this.nodeSchedulingType !== "" && typeof this.nodeSchedulingType !== "undefined") {
      this.scheduling_type = this.nodeSchedulingType
      this.enableSchedulingList = false
      this.enablePrioritySelect = false
    }
  },
}
</script>