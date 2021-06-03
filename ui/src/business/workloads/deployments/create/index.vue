<template>
  <layout-content header="Deployment - Create">
    <el-row :gutter="40">
      <el-col :span="8">
        <el-row>
          <el-col :span="12">
            <ko-form-item labelName="Namespace" clearable itemType="select" :selections="namespace_list" v-model="form.namespace" />
          </el-col>
          <el-col :span="12">
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
    <el-row :gutter="40" style="margin-top: 30px">
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
      <el-tab-pane label="Resources" name="Resources"></el-tab-pane>
      <el-tab-pane label="Health Check" name="Health Check"></el-tab-pane>
      <el-tab-pane label="Security Context" name="Security Context"></el-tab-pane>
      <el-tab-pane label="Networking" name="Networking"></el-tab-pane>
      <el-tab-pane label="Node Scheduling" name="Node Scheduling"></el-tab-pane>
      <el-tab-pane label="Scaling/Upgrade Policy" name="Scaling/Upgrade Policy"></el-tab-pane>
      <el-tab-pane label="Labels" name="Labels"></el-tab-pane>
      <el-tab-pane label="Annotations" name="Annotations"></el-tab-pane>
    </el-tabs>
    <!-- <el-button style="float: center" @click="getinfo">Create</el-button> -->
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import KoFormItem from "@/components/ko-form-item/index"
import KoPorts from '../../../../components/ko-work-load/ko-ports.vue'
import KoCommand from '../../../../components/ko-work-load/ko-command.vue'

export default {
  name: "DeploymentCreate",
  components: { LayoutContent, KoFormItem, KoPorts, KoCommand },
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
      activeName: "",
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