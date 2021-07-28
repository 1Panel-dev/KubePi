<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'Roles'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form" ref="form" :rules="rules">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="form.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick" v-if="form.rules">
                <el-tab-pane label="Grant Resources">
                  <ko-grant-resource :rulesArray.sync="form.rules"></ko-grant-resource>
                </el-tab-pane>
                <el-tab-pane label="Labels/Annotations">
                  <ko-key-value title="Labels" :value.sync="form.metadata.labels"></ko-key-value>
                  <ko-key-value title="Annotations" :value.sync="form.metadata.annotations"></ko-key-value>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
        </div>
        <div>
          <div class="bottom-button">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
            <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
            <el-button v-if="showYaml" @click="showYaml=false">{{ $t("commons.button.back_form") }}</el-button>
            <el-button v-loading="loading" @click="onSubmit" type="primary">
              {{ $t("commons.button.submit") }}
            </el-button>
          </div>
        </div>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"
import {getRole, updateRole} from "@/api/roles"
import KoGrantResource from "@/components/ko-rbac/grant-resource"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "RoleEdit",
  components: { KoKeyValue, KoGrantResource, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      form: {
        metadata: {},
      },
      loading: false,
      showYaml: false,
      cluster: "",
      activeName: "",
      yaml: {},
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
        }
      }
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Roles" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    onSubmit () {
      if (this.showYaml) {
        this.onUpdate(this.$refs.yaml_editor.getValue())
      } else {
        this.$refs["form"].validate((valid) => {
          if (valid) {
            this.onUpdate(this.transformYaml())
          }
        })
      }
    },
    onUpdate (data) {
      this.loading = true
      updateRole(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Roles" })
      }).finally(() => {
        this.loading = false
      })
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.form))
    },
    getDetail () {
      this.loading = true
      getRole(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
