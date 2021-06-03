<template>
  <div style="margin-top: 20px">
    <el-row :gutter="40">
      <el-col :span="12">
        <ko-form-item labelName="Node Scheduling" clearable itemType="radio" v-model="form.scheduling_type" :radios="scheduling_type_list" />
      </el-col>
    </el-row>

    <el-row :gutter="40" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Node Name" clearable itemType="select" v-model="form.nodename" :selections="node_list" />
      </el-col>
    </el-row>

    <div style="margin-top: 20px">
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
            <ko-form-item placeholder="e.g. foo" clearable itemType="input" v-model="row.name" />
          </template>
        </el-table-column>
        <el-table-column min-width="80" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item placeholder="e.g. bar" clearable itemType="input" v-model="row.value" />
          </template>
        </el-table-column>
        <el-table-column min-width="20">
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
        { label: "Run pods on specific node(s)", value: "No" },
        { label: "Run pods on node(s) matching these scheduling rules", value: "Yes" },
      ],
      node_list: [],
      form: {
        scheduling_type: "",
        nodename: "",
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
  },
}
</script>