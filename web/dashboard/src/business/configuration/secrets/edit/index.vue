<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'Secrets'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="item">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="item.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="item.metadata.namespace" disabled>
                  <el-option :label="item.metadata.namespace"
                             :value="item.metadata.namespace"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.configuration.type')" required>
                <el-select v-model="item.type" disabled>
                  <el-option label="Opaque" value="Opaque"></el-option>
                  <!--                  <el-option label="Service Account Token	" value="kubernetes.io/service-account-token"></el-option>-->
                  <el-option label="Docker Registry" value="kubernetes.io/dockerconfigjson"></el-option>
                  <el-option :label="$t('business.configuration.basic_auth')"
                             value="kubernetes.io/basic-auth"></el-option>
                  <el-option :label="$t('business.configuration.ssh_auth')" value="kubernetes.io/ssh-auth"></el-option>
                  <el-option :label="$t('business.configuration.tls_auth')" value="kubernetes.io/tls"></el-option>
                  <!--                  <el-option :label="$t('business.configuration.token_auth')"-->
                  <!--                             value="bootstrap.kubernetes.io/token"></el-option>-->
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane key="data" label="Data" v-if="item.type==='Opaque'">
                  <ko-secret-data :dataObj.sync="item.data"></ko-secret-data>
                </el-tab-pane>
                <el-tab-pane key="data" label="Data" v-if="item.type==='kubernetes.io/dockerconfigjson'">
                  <ko-secret-docker-data :dataObj.sync="item.data"></ko-secret-docker-data>
                </el-tab-pane>
                <el-tab-pane label="Data" v-if="item.type==='kubernetes.io/ssh-auth'">
                  <ko-secret-keys :data-obj.sync="item.data"></ko-secret-keys>
                </el-tab-pane>
                <el-tab-pane label="Authentication" v-if="item.type==='kubernetes.io/basic-auth'">
                  <ko-secret-authentication :authentication-obj.sync="item.data"></ko-secret-authentication>
                </el-tab-pane>
                <el-tab-pane label="Tls" v-if="item.type==='kubernetes.io/tls'">
                  <ko-secret-certificate :certificate-obj.sync="item.data"></ko-secret-certificate>
                </el-tab-pane>
                <el-tab-pane name="1" label="Labels/Annotations">
                  <ko-labels ref="ko_labels" :labelParentObj="item.metadata"></ko-labels>
                  <ko-annotations ref="ko_annotations" :annotationsParentObj="item.metadata"></ko-annotations>
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
import {editSecret, getSecret} from "@/api/secrets"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import KoSecretData from "@/components/ko-configuration/ko-secret-data"
import YamlEditor from "@/components/yaml-editor"
import KoSecretDockerData from "@/components/ko-configuration/ko-secret-docker-data"
import KoSecretKeys from "@/components/ko-configuration/ko-secret-keys"
import KoSecretAuthentication from "@/components/ko-configuration/ko-secret-authentication"
import KoSecretCertificate from "@/components/ko-configuration/ko-secret-certificate"

export default {
  name: "SecretEdit",
  components: {
    KoSecretCertificate,
    KoSecretAuthentication,
    KoSecretKeys,
    KoSecretDockerData,
    YamlEditor,
    KoSecretData,
    LayoutContent,
    KoAnnotations,
    KoLabels
  },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      item: {
        metadata: {},
        type: ""
      },
      loading: false,
      cluster: "",
      showYaml: false,
      activeName: "",
      yaml: {}
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getSecret(this.cluster, this.namespace, this.name).then(res => {
        this.loading = false
        this.item = res
        if (this.showYaml) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
      })
    },
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
      this.showYaml = false
    },
    transformYaml () {
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.item))
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
      editSecret(this.cluster, this.form.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Secrets" })
      }).finally(() => {
        this.loading = false
      })
    },
  },
  watch: {
    yamlShow: function (newValue) {
      if (newValue) {
        this.yaml = this.transformYaml()
      }
      this.$router.push({
        path: "/" + this.namespace + "/secrets/edit/" + this.name,
        query: { yamlShow: newValue }
      })
      this.getDetail()
    }
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