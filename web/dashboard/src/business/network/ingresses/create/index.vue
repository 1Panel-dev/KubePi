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
              <ko-select :namespace.sync="form.metadata.namespace"></ko-select>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-tabs v-model="activeName" tab-position="top" type="border-card"
                     @tab-click="handleClick" ref=tabs>
              <el-tab-pane :label="$t('business.network.rule')">
                <ko-ingress-rule :cluster="cluster" :namespace="form.metadata.namespace"
                                 :rulesArray.sync="form.spec.rules"></ko-ingress-rule>
              </el-tab-pane>
              <el-tab-pane :label="$t('business.network.default_backend')">
                <ko-ingress-default-backend :cluster="cluster" :namespace="form.metadata.namespace"
                                            :defaultBackendObj.sync="form.spec.defaultBackend"></ko-ingress-default-backend>
              </el-tab-pane>
              <el-tab-pane :label="$t('business.configuration.certificate')">
                <ko-ingress-tls :cluster="cluster" :namespace="form.metadata.namespace"
                                :tlsArray.sync="form.spec.tls"></ko-ingress-tls>
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
import KoIngressRule from "@/components/ko-network/ingress-rules"
import YamlEditor from "@/components/yaml-editor"
import KoIngressDefaultBackend from "@/components/ko-network/ingress-defaultbackend"
import KoIngressTls from "@/components/ko-network/ingress-tls"
import {createIngress} from "@/api/ingress"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import KoSelect from "@/components/ko-select"

export default {
  name: "IngressCreate",
  components: {
    KoKeyValue,
    KoIngressTls,
    KoIngressDefaultBackend,
    YamlEditor,
    KoIngressRule,
    LayoutContent,
    KoSelect,
  },
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
          namespace: "",
        },
        spec: {}
      },
      rules: {},
      activeName: "",
      yaml: undefined,
      showYaml: false
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "Ingresses" })
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
      createIngress(this.cluster, data.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Ingresses" })
      }).finally(() => {
        this.loading = false
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
  }
}
</script>

<style scoped>

</style>
