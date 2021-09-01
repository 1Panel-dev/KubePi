<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: (toggleCase()+'s') }" v-loading="loading">
    <br>
    <div v-if="!showYaml">
      <el-row :gutter="10" class="row-box">
        <el-col :span="spanWidth">
          <el-card class="el-card">
            <span style="font-size: 14px;font-weight: bold;">{{$t('business.workload.general')}}</span>
            <el-form label-position="top" ref="form" :model="form">
              <el-row>
                <el-form-item :label="$t('commons.table.name')" prop="metadata.name" :rules="nameRules">
                  <ko-form-item :disabled="readOnly || !isCreateOperation()" itemType="input" v-model="form.metadata.name" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace" :rules="selectRules">
                  <ko-select v-if="isCreateOperation()" style="width: 100%" @change="changeNs" :namespace.sync="form.metadata.namespace"></ko-select>
                  <ko-form-item v-else disabled :noClear="true" itemType="select2" :selections="namespace_list" v-model="form.metadata.namespace" />
                </el-form-item>
              </el-row>
              <el-row>
                <span style="line-height: 32px;font-size: 14px;color: #ebeef5;">{{$t('business.workload.annotations')}}</span>
                <ko-kv-table @updateKvDatas="updateKvDatas" kvType="annotation" :isReadOnly="readOnly" :dataObj="form.metadata.annotations" />
              </el-row>
              <el-row style="margin-top:10px; margin-bottom:10px">
                <span style="line-height: 32px;font-size: 14px;color: #ebeef5;">{{$t('business.workload.label')}}</span>
                <ko-kv-table @updateKvDatas="updateKvDatas" :resourceName="form.metadata.name" kvType="label" :isReadOnly="readOnly" :dataObj="form.metadata.labels" />
              </el-row>
              <el-row v-if="isReplicasShow()">
                <el-form-item :label="$t('business.workload.replicas')" prop="spec.replicas" :rules="numberRules">
                  <ko-form-item :disabled="readOnly" placeholder="Any text you want that better describes this resource" itemType="number" v-model.number="form.spec.replicas" />
                </el-form-item>
              </el-row>
              <el-row v-if="isStatefulSet()">
                <el-form-item :label="$t('business.network.service_name')" prop="spec.serviceName">
                  <ko-form-item :disabled="readOnly" itemType="select2" v-model="form.spec.serviceName" :selections="headless_service" />
                </el-form-item>
              </el-row>
              <el-row v-if="isCronJob()">
                <el-form-item :label="$t('business.workload.schedule')" prop="spec.schedule" :rules="requiredRules">
                  <ko-form-item :disabled="readOnly" placeholder="0 * * * *" itemType="input" v-model="form.spec.schedule" />
                </el-form-item>
              </el-row>
            </el-form>
          </el-card>
        </el-col>
        <el-col :span="spanWidth" v-if="isStatefulSet()">
          <el-card class="el-card">
            <span style="font-size: 14px;font-weight: bold;">{{$t('business.workload.volume_claim_template')}}</span>
            <ko-volume-claim @loadVolumes="loadVolumes" :key="isRefresh" :isReadOnly="readOnly" ref="ko_volume_claim" :volumeClaimParentObj="form.spec" :pvcList="pvc_list" :scList="sc_list" />
          </el-card>
        </el-col>
        <el-col :span="spanWidth">
          <el-card class="el-card">
            <span style="font-size: 14px;font-weight: bold;">{{$t('business.workload.volume')}}</span>
            <ko-volume @loadVolumes="loadVolumes" :key="isRefresh" :isReadOnly="readOnly" ref="ko_volume" :volumeParentObj="podSpec" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 20px" class="el-card">
        <span style="font-size: 14px;font-weight: bold;">{{$t('business.workload.spec')}}</span>
        <el-row :gutter="20">
          <el-col :span="18">
            <ko-base :isReadOnly="readOnly" :baseParentObj="podSpec" @refreshContainer="refreshContainer" @gatherFormData="gatherFormData" @addContainer="addContainer" @deleteContainer="deleteContainer" />
          </el-col>
        </el-row>
        <el-tabs :key="isRefresh" style="background-color: #141418;" type="border-card" v-model="activeName">
          <el-tab-pane :label="$t('business.workload.general')" name="General">
            <ko-container :isReadOnly="readOnly" ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="currentContainer" :secretList="secret_list_of_ns" />
            <ko-ports :isReadOnly="readOnly" ref="ko_ports" :portParentObj="currentContainer" />
            <ko-command :isReadOnly="readOnly" ref="ko_command" :commandParentObj="currentContainer" />
            <ko-environment :isReadOnly="readOnly" ref="ko_environment" :envParentObj="currentContainer" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
          </el-tab-pane>

          <el-tab-pane v-if="isStandardContainer()" :label="$t('business.workload.health_check')" name="Health Check">
            <ko-health-check :isReadOnly="readOnly" ref="ko_health_readiness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.readiness_check')" />
            <ko-health-check :isReadOnly="readOnly" ref="ko_health_liveness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.liveness_check')" />
            <ko-health-check :isReadOnly="readOnly" ref="ko_health_startup_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.startup_check')" />
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

          <el-tab-pane :label="$t('business.workload.mount')" name="Mount">
            <ko-volume-mount :isReadOnly="readOnly" ref="ko_volume_mount" :mountParentObj="currentContainer" :volumeList="volume_list" />
          </el-tab-pane>
        </el-tabs>
      </el-card>
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
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"

