<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Deployments' }" v-loading="loading">
    <br>
    <div v-if="showYaml === 'false'">
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
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.replicas')" prop="spec.replicas" :rules="numberRules">
              <ko-form-item placeholder="Any text you want that better describes this resource" clearable itemType="number" v-model="form.spec.replicas" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-top: 10px">
          <el-col :span="12">
            <el-row>
              <el-col :span="20">
                <el-form-item :label="$t('business.workload.container')">
                  <el-select @change="selectContainer" style="width:100%" itemType="select" v-model="selectContainerIndex">
                    <el-option v-for="(item, index) in form.spec.template.spec.containers" :key="index" :label="item.name" :value="index" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="2">
                <div style="margin-top: 33px">
                  <el-button style="width:100%" @click="handleAddContainer">+</el-button>
                </div>
              </el-col>
              <el-col :span="2">
                <div style="margin-top: 33px">
                  <el-button style="width:100%" @click="handleDeleteContainer">-</el-button>
                </div>
              </el-col>
            </el-row>
          </el-col>
        </el-row>
      </el-form>

      <el-tabs style="margin-top: 30px;background-color: #141418;" type="border-card" v-model="activeName">
        <el-tab-pane label="General" name="General">
          <div :key="isRefresh">
            <ko-container ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="form.spec.template.spec.containers[currentContainerIndex]" :secretList="secret_list_of_ns" />
            <ko-ports ref="ko_ports" :portParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
            <ko-command ref="ko_command" :commandParentObj="form.spec.template.spec.containers[currentContainerIndex]" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Health Check" name="Health Check">
          <div :key="isRefresh">
            <ko-health-check ref="ko_health_readiness_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Readiness Check" />
            <ko-health-check ref="ko_health_liveness_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Liveness Check" />
            <ko-health-check ref="ko_health_startup_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Startup Check" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Labels/Annotations" name="Labels/Annotations">
          <div :key="isRefresh">
            <ko-labels ref="ko_labels" :labelParentObj="form.metadata" labelTitle="Labels" />
            <ko-annotations ref="ko_annotations" :annotationsParentObj="form.metadata" annotationsTitle="Annotations" />
            <ko-labels ref="ko_pod_labels" :labelParentObj="form.spec.template.metadata" labelTitle="Pod Labels" />
            <ko-annotations ref="ko_pod_annotations" :annotationsParentObj="form.spec.template.metadata" annotationsTitle="Pod Annotations" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Networking" name="Networking">
          <ko-networking ref="ko_networking" :key="isRefresh" :networkingParentObj="form.spec.template.spec" />
        </el-tab-pane>
        <el-tab-pane label="Scheduling" name="Scheduling">
          <div :key="isRefresh">
            <ko-pod-scheduling ref="ko_pod_scheduling" :podSchedulingParentObj="form.spec.template.spec" :namespaceList="namespace_list" />
            <ko-node-scheduling ref="ko_node_scheduling" :nodeSchedulingParentObj="form.spec.template.spec" :nodeList="node_list" />
            <ko-tolerations ref="ko_toleration" :tolerationsParentObj="form.spec.template.spec" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Resources" name="Resources">
          <ko-resources ref="ko_resource" :key="isRefresh" :resourceParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy ref="ko_upgrade_policy" :key="isRefresh" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>
        <el-tab-pane label="Security Context" name="Security Context">
          <ko-security-context ref="ko_security_context" :key="isRefresh" :securityContextParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Storage" name="Storage">
          <ko-storage :key="isRefresh" ref="ko_storage" :storageParentObj="form.spec.template.spec" :currentContainerIndex="currentContainerIndex" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="grid-content bg-purple-light" v-if="showYaml === 'true'">
      <yaml-editor :value="yaml"></yaml-editor>
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
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"
import KoStorage from "@/components/ko-workloads/ko-storage.vue"

import YamlEditor from "@/components/yaml-editor"

