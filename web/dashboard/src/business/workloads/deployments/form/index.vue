<template>
  <layout-content :header="'Deployment - ' + operation + ' - With - Form'">
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
        <ko-form-item labelName="Description" placeholder="Any text you want that better describes this resource" clearable itemType="input" v-model="description" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 30px">
      <el-col :span="12">
        <ko-form-item labelName="Container Image" clearable itemType="input" v-model="form.spec.template.spec.containers[0].image" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Image Pull Policy" clearable itemType="select" v-model="form.spec.template.spec.containers.imagePullPolicy" :selections="image_pull_policy_list" />
      </el-col>
    </el-row>

    <el-tabs style="margin-top: 50px" v-model="activeName">
      <el-tab-pane label="Ports" name="Ports">
        <ko-ports ref="ko_ports" :key="isRefresh" :portParentObj="form.spec.template.spec.containers[0]" />
      </el-tab-pane>
      <el-tab-pane label="Command" name="Command">
        <ko-command ref="ko_command" :key="isRefresh" :commandParentObj="form.spec.template.spec.containers[0]" />
      </el-tab-pane>
      <el-tab-pane label="Resources" name="Resources">
        <ko-resources ref="ko_resource" :key="isRefresh" :resourceParentObj="form.spec.template.spec.containers[0]" />
      </el-tab-pane>
      <el-tab-pane label="Health Check" :key="isRefresh" name="Health Check">
        <ko-health-check ref="ko_health_readiness_check" :healthCheckParentObj="form.spec.template.spec.containers[0]" style="margin-top=30px" health_check_type="Readiness Check" health_check_helper="Containers will be removed from service endpoints when this check is failing. Recommended." />
        <ko-health-check ref="ko_health_liveness_check" :healthCheckParentObj="form.spec.template.spec.containers[0]" style="margin-top=30px" health_check_type="Liveness Check" health_check_helper="Containers will be restarted when this check is failing. Not recommended for most uses." />
        <ko-health-check ref="ko_health_startup_check" :healthCheckParentObj="form.spec.template.spec.containers[0]" style="margin-top=30px" health_check_type="Startup Check" health_check_helper="Containers will wait until this check succeeds before attempting other health checks." />
      </el-tab-pane>
      <el-tab-pane label="Security Context" name="Security Context">
        <ko-security-context ref="ko_security_context" :key="isRefresh" :securityContextParentObj="form.spec.template.spec.containers[0]" />
      </el-tab-pane>
      <el-tab-pane label="Networking" name="Networking">
        <ko-networking ref="ko_networking" :key="isRefresh" :networkingParentObj="form.spec.template.spec" />
      </el-tab-pane>
      <el-tab-pane label="Node Scheduling" name="Node Scheduling">
        <ko-node-scheduling ref="ko_node_scheduling" :key="isRefresh" :nodeSchedulingParentObj="form.spec.template.spec" />
      </el-tab-pane>
      <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
        <ko-upgrade-policy ref="ko_upgrade_policy" :key="isRefresh" :upgradePolicyParentObj="form.spec" />
      </el-tab-pane>
      <el-tab-pane label="Labels" name="Labels">
        <ko-labels ref="ko_labels" :key="isRefresh" :labelParentObj="form.spec.template.metadata" />
      </el-tab-pane>
      <el-tab-pane label="Annotations" name="Annotations">
        <ko-annotations ref="ko_annotations" :key="isRefresh" :annotationsParentObj="form.spec.template.metadata" />
      </el-tab-pane>
    </el-tabs>
    <div class="grid-content bg-purple-light" v-if="showYaml">
      <yaml-editor :value="yaml"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="loading" @click="onSubmit()" type="primary">
          {{ $t("commons.button.create") }}
        </el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoFormItem from "@/components/ko-form-item/index"
import KoPorts from "@/components/ko-workloads/ko-ports.vue"
import KoCommand from "@/components/ko-workloads/ko-command.vue"
import KoResources from "@/components/ko-workloads/ko-resources.vue"
import KoHealthCheck from "@/components/ko-workloads/ko-health-check.vue"
import KoSecurityContext from "@/components/ko-workloads/ko-security-context.vue"
import KoNetworking from "@/components/ko-workloads/ko-networking.vue"
import KoNodeScheduling from "@/components/ko-workloads/ko-node-scheduling.vue"
import KoUpgradePolicy from "@/components/ko-workloads/ko-upgrade-policy.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"

import YamlEditor from "@/components/yaml-editor"
import { getDeploymentByName } from "@/api/workloads"

export default {
  name: "DeploymentForm",
  components: { LayoutContent, KoFormItem, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoNodeScheduling, KoUpgradePolicy, KoLabels, KoAnnotations, YamlEditor },
  data() {
    return {
      showYaml: false,
      yaml: "",
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
      activeName: "Ports",
      description: "",
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
            metadata: {
              labels: {
                "workload.user.cattle.io/workloadselector": "deployment-default-undefined",
              },
            },
            spec: {
              containers: [{
                name: "",
                image: "",
                imagePullPolicy: "",
              }],
              restartPolicy: "Always",
            },
          },
        },
        type: "apps.deployment",
      },
    }
  },
  methods: {
    search() {
      getDeploymentByName(this.$route.params.cluster, this.$route.params.namespace, this.$route.params.name)
        .then((res) => {
          this.form = res
          this.isRefresh = !this.isRefresh
        })
    },
    transformYaml() {
      // ports
      this.$refs.ko_ports.transformation(this.form.spec.template.spec.containers)
      // command
      this.$refs.ko_command.transformation(this.form.spec.template.spec.containers)
      // resource
      this.$refs.ko_resource.transformation(this.form.spec.template.spec.containers)
      // health_check
      this.$refs.ko_health_readiness_check.transformation(this.form.spec.template.spec.containers)
      this.$refs.ko_health_liveness_check.transformation(this.form.spec.template.spec.containers)
      this.$refs.ko_health_startup_check.transformation(this.form.spec.template.spec.containers)
      // security context
      this.$refs.ko_security_context.transformation(this.form.spec.template.spec.containers)
      // networking
      this.$refs.ko_networking.transformation(this.form.spec.template.spec)
      // node scheduling
      this.$refs.ko_node_scheduling.transformation(this.form.spec.template.spec)
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
      // const data = this.transformYaml()
      // this.loading = true
      // createNamespace("test1", data).then(() => {
      //   this.$message({
      //     type: "success",
      //     message: this.$t("commons.msg.create_success"),
      //   })
      //   this.$router.push({ name: "Namespaces" })
      // }).finally(() => {
      //   this.loading = false
      // })
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
    if (this.$route.params.name) {
      this.search()
      this.operation = "Edit"
    } else {
      this.operation = "Create"
    }
  },
}
</script>

<style scoped>
</style>