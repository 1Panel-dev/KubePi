<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'Secrets'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
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
            <el-col :span="3">
              <el-form-item :label="$t('business.configuration.type')" required>
                <el-select v-model="form.type" @change="changeType">
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
                <el-tab-pane :label="$t('business.configuration.data')" v-if="form.type==='Opaque'">
                  <ko-secret-data :dataObj.sync="form.data"></ko-secret-data>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.data')"
                             v-if="form.type==='kubernetes.io/dockerconfigjson'">
                  <ko-secret-docker-data :dataObj.sync="form.data"></ko-secret-docker-data>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.data')" v-if="form.type==='kubernetes.io/ssh-auth'">
                  <ko-secret-keys :data-obj.sync="form.data"></ko-secret-keys>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.data')" v-if="form.type==='kubernetes.io/basic-auth'">
                  <ko-secret-authentication
                          :authentication-obj.sync="form.data"></ko-secret-authentication>
                </el-tab-pane>
                <el-tab-pane label="TLS" v-if="form.type==='kubernetes.io/tls'">
                  <ko-secret-certificate :certificate-obj.sync="form.data"></ko-secret-certificate>
                </el-tab-pane>
                <el-tab-pane name="1" :label="$t('business.workload.labels_annotations')">
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
        <ko-alert :messages="messages"></ko-alert>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoSecretData from "@/components/ko-configuration/ko-secret-data"
import YamlEditor from "@/components/yaml-editor"
import {createSecret} from "@/api/secrets"
import KoSecretDockerData from "@/components/ko-configuration/ko-secret-docker-data"
import KoSecretKeys from "@/components/ko-configuration/ko-secret-keys"
import KoSecretAuthentication from "@/components/ko-configuration/ko-secret-authentication"
import KoSecretCertificate from "@/components/ko-configuration/ko-secret-certificate"
import KoAlert from "@/components/ko-alert"
import Rule from "@/utils/rules"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import KoSelect from "@/components/ko-select"

export default {
  name: "SecretCreate",
  components: {
    KoKeyValue,
    KoAlert,
    KoSecretCertificate,
    KoSecretAuthentication,
    KoSecretKeys,
    YamlEditor,
    KoSecretData,
    LayoutContent,
    KoSecretDockerData,
    KoSelect,
  },
  props: {},
  data () {
    return {
      cluster: "",
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "Secret",
        metadata: {
          name: "",
          namespace: "",
          labels: {},
          annotations: {},
        },
        data: {},
        type: "Opaque"
      },
      showYaml: false,
      yaml: undefined,
      activeName: "",
      messages: [],
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        }
      }
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
      createSecret(this.cluster, data.metadata.namespace, data).then(() => {
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
    this.showYaml = this.$route.query.yamlShow === "true"
  }
}
</script>

<style scoped>

</style>
