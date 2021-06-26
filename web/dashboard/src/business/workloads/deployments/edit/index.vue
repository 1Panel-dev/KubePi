<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Deployments' }" v-loading="loading">
    <br>
    <div v-if="!showYaml">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-row>
            <el-col :span="10">
              <ko-form-item labelName="Namespace" clearable itemType="select" :selections="namespace_list" v-model="form.metadata.namespace" />
            </el-col>
            <el-col :span="14">
              <ko-form-item labelName="Name" clearable itemType="input" v-model="form.metadata.name" />
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Replicas" placeholder="Any text you want that better describes this resource" clearable itemType="number" v-model="form.spec.replicas" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 10px">
        <el-col :span="12">
          <el-row>
            <el-col :span="20">
              <ko-form-item @change="selectContainer" labelName="Container" itemType="select" v-model="currentIndex" :selections="container_lists" />
            </el-col>
            <el-col :span="2">
              <div style="margin-top: 16px">
                <el-button style="width:100%" size="mini" @click="handleAddContainer">+</el-button>
              </div>
            </el-col>
            <el-col :span="2">
              <div style="margin-top: 16px">
                <el-button style="width:100%" size="mini" @click="handleDeleteContainer">-</el-button>
              </div>
            </el-col>
          </el-row>
        </el-col>
      </el-row>

      <el-tabs style="margin-top: 30px;background-color: #141418;" v-model="activeName">
        <el-tab-pane label="General" name="General">
          <div :key="isRefresh">
            <ko-container ref="ko_general" :containerParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
            <ko-ports ref="ko_ports" :portParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
            <ko-command ref="ko_command" :commandParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Resources" name="Resources">
          <ko-resources ref="ko_resource" :key="isRefresh" :resourceParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Health Check" name="Health Check">
          <div :key="isRefresh">
            <ko-health-check ref="ko_health_readiness_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Readiness Check" health_check_helper="Containers will be removed from service endpoints when this check is failing. Recommended." />
            <ko-health-check ref="ko_health_liveness_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Liveness Check" health_check_helper="Containers will be restarted when this check is failing. Not recommended for most uses." />
            <ko-health-check ref="ko_health_startup_check" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" style="margin-top=30px" health_check_type="Startup Check" health_check_helper="Containers will wait until this check succeeds before attempting other health checks." />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Security Context" name="Security Context">
          <ko-security-context ref="ko_security_context" :securityContextParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Networking" name="Networking">
          <ko-networking ref="ko_networking" :key="isRefresh" :networkingParentObj="form.spec.template.spec" />
        </el-tab-pane>
        <el-tab-pane label="Scheduling" name="Scheduling">
          <div :key="isRefresh">
            <ko-pod-scheduling ref="ko_pod_scheduling" :podSchedulingParentObj="form.spec.template.spec" />
            <ko-node-scheduling ref="ko_node_scheduling" :podSchedulingParentObj="form.spec.template.spec" />
            <ko-tolerations ref="ko_tolerations" :tolerationsParentObj="form.spec.template.spec" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy ref="ko_upgrade_policy" :key="isRefresh" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>
        <el-tab-pane label="Labels/Annotations" name="Labels/Annotations">
          <div :key="isRefresh">
            <ko-labels ref="ko_labels" :labelParentObj="form.spec.template.metadata" />
            <ko-annotations ref="ko_annotations" :annotationsParentObj="form.spec.template.metadata" />
          </div>
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
import KoUpgradePolicy from "@/components/ko-workloads/ko-upgrade-policy.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"

import YamlEditor from "@/components/yaml-editor"
import { getDeploymentByName, createDeployment } from "@/api/workloads"

export default {
  name: "DeploymentForm",
  components: { LayoutContent, KoFormItem, KoContainer, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoPodScheduling, KoNodeScheduling, KoTolerations, KoUpgradePolicy, KoLabels, KoAnnotations, YamlEditor },
  data() {
    return {
      showYaml: false,
      yaml: {},
      isRefresh: false,
      operation: "",
      loading: false,
      namespace_list: [
        { label: "kube-system", value: "kube-system" },
        { label: "kube-public", value: "kube-public" },
        { label: "kube-operator", value: "kube-operator" },
        { label: "default", value: "default" },
      ],
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "ifNotPresent", value: "ifNotPresent" },
        { label: "Never", value: "Never" },
      ],
      container_lists: [{ label: "Container-0", value: 0 }],
      activeName: "General",
      description: "",
      currentIndex: 0,
      currentContainerIndex: 0,
      form: {
        apiVersion: "apps/v1",
        kind: "Deployment",
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
                  name: "",
                  image: "",
                  imagePullPolicy: "ifNotPresent",
                },
              ],
              restartPolicy: "Always",
            },
          },
        },
        type: "apps.deployment",
      },
      operationLoading: false,
    }
  },
  methods: {
    search() {
      getDeploymentByName(this.$route.params.cluster, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        this.isRefresh = !this.isRefresh
      })
    },
    selectContainer() {
      this.form = this.transformYaml()
      this.currentContainerIndex = this.currentIndex
      this.isRefresh = !this.isRefresh
    },
    handleAddContainer() {
      this.container_lists.push({ label: "Container-" + this.container_lists.length.toString(), value: this.container_lists.length })
      this.form.spec.template.spec.containers.push({ name: "containers1", image: "", imagePullPolicy: "ifNotPresent" })
    },
    handleDeleteContainer() {
      this.currentContainerName = 0
      for (let i = 0; i < this.container_lists.length; i++) {
        if (this.container_lists[i].value === this.container_lists.length - 1) {
          this.container_lists.splice(i, 1)
        }
      }
    },
    transformYaml() {
      // general
      this.$refs.ko_general.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // ports
      this.$refs.ko_ports.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // command
      this.$refs.ko_command.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // resource
      this.$refs.ko_resource.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // health_check
      this.$refs.ko_health_readiness_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_health_liveness_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      this.$refs.ko_health_startup_check.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // security context
      this.$refs.ko_security_context.transformation(this.form.spec.template.spec.containers[this.currentContainerIndex])
      // networking
      this.$refs.ko_networking.transformation(this.form.spec.template.spec)
      // node scheduling
      this.$refs.ko_pod_scheduling.transformation(this.form.spec.template.spec)
      // upgrade policy
      this.$refs.ko_upgrade_policy.transformation(this.form.spec)
      // labels
      this.$refs.ko_labels.transformation(this.form.spec.template.metadata)
      // annotations
      this.$refs.ko_annotations.transformation(this.form.spec.template.metadata)
      return this.form
    },
    onCancel() {
      this.$router.push({ name: "Deployments" })
    },
    onSubmit() {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
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
      this.yaml = this.transformYaml()
      this.showYaml = true
    },
    backToForm() {
      this.showYaml = false
    },
  },
  mounted() {
    this.currentIndex = 0
    this.currentContainerIndex = 0
    if (this.$route.params.name) {
      this.search()
    }
  },
}
</script>

<style scoped>
</style>