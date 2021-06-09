<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Entrypoint" placeholder="e.g. /bin/sh" clearable itemType="input" v-model="form.command" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Arguments" placeholder="e.g. /usr/sbin/httpd -f httpd.conf" clearable itemType="input" v-model="form.args" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="WorkingDir" placeholder="e.g. /myapp" clearable itemType="input" v-model="form.workingDir" />
      </el-col>
      <el-col :span="6">
        <ko-form-item labelName="Stdin" clearable itemType="select" v-model="form.stdin" :selections="stdin_list" />
      </el-col>
      <el-col :span="6">
        <el-checkbox :disabled="form.stdin === 'No'" v-model="form.tty">TTY</el-checkbox>
      </el-col>
    </el-row>
    <div style="margin-top: 30px">
      <label>Environment Variables
        <el-tooltip class="item" effect="dark" content="ProTip: Paste lines of key=value or key: value into any key field for easy bulk entry" placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>
      </label>
    </div>

    <el-button style="margin-top: 20px" @click="handleAdd">Add</el-button>
    <el-table v-if="form.env.length !== 0" :data="form.env">
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

    <el-button style="margin-top: 20px" @click="handleResourceAdd">Add form Resource</el-button>
    <el-table v-if="form.envFromResource.length !== 0" :data="form.envFromResource">
      <el-table-column min-width="40" label="Type">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" @change="changeType(row, 'type')" itemType="select" v-model="row.type" :selections="type_list" />
        </template>
      </el-table-column>
      <el-table-column min-width="40" label="Source">
        <template v-slot:default="{row}">
          <ko-form-item v-if="row.type === 'Resource' || row.type === 'Field'" :withoutLabel="true" clearable itemType="input" v-model="row.source" />
          <ko-form-item v-if="row.type === 'ConfigMap'" :withoutLabel="true" clearable itemType="select" v-model="row.source" :selections="config_map_list" />
          <ko-form-item v-if="row.type === 'Secret key'" :withoutLabel="true" clearable itemType="select" v-model="row.source" :selections="secret_key_list" />
          <ko-form-item v-if="row.type === 'Secret'" :withoutLabel="true" clearable itemType="select" v-model="row.source" :selections="secret_list" />
        </template>
      </el-table-column>
      <el-table-column min-width="40" label="Key">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" placeholder="e.g. requests.cpu" v-if="row.type ==='Resource'" clearable itemType="input" v-model="row.key" />
          <ko-form-item :withoutLabel="true" v-if="row.type === 'ConfigMap' || row.type === 'Secret key' || row.type === 'Secret'" @change="changeType(row, 'key')" clearable itemType="select" v-model="row.key" :selections="type_list" />
          <ko-form-item :withoutLabel="true" v-if="row.type ==='Field'" disabled clearable itemType="input" value="n/a" />
        </template>
      </el-table-column>
      <el-table-column min-width="10">
        <span style="font-size: 20px; color: white;">AS</span>
      </el-table-column>
      <el-table-column min-width="40" label="Prefix or Alias">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.prefix_or_alias" />
        </template>
      </el-table-column>
      <el-table-column width="120px">
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
  props: {
    commandObj: Object,
  },
  data() {
    return {
      form: {
        args: "",
        command: "",
        workingDir: "",
        stdin: "",
        tty: false,
        env: [],
        envFromResource: [],
      },
      config_map_list: [{ label: "kube-root-ca.crt", value: "kube-root-ca.crt" }],
      secret_key_list: [{ label: "default-token-fkpdx", value: "default-token-fkpdx" }],
      secret_list: [{ label: "default-token-fkpdx", value: "default-token-fkpdx" }],
      type_list: [
        { label: "Resource", value: "Resource" },
        { label: "ConfigMap", value: "ConfigMap" },
        { label: "Secret key", value: "Secret key" },
        { label: "Field", value: "Field" },
        { label: "Secret", value: "Secret" },
      ],
      stdin_list: [
        { label: "No", value: "No" },
        { label: "Ones", value: "Ones" },
        { label: "Yes", value: "Yes" },
      ],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.form.env.length; i++) {
        if (this.form.env[i] === row) {
          this.form.env.splice(i, 1)
        }
      }
    },
    handleResourceDelete(row) {
      for (let i = 0; i < this.form.envFromResource.length; i++) {
        if (this.form.envFromResource[i] === row) {
          this.form.envFromResource.splice(i, 1)
        }
      }
    },
    handleAdd() {
      var item = {
        name: "",
        value: "",
      }
      this.form.env.unshift(item)
    },
    handleResourceAdd() {
      var item = {
        type: "Resource",
        source: "",
        key: "",
        prefix_or_alias: "",
      }
      this.form.envFromResource.unshift(item)
    },
    changeType(row, operator) {
      if (operator === "type") {
        row.key = row.type
      } else {
        row.type = row.key
      }
    },
    transformation(parentFrom) {
      if (this.form.args) {
        parentFrom.args = this.form.args.split(",")
      } else {
        delete parentFrom.args
      }
      if (this.form.workingDir) {
        parentFrom.workingDir = this.form.workingDir
      } else {
        delete parentFrom.workingDir
      }
      if (this.form.stdin) {
        switch (this.form.stdin) {
          case "No":
            parentFrom.stdin = false
            break
          case "Yes":
            parentFrom.stdin = true
            break
          case "Ones":
            parentFrom.stdin = true
            parentFrom.stdinOnce = true
            break
        }
      } else {
        delete parentFrom.stdin
      }
      if (this.form.tty) {
        parentFrom.tty = this.form.tty
      } else {
        delete parentFrom.tty
      }
      let envList = []
      let envFromList = []
      if (this.form.env || this.form.envFromResource) {
        for (const en of this.form.env) {
          envList.push(en)
        }
        for (const en of this.form.envFromResource) {
          switch (en.type) {
            case "Resource":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  resourceFieldRef: {
                    containerName: en.source,
                    divisor: 0,
                    resource: en.key,
                  },
                },
              })
              break
            case "ConfigMap":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  configMapKeyRef: {
                    name: en.source,
                    optional: false, // 这个false是什么意思不知道
                  },
                },
              })
              break
            case "Secret Key":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  secretKeyRef: {
                    name: en.source,
                    optional: false,
                  },
                },
              })
              break
            case "Field":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  fieldRef: {
                    apiVersion: "v1",
                    fieldPath: en.source,
                  },
                },
              })
              break
            case "Secret":
              envFromList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  secretRef: {
                    name: en.source,
                    optional: false,
                  },
                },
              })
              break
          }
        }
      }
      parentFrom.env = envList
      parentFrom.envFrom = envFromList
    }
  },
  mounted() {
    if (this.commandObj) {
      // 给form赋值
    }
  },
}
</script>