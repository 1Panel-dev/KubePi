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
                       @tab-click="handleClick" v-if="Object.keys(item.metadata).length!==0">
                <el-tab-pane key="data" :label="$t('business.configuration.data')" v-if="item.type==='Opaque'">
                  <ko-secret-data :dataObj.sync="item.data"></ko-secret-data>
                </el-tab-pane>
                <el-tab-pane key="data" :label="$t('business.configuration.data')"
                             v-if="item.type==='kubernetes.io/dockerconfigjson'">
                  <ko-secret-docker-data :dataObj.sync="item.data"></ko-secret-docker-data>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.data')" v-if="item.type==='kubernetes.io/ssh-auth'">
                  <ko-secret-keys :data-obj.sync="item.data"></ko-secret-keys>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.data')" v-if="item.type==='kubernetes.io/basic-auth'">
                  <ko-secret-authentication :authentication-obj.sync="item.data"></ko-secret-authentication>
                </el-tab-pane>
                <el-tab-pane label="TLS" v-if="item.type==='kubernetes.io/tls'">
                  <ko-secret-certificate :certificate-obj.sync="item.data"></ko-secret-certificate>
                </el-tab-pane>
                <el-tab-pane name="1" :label="$t('business.workload.labels_annotations')">
                  <ko-key-value :title="$t('business.workload.label')"
                                :value.sync="item.metadata.labels"></ko-key-value>
                  <ko-key-value :title="$t('business.workload.annotations')"
                                :value.sync="item.metadata.annotations"></ko-key-value>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor ref="yaml_editor" :is-edit="true" :value="yaml"></yaml-editor>
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
import KoSecretData from "@/components/ko-configuration/ko-secret-data"
import YamlEditor from "@/components/yaml-editor"
import KoSecretDockerData from "@/components/ko-configuration/ko-secret-docker-data"
import KoSecretKeys from "@/components/ko-configuration/ko-secret-keys"
import KoSecretAuthentication from "@/components/ko-configuration/ko-secret-authentication"
import KoSecretCertificate from "@/components/ko-configuration/ko-secret-certificate"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "SecretEdit",
  components: {
    KoKeyValue,
    KoSecretCertificate,
    KoSecretAuthentication,
    KoSecretKeys,
    KoSecretDockerData,
    YamlEditor,
    KoSecretData,
    LayoutContent,
  },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      item: {
        metadata: {},
        type: "",
        data: {}
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
          this.yaml = this.transformYaml()
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
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.item))
    },
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      editSecret(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Secrets" })
      }).finally(() => {
        this.loading = false
      })
    },
  },
  watch: {
    showYaml: function (newValue) {
      this.$router.push({
        name: "SecretEdit",
        params: {
          name: this.name,
          namespace: this.namespace,
        },
        query: { yamlShow: newValue }
      })
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
