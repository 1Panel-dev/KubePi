<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Node Scheduling" clearable itemType="radio" v-model="form.scheduling_type" :radios="scheduling_type_list" />
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Node Name" clearable itemType="select" v-model="form.nodeName" :selections="node_list" />
      </el-col>
    </el-row>

    <div v-if="form.scheduling_type==='matching_rules'">
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <label>Require any of:</label>
          <div style="margin-top: 10px">
            <el-button @click="handleSelectorRequiredAdd()">Add Node Selector</el-button>
          </div>
          <div v-for="(item, index) in nodeSelectorTerms" :key="index">
            <div style="border: solid 1px rgb(228, 204, 204); border-radius: 5px;padding: 5px; margin-top: 5px;">
              <div>
                <el-button type="text" style="float: right;margin-right: 20px; font-size: 16px" @click="handleSelectorRequiredClose(nodeSelectorTerms[index])">X</el-button>
                <el-button @click="handleRuleRequiredAdd(nodeSelectorTerms[index])">Add</el-button>
              </div>
              <el-table :data="item.matchExpressions">
                <el-table-column min-width="40" label="Key">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.key" />
                  </template>
                </el-table-column>
                <el-table-column min-width="35" label="Operator">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.operator" :selections="operator_list" />
                  </template>
                </el-table-column>
                <el-table-column min-width="40" label="Value">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                    <ko-form-item :withoutLabel="true" v-else clearable itemType="input" v-model="row.values" />
                  </template>
                </el-table-column>
                <el-table-column width="120px">
                  <template v-slot:default="{row}">
                    <el-button type="text" style="font-size: 20px" @click="handleRuleRequiredDelete(nodeSelectorTerms[index], row)">REMOVE</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-col>
        <el-col :span="12">
          <label>Prefer any of:</label>
          <div style="margin-top: 10px">
            <el-button @click="handleSelectorPreferAdd()">Add Node Selector</el-button>
          </div>
          <div v-for="(item, index) in preference" :key="index">
            <div style="border: solid 1px rgb(228, 204, 204); border-radius: 5px;padding: 5px; margin-top: 5px;">
              <div>
                <el-button type="text" style="float: right;margin-right: 20px; font-size: 16px" @click="handleSelectorPreferClose(preference[index])">X</el-button>
                <el-button @click="handleRulePreferAdd(preference[index])">Add</el-button>
              </div>
              <el-table :data="item.matchExpressions">
                <el-table-column min-width="40" label="Key">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.name" />
                  </template>
                </el-table-column>
                <el-table-column min-width="35" label="Operator">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.operator" :selections="operator_list"/>
                  </template>
                </el-table-column>
                <el-table-column min-width="40" label="Value">
                  <template v-slot:default="{row}">
                    <ko-form-item :withoutLabel="true" v-if="row.operator === 'Exists' || row.operator === 'DoesNotExist'" disabled itemType="input" value="N/A" />
                    <ko-form-item :withoutLabel="true" v-else clearable itemType="input" v-model="row.values" />
                  </template>
                </el-table-column>
                <el-table-column width="120px">
                  <template v-slot:default="{row}">
                    <el-button type="text" style="font-size: 20px" @click="handleRulePreferDelete(preference[index], row)">REMOVE</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <div style="margin-top: 20px" v-if="form.scheduling_type==='specific_node'">
      <label>Nodes with these labels
        <el-tooltip class="item" effect="dark" content="ProTip: Paste lines of key=value or key: value into any key field for easy bulk entry" placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>
      </label>
      <div style="margin-top: 10px">
        <el-button @click="handleAdd">Add</el-button>
      </div>
      <el-table v-if="form.nodeSelector.length !== 0" :data="form.nodeSelector">
        <el-table-column min-width="80" label="Key">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo" clearable itemType="input" v-model="row.name" />
          </template>
        </el-table-column>
        <el-table-column min-width="80" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" placeholder="e.g. bar" clearable itemType="input" v-model="row.value" />
          </template>
        </el-table-column>
        <el-table-column width="120px">
          <template v-slot:default="{row}">
            <el-button type="text" style="font-size: 20px" @click="handleDelete(row)">REMOVE</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoNodeScheduling",
  components: { KoFormItem },
  data() {
    return {
      scheduling_type_list: [
        { label: "Run pods on specific node(s)", value: "specific_node" },
        { label: "Run pods on node(s) matching these scheduling rules", value: "matching_rules" },
      ],
      operator_list: [
        { label: "<", value: "Lt" },
        { label: ">", value: "Gt" },
        { label: "Exists", value: "Exists" },
        { label: "Not Exist", value: "DoesNotExist" },
        { label: "=", value: "In" },
        { label: "!=", value: "NotIn" },
      ],
      node_list: [],
      nodeSelectorTerms: [],
      preference: [],
      form: {
        scheduling_type: "specific_node",
        nodeName: "",
        nodeSelector: [],
      },
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.form.nodeSelector.length; i++) {
        if (this.form.nodeSelector[i] === row) {
          this.form.nodeSelector.splice(i, 1)
        }
      }
    },
    handleAdd() {
      var item = {
        name: "",
        value: "",
      }
      this.form.nodeSelector.unshift(item)
    },
    handleRuleRequiredDelete(requireItem, row) {
      for (let i = 0; i < requireItem.matchExpressions.length; i++) {
        if (requireItem.matchExpressions[i] === row) {
          requireItem.matchExpressions.splice(i, 1)
        }
      }
    },
    handleRuleRequiredAdd(requireItem) {
      var item = {
        key: "",
        operator: "",
        value: "",
      }
      requireItem.matchExpressions.unshift(item)
    },
    handleSelectorRequiredAdd() {
      var item = {
        matchExpressions: [],
      }
      this.nodeSelectorTerms.unshift(item)
    },
    handleSelectorRequiredClose(requireItem) {
      for (let i = 0; i < this.nodeSelectorTerms.length; i++) {
        if (this.nodeSelectorTerms[i] === requireItem) {
          this.nodeSelectorTerms.splice(i, 1)
        }
      }
    },

    handleRulePreferDelete(preferItem, row) {
      for (let i = 0; i < preferItem.matchExpressions.length; i++) {
        if (preferItem.matchExpressions[i] === row) {
          preferItem.matchExpressions.splice(i, 1)
        }
      }
    },
    handleRulePreferAdd(preferItem) {
      var item = {
        key: "",
        operator: "",
        value: "",
      }
      preferItem.matchExpressions.unshift(item)
    },
    handleSelectorPreferAdd() {
      var item = {
        matchExpressions: [],
      }
      this.preference.unshift(item)
    },
    handleSelectorPreferClose(preferItem) {
      for (let i = 0; i < this.preference.length; i++) {
        if (this.preference[i] === preferItem) {
          this.preference.splice(i, 1)
        }
      }
    },
  },
}
</script>