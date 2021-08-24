<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: (toggleCase()+'s') }" v-loading="loading">
    <br>
    <div v-if="!showYaml">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="5">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace" :rules="selectRules">
              <ko-form-item :disabled="readOnly" :noClear="true" @change="changeNs" itemType="select2" :selections="namespace_list" v-model="form.metadata.namespace" />
            </el-form-item>
          </el-col>
          <el-col :span="7">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name" :rules="nameRules">
              <ko-form-item :disabled="readOnly" itemType="input" v-model="form.metadata.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="isReplicasShow()">
            <el-form-item :label="$t('business.workload.replicas')" prop="spec.replicas" :rules="numberRules">
              <ko-form-item :disabled="readOnly" placeholder="Any text you want that better describes this resource" itemType="number" v-model.number="form.spec.replicas" />
            </el-form-item>
          </el-col>
          <el-col :span="12" v-if="isCronJob()">
            <el-form-item :label="$t('business.workload.schedule')" prop="spec.schedule" :rules="requiredRules">
              <ko-form-item :disabled="readOnly" placeholder="0 * * * *" itemType="input" v-model="form.spec.schedule" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <ko-base :isReadOnly="readOnly" :baseParentObj="podSpec" @refreshContainer="refreshContainer" @gatherFormData="gatherFormData" @addContainer="addContainer" @deleteContainer="deleteContainer" />
          </el-col>
          <el-col :span="12" v-if="isStatefulSet()">
            <el-form-item :label="$t('business.network.service_name')" prop="spec.serviceName" :rules="selectRules">
              <ko-form-item :disabled="readOnly" :noClear="true" itemType="select2" v-model="form.spec.serviceName" :selections="service_list_of_ns" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <el-tabs :key="isRefresh" style="margin-top: 30px;background-color: #141418;" type="border-card" v-model="activeName">
        <el-tab-pane :label="$t('business.workload.general')" name="General">
          <ko-container :isReadOnly="readOnly" ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="currentContainer" :secretList="secret_list_of_ns" />
          <ko-ports :isReadOnly="readOnly" ref="ko_ports" :portParentObj="currentContainer" />
          <ko-command :isReadOnly="readOnly" ref="ko_command" :commandParentObj="currentContainer" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.health_check')" name="Health Check">
          <ko-health-check :isReadOnly="readOnly" ref="ko_health_readiness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.readiness_check')" />
          <ko-health-check :isReadOnly="readOnly" ref="ko_health_liveness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.liveness_check')" />
          <ko-health-check :isReadOnly="readOnly" ref="ko_health_startup_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.startup_check')" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.labels_annotations')" name="Labels/Annotations">
          <ko-labels :isReadOnly="readOnly" :label-obj.sync="form.metadata.labels" :labelTitle="$t('business.workload.label')" />
          <ko-annotations :isReadOnly="readOnly" :annotations-obj.sync="form.metadata.annotations" :annotationsTitle="$t('business.workload.annotations')" />
          <ko-labels :isReadOnly="readOnly" :label-obj.sync="podMetadata.labels" :labelTitle="'Pod ' + $t('business.workload.label')" />
          <ko-annotations :isReadOnly="readOnly" :annotations-obj.sync="podMetadata.annotations" :annotationsTitle="'Pod ' + $t('business.workload.annotations')" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.network')" name="Networking">
          <ko-networking :isReadOnly="readOnly" ref="ko_networking" :networkingParentObj="podSpec" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.schedule')" name="Scheduling">
          <ko-pod-scheduling :isReadOnly="readOnly" ref="ko_pod_scheduling" :podSchedulingParentObj="podSpec" :namespaceList="namespace_list" />
          <ko-node-scheduling :isReadOnly="readOnly" ref="ko_node_scheduling" :nodeSchedulingParentObj="podSpec" :nodeList="node_list" />
          <ko-tolerations ref="ko_toleration" :tolerationsParentObj="podSpec" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.resource')" name="Resources">
          <ko-resources :isReadOnly="readOnly" ref="ko_resource" :resourceParentObj="currentContainer" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.upgrade_policy')" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy :isReadOnly="readOnly" v-if="type === 'deployments'" ref="ko_upgrade_policy" :upgradePolicyParentObj="form.spec" />
          <ko-upgrade-policy-statefulset :isReadOnly="readOnly" v-if="type === 'statefulsets'" ref="ko_upgrade_policy_statefulset" :upgradePolicyParentObj="form.spec" />
          <ko-upgrade-policy-cronjob :isReadOnly="readOnly" v-if="type === 'cronjobs'" ref="ko_upgrade_policy_cronjob" :upgradePolicyParentObj="form.spec" />
          <ko-upgrade-policy-job :isReadOnly="readOnly" v-if="type === 'jobs'" ref="ko_upgrade_policy_job" :upgradePolicyParentObj="form.spec" />
          <ko-upgrade-policy-daemonset :isReadOnly="readOnly" v-if="type === 'daemonsets'" ref="ko_upgrade_policy_daemonset" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.security_context')" name="Security Context">
          <ko-security-context :isReadOnly="readOnly" ref="ko_security_context" :securityContextParentObj="currentContainer" />
        </el-tab-pane>

        <el-tab-pane :label="$t('business.workload.storage')" name="Storage">
          <ko-storage :isReadOnly="readOnly" ref="ko_storage" :storageParentObj="podSpec" :currentContainerIndex="currentContainerIndex" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>

        <el-tab-pane v-if="isStatefulSet()" :label="$t('business.workload.volume_claim_template')" name="Volume Claim Templates">
          <ko-volume-claim :isReadOnly="readOnly" ref="ko_volume_claim" :volumeClaimParentObj="form.spec" :currentNamespace="form.metadata.namespace" :pvcList="pvc_list" :scList="sc_list" />
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="grid-content bg-purple-light" v-if="showYaml">
      <yaml-editor :value="yaml" :is-edit="true" :readOnly="readOnly" ref="yaml_editor"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="operationLoading" @click="onSubmit()" type="primary">{{ $t("commons.button.confirm") }}</el-button>
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
import KoUpgradePolicy from "@/components/ko-workloads/ko-upgrade-policy.vue"
import KoUpgradePolicyStatefulset from "@/components/ko-workloads/ko-upgrade-policy-statefulset.vue"
import KoUpgradePolicyDaemonset from "@/components/ko-workloads/ko-upgrade-policy-daemonset.vue"
import KoUpgradePolicyCronjob from "@/components/ko-workloads/ko-upgrade-policy-cronjob.vue"
import KoUpgradePolicyJob from "@/components/ko-workloads/ko-upgrade-policy-job.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"
import KoStorage from "@/components/ko-workloads/ko-storage.vue"
import KoVolumeClaim from "@/components/ko-workloads/ko-volume-claim.vue"

