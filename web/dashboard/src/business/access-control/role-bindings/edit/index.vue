<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'RoleBindings'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="form.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="form.metadata.namespace" disabled>
                  <el-option :key="form.metadata.namespace"
                             :label="form.metadata.namespace"
                             :value="form.metadata.namespace"
                  ></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       v-if="Object.keys(form.metadata).length!==0">
                <el-tab-pane :label="$t('business.access_control.authorized')">
                  <ko-role-object :cluster-name="cluster" :namespace="form.metadata.namespace"
                                  :role-ref-obj.sync="form.roleRef"
                                  :subject-array.sync="form.subjects"></ko-role-object>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.workload.labels_annotations')">
                  <ko-key-value :title="$t('business.workload.label')"
                                :value.sync="form.metadata.labels"></ko-key-value>
                  <ko-key-value :title="$t('business.workload.annotations')"
                                :value.sync="form.metadata.annotations"></ko-key-value>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
        </div>
      </el-row>
    </div>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
      <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
      <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
      <el-button v-loading="loading" @click="onSubmit" type="primary">
        {{ $t("commons.button.submit") }}
      </el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getRoleBinding, updateRoleBinding} from "@/api/rolebings"
import KoRoleObject from "@/components/ko-rbac/role-object"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "RoleBindingEdit",
  components: { YamlEditor, LayoutContent, KoRoleObject, KoKeyValue },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      loading: false,
      cluster: "",
      form: {
        metadata: {}
      },
      showYaml: false,
      yaml: {},
      activeName: ""
    }
  },
  methods: {
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updateRoleBinding(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "RoleBindings" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "RoleBindings" })
    },
    getDetail () {
      this.loading = true
      getRoleBinding(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        this.yaml = this.transformYaml()
        this.loading = false
      })
    },
    onEditYaml () {
      this.yaml = this.transformYaml()
      this.showYaml = true
    },
    backToForm () {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.form))
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
