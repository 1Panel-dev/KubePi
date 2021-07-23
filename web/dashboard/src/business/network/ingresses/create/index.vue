<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'Ingresses'}" v-loading="loading">
    <div v-if="!showYaml">
      <el-row :gutter="20">
        <el-form label-position="top" :model="form" ref="form" :rules="rules">
          <el-col :span="6">
            <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="3">
            <el-form-item :label="$t('business.namespace.namespace')" required prop="metadata.namespace">
              <el-select v-model="form.metadata.namespace">
                <el-option v-for="namespace in namespaces"
                           :key="namespace.metadata.name"
                           :label="namespace.metadata.name"
                           :value="namespace.metadata.name">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-tabs v-model="activeName" tab-position="top" type="border-card"
                     @tab-click="handleClick" ref=tabs>
              <el-tab-pane label="Rules">
                <ko-ingress-rule :cluster="cluster" :namespace="form.metadata.namespace" :rulesArray.sync="form.spec.rules"></ko-ingress-rule>
              </el-tab-pane>
            </el-tabs>
          </el-col>
        </el-form>
      </el-row>
    </div>
    <div v-if="showYaml">
      <yaml-editor ref="yaml_editor" :value="yaml"></yaml-editor>
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
import {createService} from "@/api/services"
import {listNamespace} from "@/api/namespaces"
import KoIngressRule from "@/components/ko-network/ingress-rules"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "IngressCreate",
  components: { YamlEditor, KoIngressRule, LayoutContent },
  props: {},
  data () {
    return {
      cluster: "",
      namespaces: [],
      loading: false,
      form: {
        apiVersion: "networking.k8s.io/v1",
        kind: "Ingress",
        metadata: {
          namespace: "default",
        },
        spec: {}
      },
      rules: {},
      activeName: "",
      yaml: {},
      showYaml: false
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.name
    },
    onCancel () {
      this.$router.push({ name: "Ingresses" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      this.showYaml = false
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.form))
    },
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
    onCreate (data) {
      this.loading = true
      createService(this.cluster, this.form.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Services" })
      }).finally(() => {
        this.loading = false
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    listNamespace(this.cluster).then(res => {
      this.namespaces = res.items
    })
  }
}
</script>

<style scoped>

</style>
