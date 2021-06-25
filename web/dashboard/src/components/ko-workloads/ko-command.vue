<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Entrypoint" placeholder="e.g. /bin/sh" itemType="input" v-model="form.command" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Arguments" placeholder="e.g. /usr/sbin/httpd -f httpd.conf" itemType="input" v-model="form.args" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="WorkingDir" placeholder="e.g. /myapp" itemType="input" v-model="form.workingDir" />
      </el-col>
      <el-col :span="6">
        <ko-form-item labelName="Stdin" itemType="select" v-model="form.stdin" :selections="stdin_list" />
      </el-col>
      <el-col :span="6">
        <el-checkbox style="margin-top: 20px" :disabled="form.stdin === 'No'" v-model="form.tty">TTY</el-checkbox>
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
          <ko-form-item :withoutLabel="true" placeholder="e.g. foo" itemType="input" v-model="row.name" />
        </template>
      </el-table-column>
      <el-table-column min-width="80" label="Value">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" placeholder="e.g. bar" itemType="input" v-model="row.value" />
        </template>
      </el-table-column>
      <el-table-column width="60px">
        <template v-slot:default="{row}">
          <el-button type="text" style="font-size: 10px" @click="handleDelete(row)">{{ $t("commons.button.delete") }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-button style="margin-top: 20px" @click="handleResourceAdd">Add form Resource</el-button>
    <el-table v-if="form.envFromResource.length !== 0" :data="form.envFromResource">
      <el-table-column min-width="25" label="Type">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="select" v-model="row.type" :selections="type_list" />
        </template>
      </el-table-column>
      <el-table-column min-width="40" label="Source">
        <template v-slot:default="{row}">
          <ko-form-item v-if="row.type === 'Resource' || row.type === 'Field'" :withoutLabel="true" itemType="input" v-model="row.source" />
          <ko-form-item v-if="row.type === 'ConfigMap' || row.type === 'ConfigMap key'" :withoutLabel="true" itemType="select" v-model="row.source" :selections="config_map_list" />
          <ko-form-item v-if="row.type === 'Secret' || row.type === 'Secret key'" :withoutLabel="true" itemType="select" v-model="row.source" :selections="secret_list" />
        </template>
      </el-table-column>
      <el-table-column min-width="40" label="Key">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" v-if="row.type ==='Resource'" itemType="select" v-model="row.key" :selections="resource_key_list" />
          <ko-form-item :withoutLabel="true" v-if="row.type ==='ConfigMap key'" itemType="select" v-model="row.key" :selections="cm_key_list" />
          <ko-form-item :withoutLabel="true" v-if="row.type ==='Secret key'" itemType="select" v-model="row.key" :selections="secret_key_list" />
          <ko-form-item :withoutLabel="true" v-if="row.type === 'Secret' || row.type === 'ConfigMap' || row.type === 'Field'" disabled itemType="input" v-model="row.key" placeholder="N/A" />
        </template>
      </el-table-column>
      <el-table-column min-width="10">
        <span style="font-size: 20px; color: white;">AS</span>
      </el-table-column>
      <el-table-column min-width="40" label="Prefix or Alias">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="input" v-model="row.prefix_or_alias" />
        </template>
      </el-table-column>
      <el-table-column width="60px">
        <template v-slot:default="{row}">
          <el-button type="text" style="font-size: 10px" @click="handleResourceDelete(row)">{{ $t("commons.button.delete") }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"

export default {
  name: "KoCommand",
  components: { KoFormItem },
  props: {
    commandParentObj: Object,
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
      config_map_list: [],
      secret_list: [],
      resource_key_list: [
        { label: "limits.cpu", value: "limits.cpu" },
        { label: "limits.ephemeral-storage", value: "limits.ephemeral-storage" },
        { label: "limits.memory", value: "limits.memory" },
        { label: "requests.cpu", value: "requests.cpu" },
        { label: "requests.ephemeral-storage", value: "requests.ephemeral-storage" },
        { label: "requests.memory", value: "requests.memory" },
      ],
      cm_key_list: [],
      secret_key_list: [],
      type_list: [
        { label: "Resource", value: "Resource" },
        { label: "ConfigMap", value: "ConfigMap" },
        { label: "Secret key", value: "Secret key" },
        { label: "Field", value: "Field" },
        { label: "Secret", value: "Secret" },
        { label: "ConfigMap key", value: "ConfigMap key" },
      ],
      stdin_list: [
        { label: "No", value: "No" },
        { label: "Ones", value: "Ones" },
        { label: "Yes", value: "Yes" },
      ],
    }
  },
  methods: {
    loadSecrets() {
      listSecrets(this.$route.params.cluster).then((res) => {
        this.secret_list = []
        for(const secret of res.items) {
          this.secret_list.push({label: secret.metadata.name, value: secret.metadata.name})
        }
      })
    },
    loadConfigMaps() {
      listConfigMaps(this.$route.params.cluster).then((res) => {
        this.config_map_list = []
        for(const cm of res.items) {
          this.config_map_list.push({label: cm.metadata.name, value: cm.metadata.name})
        }
      })
    },

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
    
    transformation(parentFrom) {
      if (this.form.command) {
        parentFrom.command = this.form.command.split(",")
      }
      if (this.form.args) {
        parentFrom.args = this.form.args.split(",")
      }
      if (this.form.workingDir) {
        parentFrom.workingDir = this.form.workingDir
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
      }
      if (this.form.tty) {
        parentFrom.tty = this.form.tty
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
            case "Secret key":
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
      if (envList.length !== 0) {
        parentFrom.env = envList
      }
      if (envFromList.length !== 0) {
        parentFrom.envFrom = envFromList
      }
    },
  },
  mounted() {
    this.loadSecrets()
    this.loadConfigMaps()
    if (this.commandParentObj) {
      if (this.commandParentObj.command) {
        this.form.command = this.commandParentObj.command.join(",")
      }
      if (this.commandParentObj.args) {
        this.form.args = this.commandParentObj.args.join(",")
      }
      if (this.commandParentObj.workingDir) {
        this.form.workingDir = this.commandParentObj.workingDir
      }
      if (this.commandParentObj.stdinOnce != undefined) {
        this.form.stdinOnce = true
        this.form.stdin = "Ones"
      }
      if (this.commandParentObj.stdin != undefined) {
        this.form.stdin = this.commandParentObj.stdin ? "Yes" : "No"
      }
      if (this.commandParentObj.tty) {
        this.form.tty = this.commandParentObj.tty
      }
      if (this.commandParentObj.env) {
        for (const en of this.commandParentObj.env) {
          if (en.valueFrom) {
            if (en.valueFrom.resourceFieldRef) {
              this.form.envFromResource.push({ source: en.valueFrom.resourceFieldRef.containerName, key: en.valueFrom.resourceFieldRef.resource, type: "Resource", prefix_or_alias: en.name })
            } else if (en.valueFrom.configMapKeyRef) {
              this.form.envFromResource.push({ source: en.valueFrom.configMapKeyRef.name, key: "ConfigMap", type: "ConfigMap", prefix_or_alias: en.name })
            } else if (en.valueFrom.secretKeyRef) {
              this.form.envFromResource.push({ source: en.valueFrom.secretKeyRef.name, key: "Secret key", type: "Secret key", prefix_or_alias: en.name })
            } else if (en.valueFrom.fieldRef) {
              this.form.envFromResource.push({ source: en.valueFrom.fieldRef.fieldPath, type: "Field", prefix_or_alias: en.name })
            }
          } else {
            this.form.env.push(en)
          }
        }
      }
      if (this.commandParentObj.envFrom) {
        for (const en of this.commandParentObj.envFrom) {
          if (en.valueFrom) {
            if (en.valueFrom.secretRef) {
              this.form.envFromResource.push({ source: en.valueFrom.secretRef.name, key: "Secret", type: "Secret", prefix_or_alias: en.name })
            }
          }
        }
      }
    }
  },
}
</script>