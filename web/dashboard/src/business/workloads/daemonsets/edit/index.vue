<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'DaemonSets' }" v-loading="loading">
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
          <el-col :span="12">
            <el-row>
              <el-col :span="8">
                <el-select style="margin-top: 33px;width: 100%" @change="selectContainer(true)" v-model="selectContainerType">
                  <el-option v-for="(item, index) in type_list" :key="index" :label="item.label" :value="item.value" />
                </el-select>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('business.workload.container')">
                  <el-select v-if="selectContainerType === 'standardContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
                    <el-option v-for="(item, index) in form.spec.template.spec.containers" :key="index" :label="item.name" :value="index" />
                  </el-select>
                  <el-select v-if="selectContainerType === 'initContainers'" @change="selectContainer(false)" style="width:100%" v-model="selectContainerIndex">
                    <el-option v-for="(item, index) in form.spec.template.spec.initContainers" :key="index" :label="item.name" :value="index" />
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
            <ko-container ref="ko_container" @updateContanerList="updateContainerList" :containerParentObj="currentContainer" :secretList="secret_list_of_ns" />
            <ko-ports ref="ko_ports" :portParentObj="currentContainer" />
            <ko-command ref="ko_command" :commandParentObj="currentContainer" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Health Check" name="Health Check">
          <div :key="isRefresh">
            <ko-health-check ref="ko_health_readiness_check" :healthCheckParentObj="currentContainer" health_check_type="Readiness Check" />
            <ko-health-check ref="ko_health_liveness_check" :healthCheckParentObj="currentContainer" health_check_type="Liveness Check" />
            <ko-health-check ref="ko_health_startup_check" :healthCheckParentObj="currentContainer" health_check_type="Startup Check" />
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
          <ko-resources ref="ko_resource" :key="isRefresh" :resourceParentObj="currentContainer" />
        </el-tab-pane>
        <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy-daemonset ref="ko_upgrade_policy_daemonset" :key="isRefresh" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>
        <el-tab-pane label="Security Context" name="Security Context">
          <ko-security-context ref="ko_security_context" :key="isRefresh" :securityContextParentObj="currentContainer" />
        </el-tab-pane>
        <el-tab-pane label="Storage" name="Storage">
          <ko-storage :key="isRefresh" ref="ko_storage" :storageParentObj="form.spec.template.spec" :currentContainerIndex="currentContainerIndex" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
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
import KoUpgradePolicyDaemonset from "@/components/ko-workloads/ko-upgrade-policy-daemonset.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"
import KoStorage from "@/components/ko-workloads/ko-storage.vue"

import YamlEditor from "@/components/yaml-editor"

import { getDaemonSetByName, updateDaemonSet } from "@/api/daemonsets"
import { listNamespace } from "@/api/namespaces"
import { listNodes } from "@/api/nodes"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"
import Rule from "@/utils/rules"

