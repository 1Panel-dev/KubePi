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
                <el-select v-model="form.spec.type" @change="changeType">
                  <el-option label="Cluster IP" value="ClusterIP"></el-option>
                  <el-option label="External Name" value="ExternalName"></el-option>
                  <el-option label="Load Balancer" value="LoadBalancer"></el-option>
                  <el-option label="Node Port" value="NodePort"></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane label="Service Ports">
                  <ko-service-ports :ports.sync="form.spec.ports"></ko-service-ports>
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
import {listNamespace} from "@/api/namespaces"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"
import {createService} from "@/api/services"
import KoServicePorts from "@/components/ko-network/service-ports"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "ServiceCreate",
  components: {
    KoKeyValue,
    KoServicePorts,
    YamlEditor,
    LayoutContent,
  },
  props: {},
  data () {
    return {
      cluster: "",
      namespaces: [],
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "Service",
        metadata: {
          name: "",
          namespace: "default",
          labels: {},
          annotations: {},
        },
        spec: {
          type: "ClusterIP"
        },
      },
      showYaml: false,
      yaml: {},
      activeName: "",
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
        this.$router.push({ name: "Secrets" })
      }).finally(() => {
        this.loading = false
      })
    },
    changeType (type) {
      this.form.spec = {
        type: type
      }
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
