<template>
  <div style="margin-top: 20px">
    <ko-card title="Command">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Entrypoint" prop="command">
              <ko-form-item placeholder="e.g. /bin/sh" itemType="input" v-model="form.command" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Arguments" prop="args">
              <ko-form-item placeholder="e.g. /usr/sbin/httpd -f httpd.conf" itemType="input" v-model="form.args" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="WorkingDir" prop="workingDir">
              <ko-form-item placeholder="e.g. /myapp" itemType="input" v-model="form.workingDir" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Stdin" prop="stdin">
              <ko-form-item itemType="radio" v-model="form.stdin" :radios="stdin_list" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div>
        <label>Environment Variables
          <el-tooltip class="item" effect="dark" content="ProTip: Paste lines of key=value or key: value into any key field for easy bulk entry" placement="top-start">
            <i class="el-icon-question" />
          </el-tooltip>
        </label>
      </div>
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="15%" align="left"><label>type</label></th>
          <th scope="col" width="30%" align="left"><label>prefix/variable</label></th>
          <th scope="col" width="20%" align="left"><label>source</label></th>
          <th scope="col" width="28%" align="left"><label>value</label></th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in form.envResource" v-bind:key="index">
          <td>
            <ko-form-item :withoutLabel="true" itemType="select" v-model="row.type" :selections="type_list" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="input" v-model="row.prefix_or_alias" />
          </td>
          <td>
            <ko-form-item v-if="row.type === 'Key/Value Pair' || row.type === 'Pod Field'" :withoutLabel="true" itemType="input" disabled placeholder="N/A" />
            <ko-form-item v-if="row.type === 'Resource' || row.type === 'Field'" :withoutLabel="true" itemType="input" v-model="row.source" />
            <ko-form-item v-if="row.type === 'ConfigMap key'" :withoutLabel="true" itemType="select" v-model="row.source" @change="changeConfigMap(row.source)" :selections="config_map_name_list" />
            <ko-form-item v-if="row.type === 'ConfigMap'" :withoutLabel="true" itemType="select" v-model="row.source" :selections="config_map_name_list" />
            <ko-form-item v-if="row.type === 'Secret' || row.type === 'Secret key'" :withoutLabel="true" itemType="select" v-model="row.source" :selections="secret_list" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" v-if="row.type ==='Key/Value Pair' || row.type === 'Pod Field'" itemType="input" v-model="row.value" />
            <ko-form-item :withoutLabel="true" v-if="row.type ==='Resource'" itemType="select" v-model="row.value" :selections="resource_value_list" />
            <ko-form-item :withoutLabel="true" v-if="row.type ==='ConfigMap key'" itemType="select" v-model="row.value" :selections="config_map_value_list" />
            <ko-form-item :withoutLabel="true" v-if="row.type ==='Secret key'" itemType="select" v-model="row.value" :selections="secret_value_list" />
            <ko-form-item :withoutLabel="true" v-if="row.type === 'Secret' || row.type === 'ConfigMap'" disabled itemType="input" v-model="row.key" placeholder="N/A" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleAdd">Add Variable</el-button>
          </td>
        </tr>
      </table>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"

