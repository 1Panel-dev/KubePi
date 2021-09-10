<template>
  <div>
    <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="15%" align="left"><label>{{$t('business.workload.type')}}</label></th>
          <th scope="col" width="30%" align="left"><label>{{$t('business.workload.prefix_variable')}}</label></th>
          <th scope="col" width="20%" align="left"><label>{{$t('business.workload.source')}}</label></th>
          <th scope="col" width="25%" align="left"><label>{{$t('business.workload.value')}}</label></th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in form.envResource" v-bind:key="index">
          <td>
            <ko-form-item itemType="select2" :noClear="true" v-model="row.type" :selections="type_list" />
          </td>
          <td>
            <ko-form-item itemType="input" v-model="row.prefix_or_alias" />
          </td>
          <td>
            <ko-form-item v-if="row.type === 'Key/Value Pair' || row.type === 'Pod Field'" itemType="input" disabled placeholder="N/A" />
            <ko-form-item v-if="row.type === 'Resource' || row.type === 'Field'" itemType="input" v-model="row.source" />
            <ko-form-item v-if="row.type === 'ConfigMap key'" itemType="select2" v-model="row.source" @change="changeConfigMap(row)" :selections="config_map_name_list" />
            <ko-form-item v-if="row.type === 'ConfigMap'" itemType="select2" v-model="row.source" :selections="config_map_name_list" />
            <ko-form-item v-if="row.type === 'Secret'" itemType="select2" v-model="row.source" :selections="secret_name_list" />
            <ko-form-item v-if="row.type === 'Secret key'" itemType="select2" @change="changeSecret(row)" v-model="row.source" :selections="secret_name_list" />
          </td>
          <td>
            <ko-form-item v-if="row.type ==='Key/Value Pair' || row.type === 'Pod Field'" itemType="textarea" v-model="row.value" />
            <ko-form-item v-if="row.type ==='Resource'" itemType="select2" v-model="row.value" :selections="resource_value_list" />
            <ko-form-item v-if="row.type ==='Secret key' || row.type ==='ConfigMap key'" itemType="select2" v-model="row.value" :selections="row.value_list" />
            <ko-form-item v-if="row.type === 'Secret' || row.type === 'ConfigMap'" disabled itemType="input" v-model="row.key" placeholder="N/A" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleAdd">{{$t("commons.button.add")}}{{$t("business.workload.variable")}}</el-button>
          </td>
        </tr>
      </table>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoCommand",
  components: { KoFormItem },
  props: {
    envParentObj: Object,
    currentNamespace: String,
    configMapList: Array,
    secretList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    currentNamespace: {
      handler(newName) {
        this.namespace = newName
      },
      immediate: true,
    },
    configMapList: {
      handler(newName) {
        this.config_map_name_list = []
        this.config_map_list = []
        for (const cm of newName) {
          this.config_map_name_list.push(cm.metadata.name)
          this.config_map_list.push(cm)
        }
      },
      immediate: true,
      deep: true,
    },
    secretList: {
      handler(newName) {
        this.secret_list = newName
        this.secret_name_list = []
        for (const s of newName) {
          this.secret_name_list.push(s.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      form: {
        envResource: [],
      },
      reFresh: false,
      namespace: "",
      config_map_list: [],
      config_map_name_list: [],
      secret_list: [],
      secret_name_list: [],
      resource_value_list: ["limits.cpu", "limits.memory", "requests.cpu", "requests.memory"],
      type_list: ["Key/Value Pair", "Pod Field", "Resource", "ConfigMap key", "Secret key", "Secret", "ConfigMap"],
    }
  },
  methods: {
    changeConfigMap(row) {
      row.value_list = []
      for (const cm of this.config_map_list) {
        if (row.source === cm.metadata.name && cm.metadata.namespace === this.namespace) {
          for (const item of Object.keys(cm.data)) {
            row.value_list.push(item)
          }
        }
      }
    },
    changeSecret(row) {
      row.value_list = []
      for (const se of this.secret_list) {
        if (row.source === se.metadata.name && se.metadata.namespace === this.namespace) {
          for (const item of Object.keys(se.data)) {
            row.value_list.push(item)
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
        value_list: [],
      }
      this.form.envResource.push(item)
    },
    handleDelete(index) {
      this.form.envResource.splice(index, 1)
    },

    transformation(parentFrom) {
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
                },
              })
              break
            case "ConfigMap":
              envFromList.push({
                prefix: en.prefix_or_alias,
                configMapRef: {
                  name: en.source,
                },
              })
              break
          }
        }
      }
      parentFrom.env = envList.length !== 0 ? envList : undefined
      parentFrom.envFrom = envFromList.length !== 0 ? envFromList : undefined
    },
  },
  mounted() {
    if (this.envParentObj) {
      if (this.envParentObj.env) {
        for (const en of this.envParentObj.env) {
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
      if (this.envParentObj.envFrom) {
        for (const en of this.envParentObj.envFrom) {
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