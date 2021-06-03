<template>
  <div>
    <el-row :gutter="40" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Entrypoint" clearable itemType="input" v-model="form.entrypoint" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Arguments" clearable itemType="input" v-model="form.arguments" />
      </el-col>
    </el-row>
    <el-row :gutter="40" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="WorkingDir" clearable itemType="input" v-model="form.working_dir" />
      </el-col>
      <el-col :span="6">
        <ko-form-item labelName="Stdin" clearable itemType="input" v-model="form.stdin" />
      </el-col>
      <el-col :span="6">
        <el-checkbox v-model="form.tty_check">TTY</el-checkbox>
      </el-col>
    </el-row>
    <div style="margin-top: 30px"><label>Environment Variables</label></div>

    <el-button style="margin-top: 20px" @click="handleAdd">Add</el-button>
    <el-table style="color: black; background-color: black" :show-header="false" :data="data">
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item labelName="Key" clearable itemType="input" v-model="row.key" />
        </template>
      </el-table-column>
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item labelName="Value" clearable itemType="input" v-model="row.value" />
        </template>
      </el-table-column>
      <el-table-column min-width="20">
        <template v-slot:default="{row}">
          <el-button type="text" style="font-size: 20px" @click="handleDelete(row)">REMOVE</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-button style="margin-top: 20px" @click="handleResourceAdd">Add form Resource</el-button>
    <el-table style="color: black; background-color: black" :show-header="false" :data="resource_data">
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item labelName="Type" clearable itemType="select" v-model="row.type" :selections="type_list" />
        </template>
      </el-table-column>
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item labelName="Source" clearable itemType="input" v-model="row.source" />
        </template>
      </el-table-column>
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item v-if="row.type ==='Resource'" labelName="Key" clearable itemType="input" v-model="row.key" />
          <ko-form-item v-if="row.type === 'ConfigMap' || row.type === 'Secert key' || row.type === 'Secert'" labelName="Key" clearable itemType="input" v-model="row.key" />
          <ko-form-item v-if="row.type ==='Field'" labelName="Key" disabled clearable itemType="input" v-model="row.key" />
        </template>
      </el-table-column>
      <el-table-column min-width="20">
          <span style="font-size: 20px; color: white;">AS</span>
      </el-table-column>
      <el-table-column min-width="80">
        <template v-slot:default="{row}">
          <ko-form-item labelName="Prefix or Alias" clearable itemType="input" v-model="row.prefix_or_alias" />
        </template>
      </el-table-column>
      <el-table-column min-width="20">
        <template v-slot:default="{row}">
          <el-button type="text" style="font-size: 20px" @click="handleResourceDelete(row)">REMOVE</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoCommand",
  components: { KoFormItem },
  data() {
    return {
      data: [],
      resource_data: [],
      form: {
        entrypoint: "",
        arguments: "",
        working_dir: "",
        stdin: "",
        tty_check: false,
      },
      type_list: [
        { label: "Resource", value: "Resource" },
        { label: "ConfigMap", value: "ConfigMap" },
        { label: "Secert key", value: "Secert key" },
        { label: "Field", value: "Field" },
        { label: "Secert", value: "Secert" },
      ],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.data.length; i++) {
        if (this.data[i] === row) {
          this.data.splice(i, 1)
        }
      }
    },
    handleResourceDelete(row) {
      for (let i = 0; i < this.resource_data.length; i++) {
        if (this.resource_data[i] === row) {
          this.resource_data.splice(i, 1)
        }
      }
    },
    handleAdd() {
      var item = {
        key: "",
        value: "",
      }
      this.data.unshift(item)
    },
    handleResourceAdd() {
      var item = {
        type: "Resource",
        source: "",
        key: "",
        prefix_or_alias: "",
      }
      this.resource_data.unshift(item)
    },
  },
}
</script>