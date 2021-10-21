<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'Services'}" v-loading="loading">
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
                <ko-select  :namespace.sync="form.metadata.namespace"></ko-select>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.configuration.type')" required>
                <el-select v-model="type" @change="changeType">
                  <el-option label="Cluster IP" value="ClusterIP"></el-option>
                  <el-option label="Headless" value="Headless"></el-option>
                  <el-option label="External Name" value="ExternalName"></el-option>
                  <el-option label="Load Balancer" value="LoadBalancer"></el-option>
                  <el-option label="Node Port" value="NodePort"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick" ref=tabs>
                <el-tab-pane name="ExternalName" label="External Name" v-if="form.spec.type==='ExternalName'">
                  <ko-service-external-name :external-name.sync="form.spec.externalName"></ko-service-external-name>
                </el-tab-pane>
                <el-tab-pane name="ServicePorts"  :label="$t('business.network.service_port')" v-if="form.spec.type!=='ExternalName'">
                  <ko-service-ports :type="type" :ports.sync="form.spec.ports"></ko-service-ports>
                </el-tab-pane>
                <el-tab-pane name="Selectors" :label="$t('business.network.selector')" v-if="form.spec.type!=='ExternalName' ">
                  <ko-key-value :title="$t('business.network.selector')" :value.sync="form.spec.selector"></ko-key-value>
                </el-tab-pane>
                <el-tab-pane name="SessionAffinity" label="Session Affinity" v-if="form.spec.type!=='ExternalName'">
                  <ko-service-session-affinity :specObj.sync="form.spec"></ko-service-session-affinity>
                </el-tab-pane>
                <el-tab-pane name="IPAddresses" :label="$t('business.network.ip_address')">
                  <ko-service-ip-addresses :type="type" :specObj="form.spec"></ko-service-ip-addresses>
                </el-tab-pane>
                <el-tab-pane name="Labels" :label="$t('business.workload.labels_annotations')">
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
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"
import {createService} from "@/api/services"
import KoServicePorts from "@/components/ko-network/service-ports"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import KoServiceIpAddresses from "@/components/ko-network/service-ip-addresses"
import KoServiceSessionAffinity from "@/components/ko-network/service-session-affinity"
import KoServiceExternalName from "@/components/ko-network/service-external-name"
import KoSelect from "@/components/ko-select"

export default {
  name: "ServiceCreate",
  components: {
    KoSelect,
    KoServiceExternalName,
    KoServiceSessionAffinity,
    KoServiceIpAddresses,
    KoKeyValue,
    KoServicePorts,
    YamlEditor,
    LayoutContent,
  },
  props: {},
  data () {
    return {
      cluster: "",
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "Service",
        metadata: {
          name: "",
          namespace: "",
          labels: {},
          annotations: {},
        },
        spec: {
          type: "ClusterIP"
        },
      },
      showYaml: false,
      yaml: undefined,
      activeName: "ServicePorts",
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        }
      },
      type: "ClusterIP"
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.name
    },
    onCancel () {
      this.$router.push({ name: "Services" })
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
      let data = JSON.parse(JSON.stringify(this.form))
      if (data.spec.type === 'ExternalName') {
        delete data.spec.ports
        delete data.spec.sessionAffinity
        delete data.spec.sessionAffinityConfig
        delete data.spec.clusterIP
      }else {
        delete data.spec.externalName
      }
      return data
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
      createService(this.cluster, data.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Services" })
      }).finally(() => {
        this.loading = false
      })
    },
    changeType (type) {

      if (type === 'Headless') {
        this.form.spec.type  = "ClusterIP"
        this.form.spec.clusterIP = 'None'
      }else {
        this.form.spec.type  = type
      }
      if (type === 'ExternalName') {
        this.activeName = "ExternalName"
      } else {
        this.activeName = "ServicePorts"
      }
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
