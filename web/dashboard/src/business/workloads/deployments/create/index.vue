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
        <ko-ports ref="ko_ports" />
      </el-tab-pane>
      <el-tab-pane label="Command" name="Command">
        <ko-command ref="ko_command" />
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
    <el-button style="float: center" @click="getinfo">Create</el-button>
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
import { deepClone } from "@/utils/deepClone"

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
        apiVersion: "apps/v1",
        kind: "Deployment",
        metadata: {
          name: "",
          labels: {
            app: "",
          },
        },
        spec: {
          replicas: 3,
          selector: "",
          template: {
            metadata: {
              labels: [],
            },
            spec: {
              containers: {
                name: "",
                image: "",
                ports: [],
                env: [],
                envFrom: [],
                command: [],
                args: [],
                resources: {},
                livenessProbe: {},
                readinessProbe: {},
                startupProbe: {},
                securityContext: {},
                affinity: {
                  nodeAffinity: {
                    requiredDuringSchedulingIgnoredDuringExecution: {
                      nodeSelectorTerms: [],
                    },
                    preferredDuringSchedulingIgnoredDuringExecution: {
                      preference: [],
                    },
                  },
                },
                hostname: "",
                subdomain: "",
                dnsPolicy: "",
                hostNetwork: "",
                hostAliases: [],
                dnsConfig: [],
              },
            },
          },
        },
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
      //ports
      if (this.$refs.ko_ports.ports) {
        this.form.spec.template.ports = this.$refs.ko_ports.ports
        for (const po of this.form.spec.template.ports) {
          po.expose = this.$refs.ko_ports.isExpose
        }
      } else {
        delete this.form.spec.template.ports
      }
      // command
      const commandForm = deepClone(this.$refs.ko_command.form)
      if (commandForm.args) {
        this.form.spec.template.spec.containers.args = commandForm.args.split(",")
      } else {
        delete this.form.spec.template.spec.containers.args
      }
      if (commandForm.workingDir) {
        this.form.spec.template.spec.containers.workingDir = commandForm.workingDir
      } else {
        delete this.form.spec.template.spec.containers.workingDir
      }
      if (commandForm.stdin) {
        switch (commandForm.stdin) {
          case "No":
            this.form.spec.template.spec.containers.stdin = false
            break
          case "Yes":
            this.form.spec.template.spec.containers.stdin = true
            break
          case "Ones":
            this.form.spec.template.spec.containers.stdin = true
            this.form.spec.template.spec.containers.stdinOnce = true
            break
        }
      } else {
        delete this.form.spec.template.spec.containers.stdin
      }
      if (commandForm.tty) {
        this.form.spec.template.spec.containers.tty = commandForm.tty
      } else {
        delete this.form.spec.template.spec.containers.tty
      }
      let envList = []
      let envFromList = []
      if (commandForm.env || commandForm.envFromResource) {
        for (const en of commandForm.env) {
          envList.push(en)
        }
        for (const en of commandForm.envFromResource) {
          switch (en.type) {
            case "Resource":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  resourceFieldRef: {
                    containerName: en.source,
                    divisor: 0,
                    resource: en.key,
                  },
                },
              })
              break
            case "ConfigMap":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  configMapKeyRef: {
                    name: en.source,
                    optional: false, // 这个false是什么意思不知道
                  },
                },
              })
              break
            case "Secret Key":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  secretKeyRef: {
                    name: en.source,
                    optional: false,
                  },
                },
              })
              break
            case "Field":
              envList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  fieldRef: {
                    apiVersion: "v1",
                    fieldPath: en.source,
                  },
                },
              })
              break
            case "Secret":
              envFromList.push({
                name: en.prefix_or_alias,
                valueFrom: {
                  secretRef: {
                    name: en.source,
                    optional: false,
                  },
                },
              })
              break
          }
        }
      }
      this.form.spec.template.spec.containers.env = envList
      this.form.spec.template.spec.containers.envFrom = envFromList
      console.log(this.form.spec.template.spec.containers)
    },
    handleClick() {},
  },
}
</script>

<style scoped>
</style>