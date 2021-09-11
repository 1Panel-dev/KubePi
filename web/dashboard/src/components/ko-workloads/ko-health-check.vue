<template>
  <div style="margin-top: 20px">
    <ko-card :title="health_check_type">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-row>
              <el-form-item :label="$t('business.workload.check_type')">
                <ko-form-item itemType="select" :noClear="true" v-model="check_type" :selections="type_list" />
              </el-form-item>
            </el-row>
            <div v-if="check_type === 'httpGet' || check_type == 'httpsGet'">
              <el-row>
                <el-form-item :label="$t('business.workload.check_port')" prop="httpGet.port">
                  <ko-form-item itemType="number" v-model.number="form.httpGet.port" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-form-item :label="$t('business.workload.check_path')" prop="httpGet.path">
                  <ko-form-item itemType="input" v-model="form.httpGet.path" />
                </el-form-item>
              </el-row>
            </div>
            <el-row v-if="check_type === 'tcpSocket'">
              <el-form-item :label="$t('business.workload.check_port')" prop="tcpSocket.port">
                <ko-form-item itemType="number" v-model.number="form.tcpSocket.port" />
              </el-form-item>
            </el-row>
            <el-row v-if="check_type === 'exec'">
              <table style="width: 98%" class="tab-table">
                <tr>
                  <th scope="col" width="93%" align="left">
                    <label>{{$t('business.workload.check_cmd')}}</label>
                  </th>
                  <th align="left"></th>
                </tr>
                <tr v-for="(row,index) in form.exec.command" :key="index">
                  <td>
                    <ko-form-item placeholder="e.g. /tmp/health" itemType="textarea" v-model="row.value" />
                  </td>
                  <td>
                    <el-button type="text" style="font-size: 10px" @click="handleCommandDelete(index)">
                      {{ $t("commons.button.delete") }}
                    </el-button>
                  </td>
                </tr>
                <tr>
                  <td align="left">
                    <el-button @click="handleCommandAdd">{{$t('business.workload.add')}}</el-button>
                  </td>
                </tr>
              </table>
            </el-row>
          </el-col>
          <el-col :span="12" v-if="check_type !== 'None' && check_type !== ''">
            <el-row :gutter="10">
              <el-col :span="8">
                <el-form-item :label="$t('business.workload.check_interval')" prop="periodSeconds">
                  <ko-form-item placeholder="Default: 10" deviderName="sec" itemType="number" v-model.number="form.periodSeconds" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('business.workload.initial_delay')" prop="initialDelaySeconds">
                  <ko-form-item placeholder="Default: 0" deviderName="sec" itemType="number" v-model.number="form.initialDelaySeconds" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item :label="$t('business.workload.timeout')" prop="timeoutSeconds">
                  <ko-form-item placeholder="Default: 3" deviderName="sec" itemType="number" v-model.number="form.timeoutSeconds" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-form-item :label="$t('business.workload.seccess_threshold')" prop="successThreshold">
                  <ko-form-item placeholder="Default: 1" itemType="number" v-model.number="form.successThreshold" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item :label="$t('business.workload.failure_threshold')" prop="failureThreshold">
                  <ko-form-item placeholder="Default: 3" itemType="number" v-model.number="form.failureThreshold" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <div><label>{{$t('business.workload.header')}}</label></div>
              <table style="width: 98%" class="tab-table">
                <tr>
                  <th scope="col" width="45%" align="left"><label>{{$t('business.workload.key')}}</label></th>
                  <th scope="col" width="45%" align="left"><label>{{$t('business.workload.value')}}</label></th>
                  <th align="left"></th>
                </tr>
                <tr v-for="(row, index) in form.httpHeaders" v-bind:key="index">
                  <td>
                    <ko-form-item placeholder="e.g. foo" itemType="input" v-model="row.key" />
                  </td>
                  <td>
                    <ko-form-item placeholder="e.g. bar" itemType="input" v-model="row.value" />
                  </td>
                  <td>
                    <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
                      {{ $t("commons.button.delete") }}
                    </el-button>
                  </td>
                </tr>
                <tr>
                  <td align="left">
                    <el-button @click="handleAdd">Add</el-button>
                  </td>
                </tr>
              </table>
            </el-row>
          </el-col>
        </el-row>
      </el-form>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoHealthCheck",
  props: {
    health_check_type: String,
    health_check_helper: String,
    healthCheckParentObj: Object,
    isReadOnly: Boolean,
  },
  components: { KoFormItem, KoCard },
  data() {
    return {
      type_list: [
        { label: this.$t("business.workload.none"), value: "None" },
        { label: this.$t("business.workload.http_option"), value: "httpGet" },
        { label: this.$t("business.workload.https_option"), value: "httpsGet" },
        { label: this.$t("business.workload.tcp_option"), value: "tcpSocket" },
        { label: this.$t("business.workload.cmd_option"), value: "exec" },
      ],
      check_type: "None",
      form: {
        httpGet: {
          port: "",
          path: "",
        },
        periodSeconds: 10,
        initialDelaySeconds: 0,
        timeoutSeconds: 3,
        successThreshold: 1,
        failureThreshold: 3,
        httpHeaders: [],
        exec: {
          command: [],
        },
        tcpSocket: {
          port: "",
        },
      },
    }
  },
  methods: {
    handleCommandAdd() {
      var item = {
        value: "",
      }
      this.form.exec.command.push(item)
    },
    handleCommandDelete(index) {
      this.form.exec.command.splice(index, 1)
    },
    handleDelete(index) {
      this.form.httpHeaders.splice(index, 1)
    },
    handleAdd() {
      var item = {
        key: "",
        value: "",
      }
      this.form.httpHeaders.push(item)
    },
    transformation(parentFrom) {
      if (this.check_type === "None" || this.check_type === "") {
        return
      }
      let childForm = {}
      childForm.periodSeconds = this.form.periodSeconds || undefined
      childForm.initialDelaySeconds = this.form.initialDelaySeconds || undefined
      childForm.timeoutSeconds = this.form.timeoutSeconds || undefined
      childForm.successThreshold = this.form.successThreshold || undefined
      childForm.failureThreshold = this.form.failureThreshold || undefined

      let obj = {}
      for (let i = 0; i < this.form.httpHeaders.length; i++) {
        if (this.form.httpHeaders[i].key !== "") {
          obj[this.form.httpHeaders[i].key] = this.form.httpHeaders[i].value
        }
      }
      childForm.httpHeaders = this.form.httpHeaders.length !== 0 ? obj : undefined
      let commands = []
      switch (this.check_type) {
        case "httpGet":
          childForm.httpGet = {}
          childForm.httpGet.scheme = "HTTP"
          childForm.httpGet.path = this.form.httpGet.path || undefined
          childForm.httpGet.port = this.form.httpGet.port || undefined
          break
        case "httpsGet":
          childForm.httpGet = {}
          childForm.httpGet.scheme = "HTTPS"
          childForm.httpGet.path = this.form.httpGet.path || undefined
          childForm.httpGet.port = this.form.httpGet.port || undefined
          break
        case "tcpSocket":
          childForm.tcpSocket = {}
          childForm.tcpSocket.port = this.form.tcpSocket.port || undefined
          break
        case "exec":
          for (const cmd of this.form.exec.command) {
            commands.push(cmd.value)
          }
          childForm.exec = { command: commands.length !== 0 ? commands : undefined }
          break
        default:
          break
      }
      switch (this.health_check_type) {
        case this.$t("business.workload.readiness_check"):
          parentFrom.readinessProbe = childForm
          break
        case this.$t("business.workload.liveness_check"):
          parentFrom.livenessProbe = childForm
          break
        case this.$t("business.workload.startup_check"):
          parentFrom.startupProbe = childForm
          break
      }
    },
  },
  mounted() {
    if (this.healthCheckParentObj) {
      let prodeForm = {}
      switch (this.health_check_type) {
        case this.$t("business.workload.readiness_check"):
          prodeForm = this.healthCheckParentObj.readinessProbe
          break
        case this.$t("business.workload.liveness_check"):
          prodeForm = this.healthCheckParentObj.livenessProbe
          break
        case this.$t("business.workload.startup_check"):
          prodeForm = this.healthCheckParentObj.startupProbe
          break
      }
      if (prodeForm) {
        if (prodeForm.httpGet) {
          if (prodeForm.httpGet.scheme) {
            if (prodeForm.httpGet.scheme == "HTTP") {
              this.check_type = "httpGet"
            } else {
              this.check_type = "httpsGet"
            }
          } else {
            this.check_type = "httpGet"
          }
          if (prodeForm.httpGet.port) {
            this.form.httpGet.port = prodeForm.httpGet.port
          }
          if (prodeForm.httpGet.path) {
            this.form.httpGet.path = prodeForm.httpGet.path
          }
        } else if (prodeForm.tcpSocket) {
          this.check_type = "tcpSocket"
          if (prodeForm.tcpSocket.port) {
            this.form.tcpSocket.port = prodeForm.tcpSocket.port
          }
        } else if (prodeForm.exec) {
          this.check_type = "exec"
          if (prodeForm.exec.command) {
            this.form.exec.command = []
            for (const cmd of prodeForm.exec.command) {
              this.form.exec.command.push({ value: cmd })
            }
          }
        }

        if (prodeForm.initialDelaySeconds) {
          this.form.initialDelaySeconds = prodeForm.initialDelaySeconds
        }
        if (prodeForm.timeoutSeconds) {
          this.form.timeoutSeconds = prodeForm.timeoutSeconds
        }
        if (prodeForm.periodSeconds) {
          this.form.periodSeconds = prodeForm.periodSeconds
        }
        if (prodeForm.successThreshold) {
          this.form.successThreshold = prodeForm.successThreshold
        }
        if (prodeForm.failureThreshold) {
          this.form.failureThreshold = prodeForm.failureThreshold
        }
      }
    }
  },
}
</script>