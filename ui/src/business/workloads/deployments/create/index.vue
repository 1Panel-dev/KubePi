<template>
  <layout-content header="Deployment - Create">
    <el-row :gutter="20">
      <el-col :span="8">
        <el-row>
          <el-col :span="10">
            <ko-form-item labelName="Namespace" clearable itemType="select" :selections="namespace_list" v-model="form.namespace" />
          </el-col>
          <el-col :span="14">
            <ko-form-item labelName="Name" clearable itemType="input" v-model="form.name" />
          </el-col>
        </el-row>
      </el-col>
      <el-col :span="8">
        <ko-form-item labelName="Description" placeholder="Any text you want that better describes this resource" clearable itemType="input" v-model="form.description" />
      </el-col>
      <el-col :span="8">
        <ko-form-item labelName="Replicas" clearable itemType="input" v-model="form.replicas" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 30px">
      <el-col :span="8">
        <ko-form-item labelName="Container Image" clearable itemType="input" v-model="form.description" />
      </el-col>
      <el-col :span="8">
        <ko-form-item labelName="Image Pull Policy" clearable itemType="input" v-model="form.replicas" />
      </el-col>
    </el-row>

    <el-tabs style="margin-top: 50px" v-model="activeName" @tab-click="handleClick">
      <el-tab-pane label="Ports" name="Ports">
        <ko-ports />
      </el-tab-pane>
      <el-tab-pane label="Command" name="Command">
        <ko-command />
      </el-tab-pane>
      <el-tab-pane label="Resources" name="Resources">
        <ko-resources />
      </el-tab-pane>
      <el-tab-pane label="Health Check" name="Health Check">
        <ko-health-check style="margin-top=30px" health_check_type="Readiness Check" health_check_helper="Containers will be removed from service endpoints when this check is failing. Recommended." />
        <ko-health-check style="margin-top=30px" health_check_type="Liveness Check" health_check_helper="Containers will be restarted when this check is failing. Not recommended for most uses." />
        <ko-health-check style="margin-top=30px" health_check_type="Startup Check" health_check_helper="Containers will wait until this check succeeds before attempting other health checks." />
      </el-tab-pane>
      <el-tab-pane label="Security Context" name="Security Context">
        <ko-security-context />
      </el-tab-pane>
      <el-tab-pane label="Networking" name="Networking">
        <ko-networking />
      </el-tab-pane>
      <el-tab-pane label="Node Scheduling" name="Node Scheduling">
        <ko-node-scheduling />
      </el-tab-pane>
      <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy">
        <ko-upgrade-policy />
      </el-tab-pane>
      <el-tab-pane label="Labels" name="Labels">
        <ko-labels />
      </el-tab-pane>
      <el-tab-pane label="Annotations" name="Annotations">
        <ko-annotations />
      </el-tab-pane>
    </el-tabs>
    <!-- <el-button style="float: center" @click="getinfo">Create</el-button> -->
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoFormItem from "@/components/ko-form-item/index"
import KoPorts from "../../../../components/ko-work-load/ko-ports.vue"
import KoCommand from "../../../../components/ko-work-load/ko-command.vue"
import KoResources from "../../../../components/ko-work-load/ko-resources.vue"
import KoHealthCheck from "../../../../components/ko-work-load/ko-health-check.vue"
import KoSecurityContext from "../../../../components/ko-work-load/ko-security-context.vue"
import KoNetworking from "../../../../components/ko-work-load/ko-networking.vue"
import KoNodeScheduling from "../../../../components/ko-work-load/ko-node-scheduling.vue"
import KoUpgradePolicy from "../../../../components/ko-work-load/ko-upgrade-policy.vue"
import KoLabels from "../../../../components/ko-work-load/ko-labels.vue"
import KoAnnotations from "../../../../components/ko-work-load/ko-annotations.vue"

export default {
  name: "DeploymentCreate",
  components: { LayoutContent, KoFormItem, KoPorts, KoCommand, KoResources, KoHealthCheck, KoSecurityContext, KoNetworking, KoNodeScheduling, KoUpgradePolicy, KoLabels, KoAnnotations },
  data() {
    return {
      dns_policy_list: [
        { label: "zhangsan", value: "zhangsan111" },
        { label: "lisi", value: "lisi111" },
      ],
      namespace_list: [
        { label: "kube-system", value: "kube-system" },
        { label: "kube-public", value: "kube-public" },
        { label: "kube-operator", value: "kube-operator" },
        { label: "default", value: "default" },
      ],
      activeName: "Ports",
      form: {
        namespace: "",
        name: "",
        description: "",
        replicas: "",
        network_mode: "",
        dns_policy: "",
      },
    }
  },
  methods: {
    getinfo() {
      console.log(this.network_mode)
      console.log(this.dns_policy)
    },
    handleClick() {},
  },
}
</script>

<style scoped>
</style>