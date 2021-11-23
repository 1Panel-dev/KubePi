<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <div class="grid-content bg-purple-light" v-if="!showYaml">
      <el-form label-position="top" :model="form" :rules="rules" ref="form">
        <el-row :gutter="20" style="margin-left: 5px">
          <el-col :span="4">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <el-input disabled clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-form-item disabled itemType="select2" :selections="[]" v-model="form.metadata.namespace" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row>
          <table style="width: 98%" class="tab-table">
            <tr>
              <th scope="col" width="30%" align="left"><label>{{$t('business.workload.resource')}}</label></th>
              <th scope="col" width="30%" align="left"><label>{{$t('business.workload.limit')}}</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row, index) in hards" v-bind:key="index">
              <td>
                <ko-form-item itemType="select2" @change="changeType" v-model="row.key" :selections="resource_list" />
              </td>
              <td>
                <ko-form-item itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
              </td>
            </tr>
          </table>
        </el-row>

        <el-row style="margin-top:20px">
          <table style="width: 98%" class="tab-table">
            <tr>
              <th scope="col" width="30%" align="left"><label>ScopeName</label></th>
              <th scope="col" width="30%" align="left"><label>Operator</label></th>
              <th scope="col" width="30%" align="left"><label>Values</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row, index) in selectors" v-bind:key="index">
              <td>
                <ko-form-item itemType="select" @change="changeScopeName(row)" v-model="row.scopeName" :selections="scope_list" />
              </td>
              <td>
                <ko-form-item itemType="select2" :disabled="isOperatorOnlyExist(row.scopeName)" v-model="row.operator" :selections="operator_list" />
              </td>
              <td>
                <ko-form-item :disabled="row.operator === 'Exists' || row.operator === 'DoesNotExist'" :placeholder="$t('business.configuration.split_help')" itemType="input" v-model="row.values" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleScopeDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleScopeAdd">{{ $t("commons.button.add") }}</el-button>
              </td>
            </tr>
          </table>
        </el-row>

      </el-form>
    </div>

    <div class="grid-content bg-purple-light" v-if="showYaml">
      <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="loading" @click="onSubmit" type="primary">
          {{ $t("commons.button.submit") }}
        </el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"
import KoFormItem from "@/components/ko-form-item/index"
import { parseArryToObj, parseObjToArry } from "@/utils/objArryParse"
import { getResourceQuota, updateResourceQuota } from "@/api/resourcequota"

export default {
  name: "ResourceQuotaEdit",
  components: { YamlEditor, LayoutContent, KoFormItem },
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "ResourceQuota",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          hard: {},
          scopeSelector: {
            matchExpressions: [],
          },
        },
      },
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        },
      },
      hards: [],
      all_resource_list: ["limits.cpu", "limits.memory", "requests.cpu", "requests.memory", "configmaps", "pods", "replicationcontrollers", "resourcequotas", "services", "services.loadbalancers", "services.nodeports", "secrets", "requests.storage", "persistentvolumeclaims"],
      resource_list: [],

      selectors: [],
      scope_list: [
        { label: "Terminating (.spec.activeDeadlineSeconds >= 0)", value: "Terminating" },
        { label: "NotTerminating (.spec.activeDeadlineSeconds is nil)", value: "NotTerminating" },
        { label: "BestEffort (have best effort quality of service)", value: "BestEffort" },
        { label: "NotBestEffort (not have best effort quality of service)", value: "NotBestEffort" },
        { label: "PriorityClass (references the specified priority class)", value: "PriorityClass" },
      ],
      operator_list: ["In", "NotIn", "Exists", "DoesNotExist"],

      showYaml: false,
      yaml: undefined,
      cluster: "",
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getResourceQuota(this.cluster, this.namespace, this.name)
        .then((res) => {
          this.loading = true
          this.form = res
          this.hards = parseObjToArry(this.form.spec.hard)
          if (this.form.spec?.scopeSelector?.matchExpressions) {
            for (const item of this.form.spec.scopeSelector.matchExpressions) {
              this.selectors.push({
                scopeName: item.scopeName,
                operator: item.operator,
                values: item.values ? item.values.join(",") : "",
              })
            }
          }
        })
        .finally(() => {
          this.loading = false
        })
    },
    handleAdd() {
      var item = {
        key: "",
        value: "",
      }
      this.hards.push(item)
    },
    handleDelete(index) {
      this.hards.splice(index, 1)
    },
    handleScopeAdd() {
      var item = {
        scopeName: "",
        operator: "",
        value: "",
      }
      this.selectors.push(item)
    },
    handleScopeDelete(index) {
      this.selectors.splice(index, 1)
    },
    changeType() {
      let newTypeList = []
      for (const t of this.all_resource_list) {
        let isExist = false
        for (const item of this.hards) {
          if (item.key === t) {
            isExist = true
            break
          }
        }
        if (!isExist) {
          newTypeList.push(t)
        }
      }
      this.resource_list = newTypeList
    },
    changeScopeName(row) {
      if (this.isOperatorOnlyExist(row.scopeName)) {
        row.operator = "Exists"
      }
    },
    isOperatorOnlyExist(scopeName) {
      return scopeName === "Terminating" || scopeName === "NotTerminating" || scopeName === "BestEffort" || scopeName === "NotBestEffort"
    },
    beforeSubmit() {
      this.form.spec.hard = parseArryToObj(this.hards)
      if (this.selectors.length === 0) {
        delete this.form.spec.scopeSelector
        return
      }
      this.form.spec.scopeSelector = { matchExpressions: [] }
      for (const scope of this.selectors) {
        this.form.spec.scopeSelector.matchExpressions.push({
          scopeName: scope.scopeName,
          operator: scope.operator,
          values: scope.values ? scope.values.split(",") : undefined,
        })
      }
    },
    onSubmit() {
      this.loading = true
      let data
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        this.beforeSubmit()
        data = this.form
      }
      updateResourceQuota(this.cluster, data.metadata.namespace, data.metadata.name, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "ResourceQuotas" })
        })
        .finally(() => {
          this.loading = false
        })
    },
    onCancel() {
      this.$router.push({ name: "ResourceQuotas" })
    },
    transformYaml() {
      return JSON.parse(JSON.stringify(this.form))
    },
    onEditYaml() {
      this.beforeSubmit()
      this.yaml = this.form
      this.showYaml = true
    },
    backToForm() {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    handleClick(tab) {
      this.activeName = tab.index
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
    this.changeType()
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