export default {
  name: "DaemonSetEdit",
  components: { LayoutContent, KoFormItem, KoContainer, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoPodScheduling, KoNodeScheduling, KoTolerations, KoUpgradePolicyDaemonset, KoLabels, KoAnnotations, KoStorage, YamlEditor },
  data() {
    return {
      showYaml: false,
      type_list: [
        { label: "Init Container", value: "initContainers" },
        { label: "Standard Container", value: "standardContainers" },
      ],
      selectContainerIndex: 0,
      currentContainerIndex: 0,
      selectContainerType: "standardContainers",
      currentContainerType: "standardContainers",
      currentContainer: {},
      yaml: {},
      isRefresh: false,
      loading: false,
      namespace_list: [],
      secret_list: [],
      secret_list_of_ns: [],
      config_map_list: [],
      config_map_list_of_ns: [],
      node_list: [],
      containersList: [],
      activeName: "General",
      isValid: true,
      unValidInfo: "",
      form: {
        apiVersion: "apps/v1",
        kind: "DaemonSet",
        metadata: {
          name: "my-test-daemonset",
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
      },
      clusterName: "",
      operationLoading: false,
      numberRules: [Rule.NumberRule],
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
    search() {
      getDaemonSetByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        this.currentContainerIndex = 0
        this.currentContainer = this.form.spec.template.spec.containers[0]
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

    selectContainer(isChangeType) {
      this.gatherFormValid()
      if (!this.isValid) {
        this.$notify({ title: this.$t("commons.message_box.prompt"), message: this.unValidInfo })
        this.selectContainerIndex = this.currentContainerIndex
        this.selectContainerType = this.currentContainerType
        return
      }
      this.form = this.gatherFormData()
      this.selectContainerIndex = isChangeType ? 0 : this.selectContainerIndex
      if (this.selectContainerType === "initContainers") {
        if (!this.form.spec.template.spec.initContainers) {
          this.$confirm(this.$t("commons.confirm_message.add_init_container"), this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "warning",
          })
            .then(() => {
              this.form.spec.template.spec.initContainers = [{ name: "Container-0", image: "", imagePullPolicy: "ifNotPresent" }]
              this.currentContainerType = this.selectContainerType
              this.currentContainerIndex = isChangeType ? 0 : this.selectContainerIndex
              this.currentContainer = this.form.spec.template.spec.initContainers[this.currentContainerIndex]
              this.isRefresh = !this.isRefresh
            })
            .catch(() => {
              this.selectContainerType = this.currentContainerType
              this.selectContainerIndex = this.currentContainerIndex
            })
        } else {
          this.currentContainerType = this.selectContainerType
          this.currentContainerIndex = this.selectContainerIndex
          this.currentContainer = this.form.spec.template.spec.initContainers[this.currentContainerIndex]
          this.isRefresh = !this.isRefresh
        }
      } else {
        this.currentContainerType = this.selectContainerType
        this.currentContainerIndex = this.selectContainerIndex
        this.currentContainer = this.form.spec.template.spec.containers[this.currentContainerIndex]
        this.isRefresh = !this.isRefresh
      }
    },
    handleAddContainer() {
      this.dialogContainerVisible = false
      let itemContainer = {}
      itemContainer.image = ""
      itemContainer.imagePullPolicy = "ifNotPresent"
      if (this.selectContainerType === "initContainers") {
        if (this.form.spec.template.spec.initContainers) {
          itemContainer.name = "initContainer-" + this.form.spec.template.spec.initContainers.length.toString()
        } else {
          this.form.spec.template.spec.initContainers = []
          itemContainer.name = "initContainer-0"
        }
        this.form.spec.template.spec.initContainers.push(itemContainer)
      } else {
        itemContainer.name = "Container-" + this.form.spec.template.spec.containers.length.toString()
        this.form.spec.template.spec.containers.push(itemContainer)
      }
    },
    handleDeleteContainer() {
      if (this.selectContainerType === "initContainers") {
        if (this.form.spec.template.spec.initContainers.length <= this.currentContainerIndex) {
          return
        } else if (this.form.spec.template.spec.initContainers.length - 1 === this.currentContainerIndex) {
          delete this.form.spec.template.spec.initContainers
          this.selectContainerType = "standardContainers"
          this.currentContainerType = "standardContainers"
          this.currentContainer = this.form.spec.template.spec.containers[0]
        } else {
          this.form.spec.template.spec.initContainers.splice(this.currentContainerIndex, 1)
          this.currentContainer = this.form.spec.template.spec.initContainers[0]
        }
      } else {
        if (this.form.spec.template.spec.containers.length - 1 <= this.currentContainerIndex) {
          return
        } else {
          this.form.spec.template.spec.containers.splice(this.currentContainerIndex, 1)
          this.currentContainer = this.form.spec.template.spec.containers[0]
        }
      }
      this.currentContainerIndex = 0
      this.selectContainerIndex = 0
      this.isRefresh = !this.isRefresh
    },
    gatherFormValid() {
      this.isValid = true
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
      this.$refs.ko_container.transformation(this.currentContainer)
      this.$refs.ko_ports.transformation(this.currentContainer)
      this.$refs.ko_command.transformation(this.currentContainer)
      this.$refs.ko_resource.transformation(this.currentContainer)
      this.$refs.ko_health_readiness_check.transformation(this.currentContainer)
      this.$refs.ko_health_liveness_check.transformation(this.currentContainer)
      this.$refs.ko_health_startup_check.transformation(this.currentContainer)
      this.$refs.ko_security_context.transformation(this.currentContainer)
      this.$refs.ko_networking.transformation(this.form.spec.template.spec)
      this.$refs.ko_pod_scheduling.transformation(this.form.spec.template.spec)
      this.$refs.ko_upgrade_policy_daemonset.transformation(this.form.spec)
      this.$refs.ko_labels.transformation(this.form.metadata)
      this.$refs.ko_annotations.transformation(this.form.metadata)
      this.$refs.ko_pod_labels.transformation(this.form.spec.template.metadata)
      this.$refs.ko_pod_annotations.transformation(this.form.spec.template.metadata)
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
      updateDaemonSet(this.clusterName, data)
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
    this.selectContainerIndex = 0
    this.showYaml = this.$route.query.yamlShow === "true"
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