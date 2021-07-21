<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'DaemonSets' }" v-loading="loading">
    <br>
    <div v-if="!showYaml">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="5">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-form-item @change="changeNs" itemType="select2" :selections="namespace_list" v-model="form.metadata.namespace" />
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name" :rules="requiredRules">
              <ko-form-item itemType="input" v-model="form.metadata.name" />
            </el-form-item>
          </el-col>
        </el-row>
        <ko-base :baseParentObj="form.spec.template.spec" @refreshContainer="refreshContainer" @gatherFormData="gatherFormData" @addContainer="addContainer" @deleteContainer="deleteContainer" />
      </el-form>

      <el-tabs :key="isRefresh" style="margin-top: 30px;background-color: #141418;" type="border-card" v-model="activeName">
        <el-tab-pane label="General" name="General">
          <ko-container ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="currentContainer" :secretList="secret_list_of_ns" />
          <ko-ports ref="ko_ports" :portParentObj="currentContainer" />
          <ko-command ref="ko_command" :commandParentObj="currentContainer" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>
        <el-tab-pane label="Health Check" name="Health Check">
          <ko-health-check ref="ko_health_readiness_check" :healthCheckParentObj="currentContainer" health_check_type="Readiness Check" />
          <ko-health-check ref="ko_health_liveness_check" :healthCheckParentObj="currentContainer" health_check_type="Liveness Check" />
          <ko-health-check ref="ko_health_startup_check" :healthCheckParentObj="currentContainer" health_check_type="Startup Check" />
        </el-tab-pane>
        <el-tab-pane label="Labels/Annotations" name="Labels/Annotations">
          <ko-labels ref="ko_labels" :label-obj.sync="form.metadata.labels" labelTitle="Labels" />
          <ko-annotations ref="ko_annotations" :annotations-obj.sync="form.metadata.annotations" annotationsTitle="Annotations" />
          <ko-labels ref="ko_pod_labels" :label-obj.sync="form.spec.template.metadata.labels" labelTitle="Pod Labels" />
          <ko-annotations ref="ko_pod_annotations" :annotations-obj.sync="form.spec.template.metadata.annotations" annotationsTitle="Pod Annotations" />
        </el-tab-pane>
        <el-tab-pane label="Networking" name="Networking">
          <ko-networking ref="ko_networking" :networkingParentObj="form.spec.template.spec" />
        </el-tab-pane>
        <el-tab-pane label="Scheduling" name="Scheduling">
          <ko-pod-scheduling ref="ko_pod_scheduling" :podSchedulingParentObj="form.spec.template.spec" :namespaceList="namespace_list" />
          <ko-node-scheduling ref="ko_node_scheduling" :nodeSchedulingParentObj="form.spec.template.spec" :nodeList="node_list" />
          <ko-tolerations ref="ko_toleration" :tolerationsParentObj="form.spec.template.spec" />
        </el-tab-pane>
        <el-tab-pane label="Resources" name="Resources">
          <ko-resources ref="ko_resource" :resourceParentObj="currentContainer" />
        </el-tab-pane>
        <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy-daemonset ref="ko_upgrade_policy_daemonset" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>
        <el-tab-pane label="Security Context" name="Security Context">
          <ko-security-context ref="ko_security_context" :securityContextParentObj="currentContainer" />
        </el-tab-pane>
        <el-tab-pane label="Storage" name="Storage">
          <ko-storage ref="ko_storage" :storageParentObj="form.spec.template.spec" :currentContainerIndex="currentContainerIndex" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="grid-content bg-purple-light" v-if="showYaml">
      <yaml-editor :value="yaml"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="operationLoading" @click="onSubmit()" type="primary">
          {{ $t("commons.button.create") }}
        </el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoFormItem from "@/components/ko-form-item/index"
import KoBase from "@/components/ko-workloads/ko-base.vue"
import KoContainer from "@/components/ko-workloads/ko-container.vue"
import KoPorts from "@/components/ko-workloads/ko-ports.vue"
import KoCommand from "@/components/ko-workloads/ko-command.vue"
import KoResources from "@/components/ko-workloads/ko-resources.vue"
import KoHealthCheck from "@/components/ko-workloads/ko-health-check.vue"
import KoSecurityContext from "@/components/ko-workloads/ko-security-context.vue"
import KoNetworking from "@/components/ko-workloads/ko-networking.vue"
import KoPodScheduling from "@/components/ko-workloads/ko-pod-scheduling.vue"
import KoNodeScheduling from "@/components/ko-workloads/ko-node-scheduling.vue"
import KoTolerations from "@/components/ko-workloads/ko-tolerations.vue"
import KoUpgradePolicyDaemonset from "@/components/ko-workloads/ko-upgrade-policy-daemonset.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"
import KoStorage from "@/components/ko-workloads/ko-storage.vue"

import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"

import { createDaemonSet } from "@/api/daemonsets"
import { listNamespace } from "@/api/namespaces"
import { listNodes } from "@/api/nodes"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"

