<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: (toggleCase()+'s') }" v-loading.fullscreen.lock="loading">
    <br>
    <div v-if="!showYaml">
      <el-row :gutter="10" class="row-box">
        <el-col :span="spanWidth">
          <el-card class="el-card">
            <span style="font-size: 14px;font-weight: bold;">{{ $t("business.workload.general") }}</span>
            <el-form label-position="top" ref="form" :model="form">
              <el-row>
                <el-form-item :label="$t('commons.table.name')" prop="metadata.name" :rules="nameRules">
                  <ko-form-item :disabled="readOnly || !isCreateOperation()" itemType="input" v-model="form.metadata.name" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace" :rules="selectRules">
                  <ko-select v-if="isCreateOperation()" style="width: 100%" @change="changeNs" :namespace.sync="form.metadata.namespace"></ko-select>
                  <ko-form-item v-else disabled itemType="select2" :selections="[]" v-model="form.metadata.namespace" />
                </el-form-item>
              </el-row>

              <el-row>
                <ko-kv-table @changeCreateMode="changeCreateMode" :key="isRefresh" ref="ko_annotation_label" :isReadOnly="readOnly" :resourceName="form.metadata.name" :metadataObj="form.metadata" :selectorObj="form.spec.selector" :isCreate="isCreateOperation()" />
              </el-row>

              <el-row v-if="isReplicasShow()">
                <el-form-item :label="$t('business.workload.replicas')" prop="spec.replicas" :rules="numberRules">
                  <ko-form-item :disabled="readOnly" itemType="number" v-model.number="form.spec.replicas" />
                </el-form-item>
              </el-row>
              <el-row v-if="isStatefulSet()">
                <el-form-item :label="$t('business.network.service_name')" prop="spec.serviceName">
                  <ko-form-item :disabled="readOnly" itemType="select2" v-model="form.spec.serviceName" :selections="headless_service_of_ns" />
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
            <span style="font-size: 14px;font-weight: bold;">{{ $t("business.workload.volume_claim_template") }}</span>
            <ko-volume-claim @loadVolumes="loadVolumes" :key="isRefresh" :isReadOnly="readOnly" ref="ko_volume_claim" :volumeClaimParentObj="form.spec" :pvcList="pvc_list_of_ns" :scList="sc_list" />
          </el-card>
        </el-col>
        <el-col :span="spanWidth">
          <el-card class="el-card">
            <span style="font-size: 14px;font-weight: bold;">{{ $t("business.workload.volume") }}</span>
            <ko-volume @loadVolumes="loadVolumes" :key="isRefresh" :isReadOnly="readOnly" ref="ko_volume" :volumeParentObj="podSpec" :configMapList="config_map_list_of_ns" :pvcList="pvc_list_of_ns" :secretList="secret_list_of_ns" />
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 20px;border: 0" class="el-card">
        <span style="font-size: 14px;font-weight: bold;">{{ $t("business.workload.containers") }}</span>
        <el-tabs tabPosition="left" v-model="activeName">
          <el-tab-pane label="Spec" name="Spec">
            <el-tabs :key="isRefresh" style="background-color: #141418;" type="border-card" v-model="activeNameSpec">
              <el-tab-pane :label="$t('business.workload.upgrade_policy')" name="Scaling/Upgrade Policy">
                <ko-upgrade-job v-if="isCronJob() || isJob()" :isReadOnly="readOnly" ref="ko_upgrade_job" :upgradePolicyParentObj="form.spec" :resourceType="type" />
                <ko-upgrade-common v-else :isReadOnly="readOnly" ref="ko_upgrade_common" :upgradePolicyParentObj="form.spec" :resourceType="type" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.node_schedule')" name="NodeSchedule">
                <ko-node-selector :isReadOnly="readOnly" ref="ko_node_scheduling" :nodeSchedulingParentObj="podSpec" :nodeList="node_list" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.affinity_or_anti')" name="PodSchedule">
                <ko-pod-scheduling :isReadOnly="readOnly" ref="ko_pod_scheduling" :podSchedulingParentObj="podSpec" :namespaceList="namespace_list" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.toleration')" name="Toleration">
                <ko-tolerations :isReadOnly="readOnly" ref="ko_toleration" :tolerationsParentObj="podSpec" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.network')" name="Networking">
                <ko-networking :isReadOnly="readOnly" ref="ko_networking" :networkingParentObj="podSpec" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.security_context')" name="Security Context">
                <ko-spec-security :isReadOnly="readOnly" ref="ko_spec_security" :securityContextParentObj="podSpec" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.others')" name="others">
                <ko-spec-base :isReadOnly="readOnly" ref="ko_spec_base" :resourceType="type" :specBaseParentObj="podSpec" :serviceList="service_of_ns" :secretList="secret_list_of_ns" :repoList="repo_list" />
              </el-tab-pane>
            </el-tabs>
          </el-tab-pane>

          <el-tab-pane label="Containers" name="Container">
            <el-row :gutter="20">
              <el-col :span="18">
                <ko-base :isReadOnly="readOnly" :baseParentObj="podSpec" @refreshContainer="refreshContainer" @gatherFormData="gatherFormData" @addContainer="addContainer" @deleteContainer="deleteContainer" />
              </el-col>
            </el-row>
            <el-tabs :key="isRefresh" style="background-color: #141418;" type="border-card" v-model="activeNameContainers">
              <el-tab-pane :label="$t('business.workload.general')" name="General">
                <ko-container :isReadOnly="readOnly" ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="currentContainer" :containerType="currentContainerType" :metadata="podMetadata" :repoList="repo_list" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.command')" name="Command">
                <ko-command :isReadOnly="readOnly" ref="ko_command" :commandParentObj="currentContainer" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.port')" name="Port">
                <ko-ports :isReadOnly="readOnly" ref="ko_ports" :portParentObj="currentContainer" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.environment_variable')" name="Environment">
                <ko-environment :isReadOnly="readOnly" ref="ko_environment" :envParentObj="currentContainer" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
              </el-tab-pane>

              <el-tab-pane v-if="isStandardContainer()" :label="$t('business.workload.health_check')" name="Health Check">
                <ko-health-check :isReadOnly="readOnly" ref="ko_health_readiness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.readiness_check')" />
                <ko-health-check :isReadOnly="readOnly" ref="ko_health_liveness_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.liveness_check')" />
                <ko-health-check :isReadOnly="readOnly" ref="ko_health_startup_check" :healthCheckParentObj="currentContainer" :health_check_type="$t('business.workload.startup_check')" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.resource')" name="Resources">
                <ko-resources :isReadOnly="readOnly" ref="ko_resource" :resourceParentObj="currentContainer" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.security_context')" name="Security Context">
                <ko-security-context :isReadOnly="readOnly" ref="ko_security_context" :securityContextParentObj="currentContainer" />
              </el-tab-pane>

              <el-tab-pane :label="$t('business.workload.mount')" name="Mount">
                <ko-volume-mount :isReadOnly="readOnly" ref="ko_volume_mount" :mountParentObj="currentContainer" :volumeList="volume_list" />
              </el-tab-pane>
            </el-tabs>
          </el-tab-pane>

          <el-tab-pane label="Service" name="Service" v-if="hasService()">
            <ko-service-add ref="service_add" :serviceObj="serviceForm" />
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
        <el-button v-loading="operationLoading" v-if="!readOnly" @click="onSubmit()" type="primary">
          {{ $t("commons.button.confirm") }}
        </el-button>
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

