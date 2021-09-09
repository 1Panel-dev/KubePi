<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'Services'}" v-loading="loading">
    <div>
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form" ref="form" :rules="rules">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="form.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.namespace.namespace')" required prop="metadata.namespace">
                <el-select v-model="form.metadata.namespace" disabled>
                  <el-option :label="form.metadata.namespace"
                             :value="form.metadata.namespace">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="3">
              <el-form-item :label="$t('business.configuration.type')" required>
                <el-select v-model="form.spec.type" @change="changeType" disabled>
                  <el-option label="Cluster IP" value="ClusterIP"></el-option>
                  <el-option label="External Name" value="ExternalName"></el-option>
                  <el-option label="Load Balancer" value="LoadBalancer"></el-option>
                  <el-option label="Node Port" value="NodePort"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick" v-if="Object.keys(form.spec).length!==0">
                <el-tab-pane name="ExternalName" label="External Name" v-if="form.spec.type==='ExternalName'">
                  <ko-service-external-name :external-name.sync="form.spec.externalName"></ko-service-external-name>
                </el-tab-pane>
                <el-tab-pane name="ServicePorts" :label="$t('business.network.service_port')" v-if="form.spec.type!=='ExternalName'">
                  <ko-service-ports :ports.sync="form.spec.ports"></ko-service-ports>
                </el-tab-pane>
                <el-tab-pane name="Selectors" :label="$t('business.network.selector')" v-if="form.spec.type!=='ExternalName'">
                  <ko-key-value :title="$t('business.network.selector')" :value.sync="form.spec.selector"></ko-key-value>
                </el-tab-pane>
                <el-tab-pane name="SessionAffinity" label="Session Affinity" v-if="form.spec.type!=='ExternalName'">
                  <ko-service-session-affinity :specObj="form.spec"></ko-service-session-affinity>
                </el-tab-pane>
                <el-tab-pane name="IPAddresses" :label="$t('business.network.ip_address')">
                  <ko-service-ip-addresses :specObj="form.spec"></ko-service-ip-addresses>
                </el-tab-pane>
                <el-tab-pane name="Labels" :label="$t('business.workload.labels_annotations')">
                  <ko-key-value :title="$t('business.workload.label')" :value.sync="form.metadata.labels"></ko-key-value>
                  <ko-key-value :title="$t('business.workload.annotations')" :value.sync="form.metadata.annotations"></ko-key-value>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor ref="yaml_editor" :is-edit="true" :value="yaml"></yaml-editor>
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
import {getService, updateService} from "@/api/services"
import KoServicePorts from "@/components/ko-network/service-ports"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import KoServiceIpAddresses from "@/components/ko-network/service-ip-addresses"
import KoServiceSessionAffinity from "@/components/ko-network/service-session-affinity"
import KoServiceExternalName from "@/components/ko-network/service-external-name"

export default {
  name: "ServiceEdit",
  components: {
    KoServiceExternalName,
    KoServiceSessionAffinity,
    KoServiceIpAddresses,
    KoKeyValue,
    KoServicePorts,
    YamlEditor,
    LayoutContent,
  },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      cluster: "",
      loading: false,
      form: {
        metadata: {},
        spec: {}
      },
      showYaml: false,
      yaml: {},
      activeName: "ServicePorts",
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        }
      }
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getService(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        if (this.showYaml) {
          this.yaml = this.transformYaml()
        }
      }).finally(() => {
        this.loading = false
      })
    },
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
      return JSON.parse(JSON.stringify(this.form))
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
      updateService(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Services" })
      }).finally(() => {
        this.loading = false
      })
    },
    changeType (type) {
      this.form.spec = {
        type: type
      }
      if (type === "ExternalName") {
        this.activeName = "ExternalName"
      } else {
        this.activeName = "ServicePorts"
      }
    }
  },
  watch: {
    showYaml: function (newValue) {
      this.$router.push({
        name: "ServiceEdit",
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
