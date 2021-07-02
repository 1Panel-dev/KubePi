<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'Secrets'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="form.metadata.namespace">
                  <el-option v-for="namespace in namespaces"
                             :key="namespace.metadata.name"
                             :label="namespace.metadata.name"
                             :value="namespace.metadata.name">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.configuration.type')" required>
                <el-select v-model="form.type" @change="changeType">
                  <el-option label="Opaque" value="Opaque"></el-option>
                  <el-option label="Service Account Token	" value="kubernetes.io/service-account-token"></el-option>
                  <el-option label="Docker Registry" value="kubernetes.io/dockerconfigjson"></el-option>
                  <el-option :label="$t('business.configuration.basic_auth')"
                             value="kubernetes.io/basic-auth"></el-option>
                  <el-option :label="$t('business.configuration.ssh_auth')" value="kubernetes.io/ssh-auth"></el-option>
                  <el-option :label="$t('business.configuration.tls_auth')" value="kubernetes.io/tls"></el-option>
                  <el-option :label="$t('business.configuration.token_auth')"
                             value="bootstrap.kubernetes.io/token"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane label="Data" v-if="form.type==='Opaque'">
                  <ko-secret-data :dataObj.sync="form.data"></ko-secret-data>
                </el-tab-pane>
                <el-tab-pane label="Data" v-if="form.type==='kubernetes.io/dockerconfigjson'">
                  <ko-secret-docker-data :dataObj.sync="form.data"></ko-secret-docker-data>
                </el-tab-pane>
                <el-tab-pane label="Data" v-if="form.type==='kubernetes.io/ssh-auth'">
                  <ko-secret-keys :data-obj.sync="form.data"></ko-secret-keys>
                </el-tab-pane>
                <el-tab-pane label="Labels/Annotations">
                  <ko-labels ref="ko_labels" :labelParentObj="form.metadata"></ko-labels>
                  <ko-annotations ref="ko_annotations" :annotationsParentObj="form.metadata"></ko-annotations>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor ref="yaml_editor" :value="yaml"></yaml-editor>
        </div>
        <div style="float: right;margin-top: 10px">
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.submit") }}
          </el-button>
        </div>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listNamespace} from "@/api/namespaces"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import KoSecretData from "@/components/ko-configuration/ko-secret-data"
import YamlEditor from "@/components/yaml-editor"
import {createSecret} from "@/api/secrets"
import KoSecretDockerData from "@/components/ko-configuration/ko-secret-docker-data"
import KoSecretKeys from "@/components/ko-configuration/ko-secret-keys"

export default {
  name: "SecretCreate",
  components: { KoSecretKeys, KoSecretDockerData, YamlEditor, KoSecretData, LayoutContent, KoAnnotations, KoLabels },
  props: {},
  data () {
    return {
      cluster: "",
      namespaces: [],
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "Secret",
        metadata: {
          name: "",
          namespace: "default",
          labels: {},
          annotations: {},
        },
        data: {},
        type: "Opaque"
      },
      showYaml: false,
      yaml: {},
      activeName: ""
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "Secrets" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      console.log(this.form)
      this.showYaml = false
    },
    transformYaml () {
      console.log(this.form)
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.form))
      // labels
      this.$refs.ko_labels.transformation(formData.metadata)
      // annotations
      this.$refs.ko_annotations.transformation(formData.metadata)
      return formData
    },
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      createSecret(this.cluster, this.form.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Secrets" })
      }).finally(() => {
        this.loading = false
      })
    },
    changeType () {
      this.form.data = {}
    }
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