import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"

import { getWorkLoadByName, createWorkLoad, updateWorkLoad } from "@/api/workloads"
import { listNamespace } from "@/api/namespaces"
import { listNodes } from "@/api/nodes"
import { listSecretsWithNs } from "@/api/secrets"
import { listConfigMapsWithNs } from "@/api/configmaps"
import { listStorageClasses } from "@/api/storageclass"
import { listNsServices } from "@/api/services"
import { listPvcs } from "@/api/pvc"

export default {
  components: {
    LayoutContent,
    KoFormItem,
    YamlEditor,
    KoBase,
    KoContainer,
    KoPorts,
    KoCommand,
    KoResources,
    KoHealthCheck,
    KoSecurityContext,
    KoNetworking,
    KoPodScheduling,
    KoNodeScheduling,
    KoTolerations,
    KoUpgradePolicy,
    KoUpgradePolicyStatefulset,
    KoUpgradePolicyDaemonset,
    KoUpgradePolicyCronjob,
    KoUpgradePolicyJob,
    KoLabels,
    KoAnnotations,
    KoStorage,
    KoVolumeClaim,
  },
  data() {
    return {
      name: "",
      type: "",
      operation: "",
      readOnly: false,
      // yaml
      showYaml: false,
      yaml: {},
      // containers
      currentContainerIndex: 0,
      currentContainerType: "standardContainers",
      currentContainer: {},
      // resources
      namespace_list: [],
      secret_list_of_ns: [],
      config_map_list_of_ns: [],
      node_list: [],
      sc_list: [],
      pvc_list: [],
      service_list_of_ns: [],
      // base form
      activeName: "General",
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "",
        kind: "",
        metadata: {
          name: "my-test",
          namespace: "",
        },
        spec: {
          replicas: 1,
          schedule: "",
          serviceName: "",
          template: {
            metadata: {},
            spec: {},
          },
        },
      },
      podSpec: {
        containers: [],
      },
      podMetadata: {},
      isRefresh: false,
      clusterName: "",
      loading: false,
      operationLoading: false,
      numberRules: [Rule.NumberRule],
      selectRules: [Rule.SelectRule],
      requiredRules: [Rule.RequiredRule],
      nameRules: [Rule.CommonNameRule],
    }
  },
  methods: {
    updateContainerList(val) {
      if (this.currentContainerType === "initContainers") {
        this.podSpec.initContainers[this.currentContainerIndex].name = val
      } else {
        this.podSpec.containers[this.currentContainerIndex].name = val
      }
    },
    search() {
      getWorkLoadByName(this.clusterName, this.type, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        if (!this.isCronJob()) {
          this.podSpec = this.form.spec.template.spec
          this.podMetadata = this.form.spec.template.metadata
        } else {
          this.podSpec = this.form.spec.jobTemplate.spec.template.spec
          this.podMetadata = this.form.spec.jobTemplate.spec.template.metadata
        }
        if (this.isStatefulSet()) {
          this.loadServicesWithNs()
        }
        this.loadSecretsWithNs()
        this.loadConfigMapsWithNs()
        this.yaml = res
        this.currentContainerIndex = 0
        this.currentContainer = this.podSpec.containers[0]
        this.isRefresh = !this.isRefresh
      })
    },
    loadNamespace() {
      this.namespace_list = []
      listNamespace(this.clusterName).then((res) => {
        for (const ns of res.items) {
          this.namespace_list.push(ns.metadata.name)
        }
      })
    },
    loadSecretsWithNs() {
      this.secret_list_of_ns = []
      listSecretsWithNs(this.clusterName, this.form.metadata.namespace).then((res) => {
        this.secret_list_of_ns = res.items
      })
    },
    loadStorageClass() {
      listStorageClasses(this.clusterName).then((res) => {
        this.sc_list = res.items
      })
    },
    loadPvcs() {
      listPvcs(this.clusterName).then((res) => {
        this.pvc_list = res.items
      })
    },
    loadServicesWithNs() {
      this.service_list_of_ns = []
      listNsServices(this.clusterName, this.form.metadata.namespace).then((res) => {
        res.items.forEach((item) => {
          this.service_list_of_ns.push(item.metadata.name)
        })
      })
    },
    loadConfigMapsWithNs() {
      this.config_map_list_of_ns = []
      listConfigMapsWithNs(this.clusterName, this.form.metadata.namespace).then((res) => {
        this.config_map_list_of_ns = res.items
      })
    },
    changeNs() {
      this.loadSecretsWithNs()
      this.loadConfigMapsWithNs()
      this.loadServicesWithNs()
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
        if (!this.podSpec.initContainers) {
          this.podSpec.initContainers = []
        }
        this.podSpec.initContainers.push(item)
      } else {
        this.podSpec.containers.push(item)
      }
    },
    deleteContainer(type, index) {
      if (type === "initContainers") {
        this.podSpec.initContainers.splice(index, 1)
      } else {
        this.podSpec.containers.splice(index, 1)
      }
    },
    gatherFormValid() {
      this.$refs["form"].validate((valid) => {
        if (!valid) {
          this.isValid = false
          this.unValidInfo = this.toggleCase() + this.$t("commons.validate.params_not_complete")
          return
        }
        this.isValid = true
      })
      if (!this.isValid) {
        return
      }
      if (!this.$refs.ko_container.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = this.$t("business.workload.container") + this.$t("commons.validate.params_not_complete")
        return
      }
      if (!this.$refs.ko_pod_scheduling.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = "Pod " + this.$t("business.workload.schedule") + this.$t("commons.validate.params_not_complete")
        return
      }
      if (!this.$refs.ko_storage.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = this.$t("business.workload.storage") + this.$t("commons.validate.params_not_complete")
        return
      }
    },
    newFormItem() {
      return {
        apiVersion: "",
        kind: this.form.kind,
        metadata: {
          name: this.form.metadata.name,
          namespace: this.form.metadata.namespace,
          labels: this.form.metadata.labels,
          annotations: this.form.metadata.annotations,
        },
        spec: {
          replicas: this.form.spec.replicas,
          schedule: this.form.spec.schedule,
          serviceName: this.form.spec.serviceName,
          template: {
            metadata: {},
            spec: {},
          },
        },
      }
    },
    gatherFormData() {
      var tempForm = this.newFormItem()
      var tempContainer = {}
      var tempPodSpec = { containers: this.podSpec.containers, initContainers: this.podSpec.initContainers }
      this.$refs.ko_container.transformation(tempContainer)
      this.$refs.ko_ports.transformation(tempContainer)
      this.$refs.ko_command.transformation(tempContainer)
      this.$refs.ko_resource.transformation(tempContainer)
      if (this.currentContainerType === "standardContainers") {
        this.$refs.ko_health_readiness_check.transformation(tempContainer)
        this.$refs.ko_health_liveness_check.transformation(tempContainer)
        this.$refs.ko_health_startup_check.transformation(tempContainer)
      }
      this.$refs.ko_security_context.transformation(tempContainer)
      this.$refs.ko_networking.transformation(tempPodSpec)
      this.$refs.ko_pod_scheduling.transformation(tempPodSpec)
      this.$refs.ko_node_scheduling.transformation(tempPodSpec)
      this.$refs.ko_toleration.transformation(tempPodSpec)
      switch (this.type) {
        case "deployments":
          tempForm.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy.transformation(tempForm.spec, tempPodSpec)
          break
        case "statefulsets":
          tempForm.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy_statefulset.transformation(tempForm.spec, tempPodSpec)
          this.$refs.ko_volume_claim.transformation(tempForm.spec)
          break
        case "cronjobs":
          tempForm.apiVersion = "batch/v1beta1"
          this.$refs.ko_upgrade_policy_cronjob.transformation(tempForm.spec, tempPodSpec)
          break
        case "jobs":
          tempForm.apiVersion = "batch/v1"
          this.$refs.ko_upgrade_policy_job.transformation(tempForm.spec, tempPodSpec)
          break
        case "daemonsets":
          tempForm.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy_daemonset.transformation(tempForm.spec, tempPodSpec)
          break
      }
      this.$refs.ko_storage.transformation(tempPodSpec)

      this.currentContainer = tempContainer
      if (this.currentContainerType === "initContainers") {
        tempPodSpec.initContainers[this.currentContainerIndex] = tempContainer
      } else {
        tempPodSpec.containers[this.currentContainerIndex] = tempContainer
      }
      if (!this.isReplicasShow()) {
        delete tempForm.spec.replicas
      }
      if (!this.isStatefulSet()) {
        delete tempForm.spec.serviceName
      }
      if (this.isCronJob()) {
        delete tempForm.spec.template
        tempForm.spec.jobTemplate = { spec: { template: { spec: tempPodSpec, metadata: this.podMetadata } } }
      } else {
        delete tempForm.spec.schedule
        tempForm.spec.template.spec = tempPodSpec
        tempForm.spec.template.metadata = this.podMetadata
      }
      this.podMetadata.labels = this.podMetadata.labels || { app: this.form.metadata.name }
      if (!this.isJob()) {
        tempForm.spec.selector = { matchLabels: this.podMetadata.labels }
      }
      this.podSpec = tempPodSpec
      let returnForm = JSON.parse(JSON.stringify(tempForm))
      this.form = returnForm
      return returnForm
    },
    isReplicasShow() {
      return this.type === "deployments" || this.type === "statefulsets"
    },
    isCronJob() {
      return this.type === "cronjobs"
    },
    isStatefulSet() {
      return this.type === "statefulsets"
    },
    isJob() {
      return this.type === "jobs"
    },
    onCancel() {
      this.$router.push({ name: this.toggleCase() + "s" })
    },
    onSubmit() {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        this.gatherFormValid()
        if (!this.isValid) {
          this.$notify({ title: this.$t("commons.message_box.prompt"), message: this.unValidInfo, type: "warning" })
          return
        }
        data = this.gatherFormData()
      }
      this.loading = true
      if (this.operation === "create") {
        createWorkLoad(this.clusterName, this.type, this.form.metadata.namespace, data)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.create_success"),
            })
            this.$router.push({ name: this.toggleCase() + "s" })
          })
          .finally(() => {
            this.loading = false
          })
      } else {
        updateWorkLoad(this.clusterName, this.type, this.form.metadata.namespace, this.form.metadata.name, data)
          .then(() => {
            this.$message({
              type: "success",
              message: this.$t("commons.msg.edit_success"),
            })
            this.$router.push({ name: this.toggleCase() + "s" })
          })
          .finally(() => {
            this.loading = false
          })
      }
    },
    onEditYaml() {
      this.gatherFormValid()
      if (!this.isValid) {
        this.$confirm(this.unValidInfo + this.$t("commons.confirm_message.open_yaml"), this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
          this.yaml = this.gatherFormData()
          this.showYaml = true
        })
      } else {
        this.yaml = this.gatherFormData()
        this.showYaml = true
      }
    },
    backToForm() {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    toggleCase() {
      switch (this.type) {
        case "deployments":
          return "Deployment"
        case "daemonsets":
          return "DaemonSet"
        case "statefulsets":
          return "StatefulSet"
        case "cronjobs":
          return "CronJob"
        case "jobs":
          return "Job"
      }
    },
  },
  mounted() {
    if (this.isStatefulSet) {
      this.loadStorageClass()
      this.loadPvcs()
    }
    this.loadNamespace()
    this.loadNodes()
  },
  created() {
    this.readOnly = this.$route.query.readOnly && this.$route.query.readOnly === "true"
    this.showYaml = this.$route.query.yamlShow === "true"
    this.clusterName = this.$route.query.cluster
    this.type = this.$route.path.split("/")[2]
    this.operation = this.$route.params.operation
    if (this.operation === "edit") {
      this.name = this.$route.params.name
      this.search()
    } else {
      this.podSpec = {
        containers: [
          {
            name: "container-0",
            image: "",
            imagePullPolicy: "IfNotPresent",
          },
        ],
        restartPolicy: "Always",
      }
      this.currentContainerIndex = 0
      this.currentContainer = this.podSpec.containers[0]
      this.isRefresh = !this.isRefresh
    }
    this.form.kind = this.toggleCase()
  },
}
</script>

<style scoped>
</style>