export default {
  name: "KoCommand",
  components: { KoFormItem, KoCard },
  props: {
    commandParentObj: Object,
    namespace: String,
  },
  watch: {
    namespace: {
      handler(newName) {
        this.namespace = newName
      },
      immediate: true,
    },
  },
  data() {
    return {
      form: {
        args: "",
        command: "",
        workingDir: "",
        stdin: "",
        tty: false,
        envResource: [],
      },
      config_map_name_list: [],
      config_map_list: [],
      config_map_value_list: [],
      secret_list: [],
      resource_value_list: [
        { label: "limits.cpu", value: "limits.cpu" },
        { label: "limits.ephemeral-storage", value: "limits.ephemeral-storage" },
        { label: "limits.memory", value: "limits.memory" },
        { label: "requests.cpu", value: "requests.cpu" },
        { label: "requests.ephemeral-storage", value: "requests.ephemeral-storage" },
        { label: "requests.memory", value: "requests.memory" },
      ],
      secret_value_list: [
        { label: "ca.crt", value: "ca.crt" },
        { label: "namespace", value: "namespace" },
        { label: "token", value: "token" },
      ],
      type_list: [
        { label: "Key/Value Pair", value: "Key/Value Pair" },
        { label: "Pod Field", value: "Pod Field" },
        { label: "Resource", value: "Resource" },
        { label: "ConfigMap key", value: "ConfigMap key" },
        { label: "Secret key", value: "Secret key" },
        { label: "Secret", value: "Secret" },
        { label: "ConfigMap", value: "ConfigMap" },
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
        for (const secret of res.items) {
          this.secret_list.push({ label: secret.metadata.name, value: secret.metadata.name })
        }
      })
    },
    loadConfigMaps() {
      listConfigMaps(this.$route.params.cluster).then((res) => {
        this.config_map_name_list = []
        for (const cm of res.items) {
          this.config_map_name_list.push({ label: cm.metadata.name, value: cm.metadata.name })
          this.config_map_list.push(cm)
        }
      })
    },
    changeConfigMap(comfigmap) {
      console.log(this.namespace)
      this.config_map_value_list = []
      for (const cm of this.config_map_list) {
        if (comfigmap === cm.metadata.name && cm.metadata.namespace === this.namespace) {
          for (const item of Object.keys(cm.data)) {
            this.config_map_value_list.push({ label: item, value: item })
          }
        }
      }
    },

    handleAdd() {
      var item = {
        type: "Resource",
        prefix_or_alias: "",
        source: "",
        value: "",
      }
      this.form.envResource.push(item)
    },
    handleDelete(index) {
      this.form.envResource.splice(index, 1)
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
      if (this.form.envResource) {
        for (const en of this.form.envResource) {
          switch (en.type) {
            case "Key/Value Pair":
              envList.push({ name: en.prefix_or_alias, value: en.value })
              break
            case "Resource":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  resourceFieldRef: {
                    containerName: en.source,
                    divisor: 0,
                    resource: en.value,
                  },
                },
              })
              break
            case "ConfigMap key":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  configMapKeyRef: {
                    name: en.source,
                    key: en.value,
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
                    key: en.value,
                    optional: false,
                  },
                },
              })
              break
            case "Pod Field":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  fieldRef: {
                    apiVersion: "v1",
                    fieldPath: en.value,
                  },
                },
              })
              break
            case "Secret":
              envFromList.push({
                prefix: en.prefix_or_alias,
                secretRef: {
                  name: en.source,
                  optional: false,
                },
              })
              break
            case "ConfigMap":
              envFromList.push({
                prefix: en.prefix_or_alias,
                configMapRef: {
                  name: en.source,
                  optional: false, // 这个false是什么意思不知道
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
              this.form.envResource.push({ source: en.valueFrom.resourceFieldRef.containerName, value: en.valueFrom.resourceFieldRef.resource, type: "Resource", prefix_or_alias: en.name })
            } else if (en.valueFrom.configMapKeyRef) {
              this.form.envResource.push({ source: en.valueFrom.configMapKeyRef.name, value: en.valueFrom.configMapKeyRef.key, type: "ConfigMap key", prefix_or_alias: en.name })
            } else if (en.valueFrom.secretKeyRef) {
              this.form.envResource.push({ source: en.valueFrom.secretKeyRef.name, value: en.valueFrom.secretKeyRef.key, type: "Secret key", prefix_or_alias: en.name })
            } else if (en.valueFrom.fieldRef) {
              this.form.envResource.push({ value: en.valueFrom.fieldRef.fieldPath, type: "Pod Field", prefix_or_alias: en.name })
            }
          } else if (en.value) {
            this.form.envResource.push({ value: en.value, type: "Key/Value Pair", prefix_or_alias: en.name })
          }
        }
      }
      if (this.commandParentObj.envFrom) {
        for (const en of this.commandParentObj.envFrom) {
          if (en.configMapRef) {
            this.form.envResource.push({ source: en.configMapRef.name, type: "ConfigMap", prefix_or_alias: en.prefix })
          } else if (en.secretRef) {
            this.form.envResource.push({ source: en.secretRef.name, type: "Secret", prefix_or_alias: en.prefix })
          }
        }
      }
    }
  },
}
</script>