import { getDeploymentByName, createDeployment } from "@/api/workloads"
import { listNamespace } from "@/api/namespaces"
import { listNodes } from "@/api/nodes"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"
import Rule from "@/utils/rules"

export default {
  name: "DeploymentEdit",
  components: { LayoutContent, KoFormItem, KoContainer, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoPodScheduling, KoNodeScheduling, KoTolerations, KoUpgradePolicy, KoLabels, KoAnnotations, KoStorage, YamlEditor },
  data() {
    return {
      showYaml: false,
      yaml: {},
      isRefresh: false,
      operation: "",
      loading: false,
      namespace_list: [],
      secret_list: [],
      secret_list_of_ns: [],
      config_map_list: [],
      config_map_list_of_ns: [],
      node_list: [],
      activeName: "General",
      selectContainerIndex: 0,
      currentContainerIndex: 0,
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "apps/v1",
        kind: "Deployment",
        metadata: {
          name: "my-test-deployment",
          namespace: "",
        },
        spec: {
          replicas: 1,
          template: {
            metadata: {},
            spec: {
              containers: [],
              restartPolicy: "Always",
            },
          },
        },
        type: "apps.deployment",
      },
      clusterName: "",
      operationLoading: false,
      numberRules: [Rule.NumberRule],
      requiredRules: [Rule.RequiredRule],
    }
  },
  methods: {
    updateContainerList(val) {
      this.form.spec.template.spec.containers[this.currentContainerIndex].name = val
    },
    search() {
      getDeploymentByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        this.yaml = res
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

    selectContainer() {
      this.gatherFormValid()
      if (!this.isValid) {
        this.$notify({ title: this.$t("commons.message_box.prompt"), message: this.unValidInfo })
        this.selectContainerIndex = this.currentContainerIndex
        return
      }
      this.form = this.gatherFormData()
      this.currentContainerIndex = this.selectContainerIndex
      this.isRefresh = !this.isRefresh
    },
    handleAddContainer() {
      this.form.spec.template.spec.containers.push({
        name: "Container-" + this.form.spec.template.spec.containers.length.toString(),
        image: "",
        imagePullPolicy: "ifNotPresent",
      })
    },
    handleDeleteContainer() {
      if (this.form.spec.template.spec.containers.length <= 1) {
        return
      }
      this.form.spec.template.spec.containers.splice(this.currentContainerIndex, 1)
      this.currentContainerIndex = 0
    },
    gatherFormValid() {
      if (!this.$refs.ko_container.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = "ko_container 参数不完整"
        return
      }
      if (!this.$refs.ko_storage.checkIsValid()) {
        this.isValid = false
        this.unValidInfo = "ko_storage 参数不完整"
        return
      }
    },
    gatherFormData() {
      this.$refs.ko_container.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_ports.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_command.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_resource.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_health_readiness_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_health_liveness_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_health_startup_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_security_context.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_networking.transformation(this.form.spec.template.spec)
      this.$refs.ko_pod_scheduling.transformation(this.form.spec.template.spec)
      this.$refs.ko_upgrade_policy.transformation(this.form.spec)
      this.$refs.ko_labels.transformation(this.form.metadata)
      this.$refs.ko_annotations.transformation(this.form.metadata)
      this.$refs.ko_pod_labels.transformation(this.form.spec.template.metadata)
      this.$refs.ko_pod_annotations.transformation(this.form.spec.template.metadata)
      this.$refs.ko_storage.transformation(this.form.spec.template.spec)
      return this.form
    },
    onCancel() {
      this.$router.push({ name: "Deployments" })
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
      createDeployment(this.clusterName, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "Deployments" })
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
    this.selectContainerIndex = 0
    this.showYaml = this.$route.query.yamlShow
    this.clusterName = this.$route.query.cluster
    if (this.$route.params.name) {
      this.search()
    }
    this.loadNamespace()
    this.loadSecrets()
    this.loadConfigMaps()
    this.loadNodes()
  },
}
</script>

<style scoped>
</style>