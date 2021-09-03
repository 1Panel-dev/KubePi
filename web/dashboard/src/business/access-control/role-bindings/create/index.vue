<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'RoleBindings'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form" ref="form" :rules="rules">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required prop="metadata.namespace">
                <ko-select :namespace.sync="form.metadata.namespace"></ko-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card">
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
      </el-row>
      <div v-if="showYaml">
        <yaml-editor :value="form" ref="yaml_editor"></yaml-editor>
      </div>
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
import {createRoleBinding} from "@/api/rolebings"
import Rule from "@/utils/rules"
import KoRoleObject from "@/components/ko-rbac/role-object"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import KoSelect from "@/components/ko-select"


export default {
  name: "RoleBindingCreate",
  components: { KoRoleObject, YamlEditor, LayoutContent,KoKeyValue,KoSelect},
  data () {
    return {
      loading: false,
      cluster: "",
      form: {
        apiVersion: "rbac.authorization.k8s.io/v1",
        kind: "RoleBinding",
        metadata: {
          name: "",
          namespace: ""
        },
      },
      showYaml: false,
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule]
        }
      },
      activeName: "",
      yaml: undefined
    }
  },
  methods: {
    onSubmit () {
      if (this.showYaml) {
        this.onCreate(this.$refs.yaml_editor.getValue())
      } else {
        this.$refs["form"].validate((valid) => {
          if (valid) {
            this.onCreate(this.transformYaml())
          }
        })
      }
    },
    onCancel () {
      this.$router.push({ name: "RoleBindings" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
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
      if (this.form?.subjects && this.form?.subjects?.length === 0) {
        delete this.form.subjects
      }
      return JSON.parse(JSON.stringify(this.form))
    },
    onCreate (data) {
      createRoleBinding(this.cluster, data.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "RoleBindings" })
      }).finally(() => {
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
  }
}
</script>

<style scoped>

</style>