import KoSpecBase from "@/components/ko-workloads/ko-spec/ko-spec-base.vue"
import KoNodeSelector from "@/components/ko-workloads/ko-spec/ko-node-selector.vue"
import KoPodScheduling from "@/components/ko-workloads/ko-spec/ko-pod-scheduling.vue"
import KoTolerations from "@/components/ko-workloads/ko-spec/ko-tolerations.vue"
import KoNetworking from "@/components/ko-workloads/ko-spec/ko-networking.vue"
import KoSpecSecurity from "@/components/ko-workloads/ko-spec/ko-spec-security.vue"
import KoUpgradeCommon from "@/components/ko-workloads/ko-spec/ko-upgrade-common.vue"
import KoUpgradeJob from "@/components/ko-workloads/ko-spec/ko-upgrade-job.vue"

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
import KoVolumeMount from "@/components/ko-workloads/ko-volume-mount.vue"

import KoServiceAdd from "@/components/ko-workloads/ko-service/ko-service-add.vue"

import { listClusterReposDetail } from "../../../../kubepi/src/api/clusters"
import { getWorkLoadByName, createWorkLoad, updateWorkLoad, deleteWorkLoad } from "@/api/workloads"
import { listSecretsWithNs, createSecret, editSecret } from "@/api/secrets"
import { listConfigMapsWithNs } from "@/api/configmaps"
import { listStorageClasses } from "@/api/storageclass"
import { listServicesWithNs } from "@/api/services"
import { listPvcsWithNs } from "@/api/pvc"
import { listNodes } from "@/api/nodes"
import { getNamespaces } from "@/api/auth"
import KoSelect from "@/components/ko-select"