export default {
  name: "DaemonSetCreate",
  components: { LayoutContent, KoFormItem, YamlEditor, KoBase, KoContainer, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoPodScheduling, KoNodeScheduling, KoTolerations, KoUpgradePolicyDaemonset, KoLabels, KoAnnotations, KoStorage },
  data() {
    return {
      // yaml
      showYaml: false,
      yaml: {},
      // containers
      currentContainerIndex: 0,
      currentContainerType: "standardContainers",
      currentContainer: {},
      // resources
      namespace_list: [],
      secret_list: [],
      secret_list_of_ns: [],
      config_map_list: [],
      config_map_list_of_ns: [],
      node_list: [],
      // base form
      activeName: "General",
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "apps/v1",
        kind: "DaemonSet",
        metadata: {
          name: "my-test-deployment",
          namespace: "kube-operator",
        },
        spec: {
          replicas: 1,
          template: {
            metadata: {},
            spec: {
              containers: [
                {
                  name: "Container-0",
                  image: "",
                  imagePullPolicy: "ifNotPresent",
                },
              ],
              restartPolicy: "Always",
            },
          },
        },
      },
      isRefresh: false,
      loading: false,
      clusterName: "",
      operationLoading: false,
      requiredRules: [Rule.RequiredRule],
    }
  },
  methods: {
    updateContainerList(val) {
      if (this.currentContainerType === "initContainers") {
        this.form.spec.template.spec.initContainers[this.currentContainerIndex].name = val
      } else {
        this.form.spec.template.spec.containers[this.currentContainerIndex].name = val
      }
    },
    loadNamespace() {
      this.namespace_list = []
      listNamespace(this.clusterName).then((res) => {
        for (const ns of res.items) {
          this.namespace_list.push(ns.metadata.name)
        }
      })
    },
    loadSecrets() {
      this.secret_list = []
      listSecrets(this.clusterName).then((res) => {
        this.secret_list = res.items
        this.loadSecretsWithNs()
      })
    },
    loadSecretsWithNs() {
      this.secret_list_of_ns = []
      for (const s of this.secret_list) {
        if (s.metadata.namespace === this.form.metadata.namespace) {
          this.secret_list_of_ns.push(s)
        }
      }
    },
    loadConfigMaps() {
      this.config_map_list = []
      listConfigMaps(this.clusterName).then((res) => {
        this.config_map_list = res.items
        this.loadConfigMapsWithNs()
      })
    },
    loadConfigMapsWithNs() {
      this.config_map_list_of_ns = []
      for (const c of this.config_map_list) {
        if (c.metadata.namespace === this.form.metadata.namespace) {
          this.config_map_list_of_ns.push(c)
        }
      }
    },
    changeNs() {
      this.loadSecretsWithNs()
      this.loadConfigMapsWithNs()
    },
    loadNodes() {
      this.node_list = []
      listNodes(this.clusterName).then((res) => {
        this.node_list = res.items
      })
    },

    refreshContainer(type, index, item) {
      this.currentContainerIndex = index
      this.currentContainerType = type
      this.currentContainer = item
      this.isRefresh = !this.isRefresh
    },
    addContainer(type, item) {
      if (type === "initContainers") {
        if (!this.form.spec.template.spec.initContainers) {
          this.form.spec.template.spec.initContainers = []
        }
        this.form.spec.template.spec.initContainers.push(item)
      } else {
        this.form.spec.template.spec.containers.push(item)
      }
    },
    deleteContainer(type, index) {
      if (type === "initContainers") {
        this.form.spec.template.spec.initContainers.splice(index, 1)
      } else {
        this.form.spec.template.spec.containers.splice(index, 1)
      }
    },
    gatherFormValid() {
      this.isValid = true
      this.$refs["form"].validate((valid) => {
        if (!valid) {
          return
        }
      })
      if (!this.$refs.ko_container.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = "ko_container 参数不完整"
        return
      }
    },
    gatherFormData() {
      this.$refs.ko_container.transformation(this.currentContainer)
      this.$refs.ko_ports.transformation(this.currentContainer)
      this.$refs.ko_command.transformation(this.currentContainer)
      this.$refs.ko_resource.transformation(this.currentContainer)
      if (this.currentContainerType === "standardContainers") {
        this.$refs.ko_health_readiness_check.transformation(this.currentContainer)
        this.$refs.ko_health_liveness_check.transformation(this.currentContainer)
        this.$refs.ko_health_startup_check.transformation(this.currentContainer)
      }
      this.$refs.ko_security_context.transformation(this.currentContainer)
      this.$refs.ko_networking.transformation(this.form.spec.template.spec)
      this.$refs.ko_node_scheduling.transformation(this.form.spec.template.spec)
      this.$refs.ko_pod_scheduling.transformation(this.form.spec.template.spec)
      this.$refs.ko_toleration.transformation(this.form.spec.template.spec)
      this.$refs.ko_upgrade_policy_daemonset.transformation(this.form.spec)
      this.$refs.ko_storage.transformation(this.form.spec.template.spec)
      if (this.currentContainerType === "initContainers") {
        this.form.spec.template.spec.initContainers[this.currentContainerIndex] = this.currentContainer
      } else {
        this.form.spec.template.spec.containers[this.currentContainerIndex] = this.currentContainer
      }
      return this.form
    },
    onCancel() {
      this.$router.push({ name: "DaemonSets" })
    },
    onSubmit() {
      let data = {}
      this.gatherFormValid()
      if (!this.isValid) {
        this.$notify({ title: this.$t("commons.message_box.prompt"), message: this.unValidInfo })
        return
      }
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.gatherFormData()
      }
      this.loading = true
      createDaemonSet(this.clusterName, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "DaemonSets" })
        })
        .finally(() => {
          this.loading = false
        })
    },
    onEditYaml() {
      this.yaml = this.gatherFormData()
      this.showYaml = true
    },
    backToForm() {
      this.showYaml = false
    },
  },
  mounted() {
    this.showYaml = this.$route.params.showYaml === "true"
    this.clusterName = this.$route.query.cluster
    this.currentContainer = this.form.spec.template.spec.containers[0]
    this.loadNamespace()
    this.loadSecrets()
    this.loadConfigMaps()
    this.loadNodes()
    this.isRefresh = !this.isRefresh
  },
}
</script>

<style scoped>
</style>
