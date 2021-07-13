<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Jobs' }" v-loading="loading">
    <br>
    <div v-if="!showYaml">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-form-item itemType="select2" disabled :selections="namespace_list" v-model="form.metadata.namespace" />
            </el-form-item>
          </el-col>
          <el-col :span="18">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <ko-form-item itemType="input" disabled v-model="form.metadata.name" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <el-tabs style="margin-top: 30px;background-color: #141418;" type="border-card" v-model="activeName">
        <el-tab-pane label="General" name="General">
          <div :key="isRefresh">
            <ko-container ref="ko_container" :isReadOnly="true" :containerParentObj="form.spec.template.spec.containers[currentContainerIndex]" :secretList="secret_list_of_ns" />
            <ko-ports ref="ko_ports" :isReadOnly="true" :portParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
            <ko-command ref="ko_command" :isReadOnly="true" :commandParentObj="form.spec.template.spec.containers[currentContainerIndex]" :currentNamespace="form.metadata.namespace" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Health Check" name="Health Check">
          <div :key="isRefresh">
            <ko-health-check ref="ko_health_readiness_check" :isReadOnly="true" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" health_check_type="Readiness Check" />
            <ko-health-check ref="ko_health_liveness_check" :isReadOnly="true" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" health_check_type="Liveness Check" />
            <ko-health-check ref="ko_health_startup_check" :isReadOnly="true" :healthCheckParentObj="form.spec.template.spec.containers[currentContainerIndex]" health_check_type="Startup Check" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Labels/Annotations" name="Labels/Annotations">
          <div :key="isRefresh">
            <ko-labels ref="ko_labels" :isReadOnly="true" :labelParentObj="form.metadata" labelTitle="Labels" />
            <ko-annotations ref="ko_annotations" :isReadOnly="true" :annotationsParentObj="form.metadata" annotationsTitle="Annotations" />
            <ko-labels ref="ko_labels" :isReadOnly="true" :labelParentObj="form.spec.template.metadata" labelTitle="Pod Labels" />
            <ko-annotations ref="ko_annotations" :isReadOnly="true" :annotationsParentObj="form.spec.template.metadata" annotationsTitle="Pod Annotations" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Networking" :isReadOnly="true" name="Networking">
          <ko-networking ref="ko_networking" :isReadOnly="true" :key="isRefresh" :networkingParentObj="form.spec.template.spec" />
        </el-tab-pane>
        <el-tab-pane label="Scheduling" name="Scheduling">
          <div :key="isRefresh">
            <ko-pod-scheduling :isReadOnly="true" ref="ko_pod_scheduling" :podSchedulingParentObj="form.spec.template.spec" :namespaceList="namespace_list" />
            <ko-node-scheduling :isReadOnly="true" ref="ko_node_scheduling" :nodeSchedulingParentObj="form.spec.template.spec" :nodeList="node_list" />
            <ko-tolerations :isReadOnly="true" ref="ko_toleration" :tolerationsParentObj="form.spec.template.spec" />
          </div>
        </el-tab-pane>
        <el-tab-pane label="Resources" name="Resources">
          <ko-resources :isReadOnly="true" ref="ko_resource" :key="isRefresh" :resourceParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
          <ko-upgrade-policy-job :isReadOnly="true" ref="ko_upgrade_policy_cronjob" :key="isRefresh" :upgradePolicyParentObj="form.spec" />
        </el-tab-pane>
        <el-tab-pane label="Security Context" name="Security Context">
          <ko-security-context :isReadOnly="true" ref="ko_security_context" :key="isRefresh" :securityContextParentObj="form.spec.template.spec.containers[currentContainerIndex]" />
        </el-tab-pane>
        <el-tab-pane label="Storage" name="Storage">
          <ko-storage :isReadOnly="true" ref="ko_storage" :key="isRefresh" :storageParentObj="form.spec.template.spec" :currentContainerIndex="currentContainerIndex" :configMapList="config_map_list_of_ns" :secretList="secret_list_of_ns" />
        </el-tab-pane>
      </el-tabs>
    </div>
    <div class="grid-content bg-purple-light" v-if="showYaml">
      <yaml-editor :read-only="true" :value="yaml"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.view_yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
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
import KoUpgradePolicyJob from "@/components/ko-workloads/ko-upgrade-policy-job.vue"
import KoLabels from "@/components/ko-workloads/ko-labels.vue"
import KoAnnotations from "@/components/ko-workloads/ko-annotations.vue"
import KoStorage from "@/components/ko-workloads/ko-storage.vue"

import YamlEditor from "@/components/yaml-editor"

import { getJobByName } from "@/api/jobs"

export default {
  name: "JobEdit",
  components: { LayoutContent, KoFormItem, YamlEditor, KoContainer, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoPodScheduling, KoNodeScheduling, KoTolerations, KoUpgradePolicyJob, KoLabels, KoAnnotations, KoStorage },
  data() {
    return {
      showYaml: false,
      yaml: {},
      isRefresh: false,
      loading: false,
      activeName: "General",
      namespace_list: [],
      secret_list: [],
      secret_list_of_ns: [],
      config_map_list: [],
      config_map_list_of_ns: [],
      node_list: [],
      currentContainerIndex: 0,
      form: {
        apiVersion: "apps/v1",
        kind: "Job",
        metadata: {
          name: "my-test-job",
          namespace: "",
        },
        spec: {
          template: {
            metadata: {},
            spec: {
              containers: [],
              restartPolicy: "Always",
            },
          },
        },
        type: "apps.job",
      },
      clusterName: "",
    }
  },
  methods: {
    search() {
      getJobByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        this.yaml = res
        this.isRefresh = !this.isRefresh
      })
    },
    onCancel() {
      this.$router.push({ name: "Jobs" })
    },
    onEditYaml() {
      this.showYaml = true
    },
    backToForm() {
      this.showYaml = false
    },
  },
  mounted() {
    this.showYaml = this.$route.query.yamlShow === "true"
    this.clusterName = this.$route.query.cluster
    if (this.$route.params.name) {
      this.search()
    }
  },
}
</script>

<style scoped>
</style>