import KoKvTable from "@/components/ko-workloads/ko-kv-table.vue"
import KoBase from "@/components/ko-workloads/ko-base.vue"
import KoVolume from "@/components/ko-workloads/ko-volume.vue"
import KoVolumeClaim from "@/components/ko-workloads/ko-volume-claim.vue"
import KoContainer from "@/components/ko-workloads/ko-container.vue"
import KoPorts from "@/components/ko-workloads/ko-ports.vue"
import KoCommand from "@/components/ko-workloads/ko-command.vue"
import KoEnvironment from "@/components/ko-workloads/ko-environment.vue"
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
import KoVolumeMount from "@/components/ko-workloads/ko-volume-mount.vue"

import { getWorkLoadByName, createWorkLoad, updateWorkLoad } from "@/api/workloads"
import { listNamespace } from "@/api/namespaces"
import { listNodes } from "@/api/nodes"
import { listSecretsWithNs } from "@/api/secrets"
import { listConfigMapsWithNs } from "@/api/configmaps"
import { listStorageClasses } from "@/api/storageclass"
import { listNsServices } from "@/api/services"
import { listPvcs } from "@/api/pvc"
import KoSelect from "@/components/ko-select"

export default {
  components: {
    KoSelect,
    LayoutContent,
    KoFormItem,
    YamlEditor,
    KoKvTable,
    KoBase,
    KoVolumeClaim,
    KoVolume,
    KoContainer,
    KoPorts,
    KoCommand,
    KoEnvironment,
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
    KoVolumeMount,
  },
  data() {
    return {
      name: "",
      type: "",
      operation: "",
      readOnly: false,
      spanWidth: 12,
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
      headless_service: [],
      volume_list: [],
      // base form
      activeName: "General",
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "",
        kind: "",
        metadata: {
          name: "",
          namespace: "",
          labels: { "k8s.kubepi.cn/name": "" },
          annotations: {},
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
      if (this.isStandardContainer()) {
        this.podSpec.containers[this.currentContainerIndex].name = val
      } else {
        this.podSpec.initContainers[this.currentContainerIndex].name = val
      }
    },
    search() {
      this.loading = true
      getWorkLoadByName(this.clusterName, this.type, this.$route.params.namespace, this.$route.params.name)
        .then((res) => {
          this.form = res
          this.changeNs(this.form.metadata.namespace)
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
          this.loadVolumes("All", this.podSpec.volumes || [], this.form.spec.volumeClaimTemplates || [])
          this.isRefresh = !this.isRefresh
        })
        .finally(() => {
          this.loading = false
        })
    },
    loadNamespace() {
      this.namespace_list = []
      listNamespace(this.clusterName).then((res) => {
        for (const ns of res.items) {
          this.namespace_list.push(ns.metadata.name)
        }
        if (this.isCreateOperation()) {
          const p = sessionStorage.getItem("namespace")
          if (p !== null) {
            this.changeNs(p)
          } else {
            this.changeNs(this.namespace_list[0])
          }
        }
      })
    },
    loadSecretsWithNs(ns) {
      this.secret_list_of_ns = []
      listSecretsWithNs(this.clusterName, ns).then((res) => {
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
    loadServicesWithNs(ns) {
      this.headless_service = []
      listNsServices(this.clusterName, ns).then((res) => {
        res.items.forEach((item) => {
          if (item.spec.clusterIP === "None") {
            this.headless_service.push(item.metadata.name)
          }
        })
      })
    },
    loadConfigMapsWithNs(ns) {
      this.config_map_list_of_ns = []
      listConfigMapsWithNs(this.clusterName, ns).then((res) => {
        this.config_map_list_of_ns = res.items
      })
    },
    loadVolumes(type, volumes, volumeClaimTemplates) {
      if (type !== "All") {
        this.volume_list = this.volume_list.filter((item) => {
          return item.belongTo !== type
        })
        for (const vo of volumes) {
          this.volume_list.push({ name: vo.name, type: vo.type, belongTo: "volume" })
        }
        for (const vo of volumeClaimTemplates) {
          this.volume_list.push({ name: vo.name, type: "VolumeClaimTemplates", belongTo: "VolumeClaimTemplates" })
        }
      } else {
        this.volume_list = []
        for (const vo of volumes) {
          this.volume_list.push({ name: vo.name, type: this.loadVolumeType(vo), belongTo: "volume" })
        }
        for (const vo of volumeClaimTemplates) {
          this.volume_list.push({ name: vo.metadata.name, type: "VolumeClaimTemplates", belongTo: "VolumeClaimTemplates" })
        }
      }
    },
    loadVolumeType(volume) {
      if (volume.configMap) {
        return "ConfigMap"
      }
      if (volume.secret) {
        return "Secret"
      }
      if (volume.persistentVolumeClaim) {
        return "PVC"
      }
      if (volume.emptyDir) {
        return "EmptyDir"
      }
      if (volume.nfs) {
        return "NFS"
      }
      if (volume.hostPath) {
        return "HostPath"
      }
    },
    changeNs(val) {
      this.loadSecretsWithNs(val)
      this.loadConfigMapsWithNs(val)
      this.loadServicesWithNs(val)
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
    updateKvDatas(type, datas) {
      if (type === "label") {
        this.form.metadata.labels = datas
        this.podMetadata.labels = datas
        if (!this.isJob()) {
          this.form.spec.selector = { matchLabels: datas }
        }
      } else {
        this.form.metadata.annotations = datas
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
    },
    gatherFormData() {
      this.$refs.ko_volume.transformation(this.podSpec)
      this.$refs.ko_container.transformation(this.currentContainer)
      this.$refs.ko_ports.transformation(this.currentContainer)
      this.$refs.ko_command.transformation(this.currentContainer)
      this.$refs.ko_environment.transformation(this.currentContainer)
      this.$refs.ko_resource.transformation(this.currentContainer)
      if (this.isStandardContainer()) {
        this.$refs.ko_health_readiness_check.transformation(this.currentContainer)
        this.$refs.ko_health_liveness_check.transformation(this.currentContainer)
        this.$refs.ko_health_startup_check.transformation(this.currentContainer)
      }
      this.$refs.ko_security_context.transformation(this.currentContainer)
      this.$refs.ko_networking.transformation(this.podSpec)
      this.$refs.ko_pod_scheduling.transformation(this.podSpec)
      this.$refs.ko_node_scheduling.transformation(this.podSpec)
      this.$refs.ko_toleration.transformation(this.podSpec)
      switch (this.type) {
        case "deployments":
          this.form.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy.transformation(this.form.spec, this.podSpec)
          break
        case "statefulsets":
          this.form.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy_statefulset.transformation(this.form.spec, this.podSpec)
          this.$refs.ko_volume_claim.transformation(this.form.spec)
          break
        case "cronjobs":
          this.form.apiVersion = "batch/v1beta1"
          this.$refs.ko_upgrade_policy_cronjob.transformation(this.form.spec, this.podSpec)
          break
        case "jobs":
          this.form.apiVersion = "batch/v1"
          this.$refs.ko_upgrade_policy_job.transformation(this.form.spec, this.podSpec)
          break
        case "daemonsets":
          this.form.apiVersion = "apps/v1"
          this.$refs.ko_upgrade_policy_daemonset.transformation(this.form.spec, this.podSpec)
          break
      }
      this.$refs.ko_volume_mount.transformation(this.currentContainer)

      if (this.isStandardContainer()) {
        this.podSpec.containers[this.currentContainerIndex] = this.currentContainer
      } else {
        this.podSpec.initContainers[this.currentContainerIndex] = this.currentContainer
      }
      if (!this.isReplicasShow()) {
        delete this.form.spec.replicas
      }
      if (!this.isStatefulSet()) {
        delete this.form.spec.serviceName
      }
      if (this.isCronJob()) {
        delete this.form.spec.template
        this.form.spec.jobTemplate = { spec: { template: { spec: this.podSpec, metadata: this.podMetadata } } }
      } else {
        delete this.form.spec.schedule
        this.form.spec.template.spec = this.podSpec
        this.form.spec.template.metadata = this.podMetadata
      }
      return JSON.parse(JSON.stringify(this.form))
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
    isStandardContainer() {
      return this.currentContainerType === "standardContainers"
    },
    isCreateOperation() {
      return this.operation === "create"
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
      if (this.isCreateOperation()) {
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
              message: this.$t("commons.msg.update_success"),
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
    if (!this.isCreateOperation()) {
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
    this.spanWidth = this.isStatefulSet() ? 8 : 12
  },
}
</script>

<style scoped>
.row-box {
  display: flex;
  flex-flow: wrap;
}
.row-box .el-card {
  min-width: 100%;
  height: 100%;
  margin-right: 20px;
  border: 0px;
  box-shadow: 0 2px 12px 0 rgb(0 0 0 / 10%);
}
</style>