import { checkPermissions } from "@/utils/permission"

export default {
  components: {
    KoSelect,
    LayoutContent,
    KoFormItem,
    YamlEditor,
    KoKvTable,
    KoSpecBase,
    KoSpecSecurity,
    KoUpgradeCommon,
    KoUpgradeJob,
    KoNodeSelector,

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
    KoTolerations,
    KoVolumeMount,

    KoServiceAdd,
  },
  data() {
    return {
      name: "",
      type: "",
      operation: "",
      recreate: false,
      readOnly: false,
      spanWidth: 12,
      // yaml
      showYaml: false,
      yaml: undefined,
      // containers
      currentContainerIndex: 0,
      currentContainerType: "standardContainers",
      currentContainer: {},
      // resources
      namespace_list: [],
      node_list: [],
      secret_list_of_ns: [],
      repo_list: [],
      config_map_list_of_ns: [],
      sc_list: [],
      pvc_list_of_ns: [],
      headless_service_of_ns: [],
      service_of_ns: [],
      volume_list: [],
      // base form
      activeName: "Container",
      activeNameSpec: "Scaling/Upgrade Policy",
      activeNameContainers: "General",
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "",
        kind: "",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          replicas: 1,
          schedule: "",
          serviceName: "",
          selector: {},
          template: {
            metadata: {},
            spec: {},
          },
        },
      },
      serviceForm: null,
      batchCreateForm: {},
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
          if (res.metadata.resourceVersion) {
            delete res.metadata.resourceVersion
          }
          this.form = res
          this.changeNs(this.form.metadata.namespace)
          if (!this.isCronJob()) {
            this.podSpec = this.form.spec.template.spec
            this.podMetadata = this.form.spec.template.metadata
          } else {
            this.podSpec = this.form.spec.jobTemplate.spec.template.spec
            this.podMetadata = this.form.spec.jobTemplate.spec.template.metadata
          }
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
    loadSecretsWithNs(ns) {
      this.secret_list_of_ns = []
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "secrets", verb: "list" })) {
        return
      }
      listSecretsWithNs(this.clusterName, ns).then((res) => {
        this.secret_list_of_ns = res.items
      })
    },
    loadStorageClass() {
      this.sc_list = []
      if (
        !checkPermissions({
          scope: "cluster",
          apiGroup: "storage.k8s.io",
          resource: "storageclasses",
          verb: "list",
        })
      ) {
        return
      }
      listStorageClasses(this.clusterName).then((res) => {
        this.sc_list = res.items
      })
    },
    loadPvcsWithNs(ns) {
      this.pvc_list_of_ns = []
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "persistentvolumeclaims", verb: "list" })) {
        return
      }
      listPvcsWithNs(this.clusterName, ns).then((res) => {
        this.pvc_list_of_ns = res.items
      })
    },
    loadServicesWithNs(ns) {
      this.headless_service_of_ns = []
      this.service_of_ns = []
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "services", verb: "list" })) {
        return
      }
      listServicesWithNs(this.clusterName, ns).then((res) => {
        res.items.forEach((item) => {
          this.service_of_ns.push(item.metadata.name)
          if (item.spec.clusterIP === "None") {
            this.headless_service_of_ns.push(item.metadata.name)
          }
        })
      })
    },
    loadConfigMapsWithNs(ns) {
      this.config_map_list_of_ns = []
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "configmaps", verb: "list" })) {
        return
      }
      listConfigMapsWithNs(this.clusterName, ns).then((res) => {
        this.config_map_list_of_ns = res.items
      })
    },
    loadRepos() {
      this.repo_list = []
      listClusterReposDetail(this.clusterName).then((res) => {
        this.repo_list = res.data
      })
    },
    loadNodes() {
      this.node_list = []
      if (!checkPermissions({ scope: "cluster", apiGroup: "", resource: "nodes", verb: "list" })) {
        return
      }
      listNodes(this.clusterName).then((res) => {
        this.node_list = res.items
      })
    },
    loadNamespace() {
      getNamespaces(this.clusterName).then((res) => {
        this.namespace_list = res.data
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
          this.volume_list.push({
            name: vo.metadata.name,
            type: "VolumeClaimTemplates",
            belongTo: "VolumeClaimTemplates",
          })
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
      if (val) {
        this.loadSecretsWithNs(val)
        this.loadConfigMapsWithNs(val)
        this.loadPvcsWithNs(val)
        if (this.isStatefulSet()) {
          this.loadServicesWithNs(val)
          this.loadStorageClass()
        }
      }
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
    changeCreateMode(val) {
      this.recreate = val
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
      this.$refs.ko_annotation_label.transformation(this.form, this.podMetadata)

      this.$refs.ko_volume.transformation(this.podSpec)
      this.$refs.ko_networking.transformation(this.podSpec)
      this.$refs.ko_pod_scheduling.transformation(this.podSpec)
      this.$refs.ko_node_scheduling.transformation(this.podSpec)
      this.$refs.ko_toleration.transformation(this.podSpec)
      this.$refs.ko_spec_security.transformation(this.podSpec)
      this.$refs.ko_container.transformation(this.currentContainer, this.podMetadata)
      this.$refs.ko_spec_base.transformation(this.podSpec, this.podMetadata)
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

      if (this.isCronJob()) {
        delete this.form.spec.template
        this.form.apiVersion = "batch/v1"
        this.form.spec.jobTemplate = { spec: {} }
        this.$refs.ko_upgrade_job.transformation(this.form.spec)
      } else if (this.isJob()) {
        this.form.apiVersion = "batch/v1"
        delete this.form.spec.selector
        this.$refs.ko_upgrade_job.transformation(this.form.spec)
      } else {
        this.form.apiVersion = "apps/v1"
        this.$refs.ko_upgrade_common.transformation(this.form.spec)
        if (this.isStatefulSet()) {
          this.$refs.ko_volume_claim.transformation(this.form.spec)
        }
      }
      this.$refs.ko_volume_mount.transformation(this.currentContainer)

      if (!this.isStatefulSet()) {
        delete this.form.spec.serviceName
      }
      if (this.isStandardContainer()) {
        this.podSpec.containers[this.currentContainerIndex] = this.currentContainer
      } else {
        this.podSpec.initContainers[this.currentContainerIndex] = this.currentContainer
      }
      if (!this.isReplicasShow()) {
        delete this.form.spec.replicas
      }

      if (this.isCronJob()) {
        delete this.form.spec.template
        this.form.spec.jobTemplate.spec.template = { spec: this.podSpec, metadata: this.podMetadata }
      } else {
        delete this.form.spec.schedule
        this.form.spec.template.spec = this.podSpec
        this.form.spec.template.metadata = this.podMetadata
      }
      if (this.hasService()) {
        this.serviceForm = this.$refs.service_add.transformation(this.form.metadata)
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
    hasService() {
      return !this.isCronJob() && !this.isJob() && this.isCreateOperation()
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
        this.onCreate(data)
      } else {
        if (!this.recreate) {
          this.onEdit(data)
        } else {
          this.onRecreate(data)
        }
      }
    },
    onCreate(data) {
      var backUrl = this.toggleCase() + "s"
      if (data.kind === "List") {
        this.batchCreateForm = data
      } else {
        this.batchCreateForm = { apiVersion: "v1", kind: "List", items: [] }
        this.batchCreateForm.items.push(data)
        if (this.serviceForm) {
          this.batchCreateForm.items.push(this.serviceForm)
        }
      }
      let ps = []
      ps = this.handleSecret(ps)
      for (const item of this.batchCreateForm.items) {
        ps.push(createWorkLoad(this.clusterName, item.kind.toLowerCase() + "s", item.metadata.namespace, item))
      }

      Promise.all(ps)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.loading = false
          this.$router.push({ name: backUrl })
        })
        .catch(() => {
          this.$message({
            type: "error",
            message: this.$t("commons.msg.create_failed"),
          })
          this.loading = false
        })
    },
    handleSecret(ps) {
      if (!this.podMetadata.labels) {
        return ps
      }
      let operation = this.podMetadata.labels["operation"]
      for (const key in this.podMetadata.labels) {
        if (key.indexOf("kubepi-repo-") !== -1) {
          let repoName = this.podMetadata.labels[key]
          let secretName = "kubepi-" + repoName + "-secret"
          let exist = false
          for (const item of this.secret_list_of_ns) {
            if (item.metadata.name === secretName) {
              exist = true
              break
            }
          }
          if (exist) {
            if (operation === "update") {
              let secret = this.addSecret(repoName, this.form.metadata.namespace, secretName)
              ps.push(editSecret(this.clusterName, this.form.metadata.namespace, secretName, secret))
            }
          } else {
            let secret = this.addSecret(repoName, this.form.metadata.namespace, secretName)
            ps.push(createSecret(this.clusterName, this.form.metadata.namespace, secret))
          }
        }
      }
      return ps
    },
    addSecret(repoName, namespace, secretName) {
      let repoInfo = null
      for (const item of this.repo_list) {
        if (item.name === repoName) {
          repoInfo = item
          break
        }
      }
      if (repoInfo === null) {
        return
      }
      const auths = {
        auths: {
          [repoInfo.downloadUrl]: {
            username: repoInfo.credential.username,
            password: repoInfo.credential.password,
          },
        },
      }
      const { Base64 } = require("js-base64")
      const data = {
        [".dockerconfigjson"]: Base64.encode(JSON.stringify(auths)),
      }
      return {
        apiVersion: "v1",
        kind: "Secret",
        metadata: {
          name: secretName,
          namespace: namespace,
        },
        data: data,
        type: "kubernetes.io/dockerconfigjson",
      }
    },
    onEdit(data) {
      let ps = []
      if (this.podMetadata.labels && this.podMetadata.labels["operation"] === "update") {
        ps = this.handleSecret(ps)
      }
      Promise.all(ps).then(() => {
        updateWorkLoad(this.clusterName, this.type, data.metadata.namespace, data.metadata.name, data)
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
      })
    },
    onRecreate(data) {
      deleteWorkLoad(this.clusterName, this.type, data.metadata.namespace, data.metadata.name)
        .then(() => {
          delete data.metadata.resourceVersion
          this.onCreate(data)
        })
        .finally(() => {
          this.loading = false
        })
    },
    onEditYaml() {
      this.gatherFormValid()
      if (!this.isValid) {
        this.$confirm(this.unValidInfo + this.$t("commons.confirm_message.open_yaml"), this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
          let data = this.gatherFormData()
          if (this.isCreateOperation() && this.serviceForm !== null) {
            this.batchCreateForm = { apiVersion: "v1", kind: "List", items: [] }
            this.batchCreateForm.items.push(data)
            this.batchCreateForm.items.push(this.serviceForm)
            this.yaml = this.batchCreateForm
          } else {
            this.yaml = data
          }
          this.showYaml = true
        })
      } else {
        let data = this.gatherFormData()
        if (this.isCreateOperation() && this.serviceForm !== null) {
          this.batchCreateForm = { apiVersion: "v1", kind: "List", items: [] }
          this.batchCreateForm.items.push(data)
          this.batchCreateForm.items.push(this.serviceForm)
          this.yaml = this.batchCreateForm
        } else {
          this.yaml = data
        }
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
  created() {
    this.readOnly = this.$route.query.readOnly && this.$route.query.readOnly === "true"
    this.showYaml = this.$route.query.yamlShow === "true"
    this.clusterName = this.$route.query.cluster
    this.type = this.$route.path.split("/")[2]
    this.operation = this.$route.params.operation
    this.loadNodes()
    this.loadRepos()
    this.loadNamespace()
    if (!this.isCreateOperation()) {
      this.name = this.$route.params.name
      this.search()
    } else {
      this.form.metadata.labels = { "kubepi.org/name": "" }
      this.podSpec = {
        containers: [
          {
            name: "container-0",
            image: "",
            imagePullPolicy: "IfNotPresent",
          },
        ],
      }
      this.podSpec.restartPolicy = this.isCronJob() || this.isJob() ? "Never" : "Always"
      this.currentContainerIndex = 0
      this.currentContainer = this.podSpec.containers[0]
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
