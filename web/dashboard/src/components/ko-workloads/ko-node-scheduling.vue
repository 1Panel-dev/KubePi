<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.node_schedule')">
      <el-form label-position="top" ref="form" :disabled="isReadOnly">
        <el-row>
          <el-col :span="24">
            <el-form-item :label="$t('business.workload.schedule_type')" v-if="enableSchedulingList">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="scheduling_type" @change="changeType" :radios="scheduling_type_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row v-if="scheduling_type === 'specific_node'">
          <el-col :span="24">
            <el-form-item :label="$t('business.workload.node_name')">
              <ko-form-item itemType="select2" v-model="nodeName" :selections="node_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <div v-if="scheduling_type === 'matching_rules'">
          <div v-for="(item, index) in nodeSchedulings" :key="index">
            <el-card style="margin-top: 10px">
              <el-row>
                <el-button style="float: right; padding: 3px 0" type="text" @click="handleNodeRulesDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </el-row>
              <el-row>
                <el-col :span="12">
                  <el-form-item :label="$t('business.workload.priority')" v-if="enablePrioritySelect">
                    <ko-form-item itemType="radio" v-model="item.priority" :radios="priority_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12" v-if="item.weight && item.priority === 'Preferred'">
                  <el-form-item :label="$t('business.workload.weight')">
                    <ko-form-item itemType="number" v-model="item.weight" />
                  </el-form-item>
                </el-col>
              </el-row>

              <table style="width: 98%" class="tab-table">
                <tr>
                  <th scope="col" width="40%" align="left">
                    <label>{{ $t('business.workload.key') }}</label>
                  </th>
                  <th scope="col" width="15%" align="left">
                    <label>{{ $t('business.workload.operator') }}</label>
                  </th>
                  <th scope="col" width="40%" align="left">
                    <label>{{ $t('business.workload.value') }}</label>
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
                    <el-button @click="handleMatchAdd(item)">
                      {{ $t('business.workload.add') }}{{ $t('business.workload.rule') }}
                    </el-button>
                  </td>
                </tr>
              </table>
            </el-card>
          </div>
        </div>
        <el-button v-if="scheduling_type === 'matching_rules'" @click="handleNodeRulesAdd()">
          {{ $t('business.workload.add') }}{{ $t('business.workload.node_selector') }}
        </el-button>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"
import { listNodes } from "@/api/nodes"

export default {
  name: "KoNodeScheduling",
  components: { KoFormItem, KoCard },
  props: {
    nodeSchedulingParentObj: Object,
    isReadOnly: Boolean,
    nodeSchedulingType: String,
  },
  data() {
    return {
      scheduling_type_list: [
        { label: this.$t("business.workload.any_node"), value: "any_node" },
        { label: this.$t("business.workload.specific_node"), value: "specific_node" },
        { label: this.$t("business.workload.match_rules"), value: "matching_rules" },
      ],
      priority_list: [
        { label: this.$t("business.workload.preferred"), value: "Preferred" },
        { label: this.$t("business.workload.required"), value: "Required" },
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
      enablePrioritySelect: true,
      priority: "Preferred",
    }
  },
  methods: {
    handleNodeRulesAdd() {
      var item = {
        priority: this.priority,
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

    changeType() {
      if (this.scheduling_type === "specific_node") {
        this.node_list = []
        listNodes(this.clusterName).then((res) => {
          for (const node of res.items) {
            this.node_list.push(node.metadata.name)
          }
        })
      }
    },

    transformation(parentFrom) {
      if (this.scheduling_type === "specific_node") {
        parentFrom.nodeName = this.nodeName || undefined
      }
      if (this.scheduling_type === "matching_rules") {
        parentFrom.nodeAffinity = {}
        if (this.nodeSchedulings.length !== 0) {
          for (const nS of this.nodeSchedulings) {
            const matchs = this.getMatchExpress(nS)
            let itemAdd = {}
            switch (nS.priority) {
              case "Preferred":
                parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution = []
                itemAdd.weight = 1
                itemAdd.preference = {}
                itemAdd.preference.matchExpressions = matchs
                parentFrom.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution.push(itemAdd)
                break
              case "Required":
                if (!parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
                  parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution = {}
                }
                parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms = []
                itemAdd.matchExpressions = matchs
                parentFrom.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms.push(itemAdd)
                break
              case "None":
                parentFrom.nodeAffinity.required = {
                  nodeSelectorTerms: [],
                }
                itemAdd.matchExpressions = matchs
                parentFrom.nodeAffinity.required.nodeSelectorTerms.push(itemAdd)
                break
            }
          }
        }
      }
    },
  },
  mounted() {
    this.clusterName = this.$route.query.cluster
    this.nodeSchedulings = []
    if (this.nodeSchedulingParentObj) {
      if (this.nodeSchedulingParentObj.nodeAffinity) {
        this.scheduling_type = "matching_rules"
        if (this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution) {
          if (this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms) {
            const schedulings = this.nodeSchedulingParentObj.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms
            for (const s of schedulings) {
              let rules = []
              if (s.matchExpressions) {
                for (const express of s.matchExpressions) {
                  rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
              this.nodeSchedulings.push({ type: "Node", priority: "Required", rules: rules })
            }
            this.priority = "Required"
          }
        }
        if (this.nodeSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution) {
          const schedulings = this.nodeSchedulingParentObj.nodeAffinity.preferredDuringSchedulingIgnoredDuringExecution
          for (const s of schedulings) {
            let rules = []
            if (s.preference) {
              if (s.preference.matchExpressions) {
                for (const express of s.preference.matchExpressions) {
                  rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
                }
              }
            }
            const weight = s.weight ? s.weight : 1
            this.nodeSchedulings.push({ type: "Node", priority: "Preferred", weight: weight, rules: rules })
          }
          this.priority = "Preferred"
        }
        if (this.nodeSchedulingParentObj.nodeAffinity.required) {
          const scheduling = this.nodeSchedulingParentObj.nodeAffinity.required.nodeSelectorTerms
          for (const s of scheduling) {
            let rules = []
            if (s.matchExpressions) {
              for (const express of s.matchExpressions) {
                rules.push({ key: express.key, operator: express.operator, value: express.values.join(",") })
              }
            }
            this.nodeSchedulings.push({ type: "Node", priority: "None", rules: rules })
          }
          this.priority = "None"
        }
      }
    }
    if (this.nodeSchedulingType) {
      this.scheduling_type = this.nodeSchedulingType
      this.enableSchedulingList = false
      this.enablePrioritySelect = false
    }
  },
}
</script>
