<template>
  <div style="margin-top: 20px">
    <ko-card :title="health_check_type">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-row>
              <el-form-item label="Type">
                <ko-form-item itemType="select" v-model="check_type" :selections="type_list" />
              </el-form-item>
            </el-row>
            <div v-if="check_type === 'httpGet' || check_type == 'httpsGet'">
              <el-row>
                <el-form-item label="Check Port" prop="httpGet.port">
                  <ko-form-item itemType="number" v-model.number="form.httpGet.port" />
                </el-form-item>
              </el-row>
              <el-row>
                <el-form-item label="Request Path" prop="httpGet.path">
                  <ko-form-item itemType="input" v-model="form.httpGet.path" />
                </el-form-item>
              </el-row>
            </div>
            <el-row v-if="check_type === 'tcpSocket'">
              <el-form-item label="Check Port" prop="tcpSocket.port">
                <ko-form-item itemType="number" v-model.number="form.tcpSocket.port" />
              </el-form-item>
            </el-row>
            <el-row v-if="check_type === 'exec'">
              <el-form-item label="Command to run" prop="exec.command">
                <ko-form-item placeholder="e.g. cat /tmp/health" itemType="input" v-model="form.exec.command" />
              </el-form-item>
            </el-row>
          </el-col>
          <el-col :span="12" v-if="check_type !== 'None' && check_type !== ''">
            <el-row :gutter="10">
              <el-col :span="8">
                <el-form-item label="Check Interval" prop="periodSeconds">
                  <ko-form-item placeholder="Default: 10" deviderName="sec" itemType="number" v-model.number="form.periodSeconds" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="Initial Delay" prop="initialDelaySeconds">
                  <ko-form-item placeholder="Default: 0" deviderName="sec" itemType="number" v-model.number="form.initialDelaySeconds" />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="Timeout" prop="timeoutSeconds">
                  <ko-form-item placeholder="Default: 3" deviderName="sec" itemType="number" v-model.number="form.timeoutSeconds" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-form-item label="Seccess Threshold" prop="successThreshold">
                  <ko-form-item placeholder="Default: 1" itemType="number" v-model.number="form.successThreshold" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="Failure Threshold" prop="failureThreshold">
                  <ko-form-item placeholder="Default: 3" itemType="number" v-model.number="form.failureThreshold" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row>
              <div><label>Header</label></div>
              <table style="width: 98%" class="tab-table">
                <tr>
                  <th scope="col" width="45%" align="left"><label>key</label></th>
                  <th scope="col" width="45%" align="left"><label>value</label></th>
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
        { label: "None", value: "None" },
        { label: "HTTP request returns a successful status (200-399)", value: "httpGet" },
        { label: "HTTPS request returns a successful status", value: "httpsGet" },
        { label: "TCP connection opens successfully", value: "tcpSocket" },
        { label: "Command run inside the container exits with status 0", value: "exec" },
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
          command: "",
        },
        tcpSocket: {
          port: "",
        },
      },
    }
  },
  methods: {
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
      if (this.form.periodSeconds) {
        childForm.periodSeconds = this.form.periodSeconds
      }
      if (this.form.initialDelaySeconds) {
        childForm.initialDelaySeconds = this.form.initialDelaySeconds
      }
      if (this.form.timeoutSeconds) {
        childForm.timeoutSeconds = this.form.timeoutSeconds
      }
      if (this.form.successThreshold) {
        childForm.successThreshold = this.form.successThreshold
      }
      if (this.form.failureThreshold) {
        childForm.failureThreshold = this.form.failureThreshold
      }
      if (this.form.httpHeaders.length !== 0) {
        let obj = {}
        for (let i = 0; i < this.form.httpHeaders.length; i++) {
          if (this.form.httpHeaders[i].key !== "") {
            obj[this.form.httpHeaders[i].key] = this.form.httpHeaders[i].value
          }
        }
        childForm.httpHeaders = obj
      }
      switch (this.check_type) {
        case "httpGet":
          childForm.httpGet = {}
          childForm.httpGet.scheme = "HTTP"
          if (this.form.httpGet.path) {
            childForm.httpGet.path = this.form.httpGet.path
          }
          if (this.form.httpGet.port) {
            childForm.httpGet.port = this.form.httpGet.port
          }
          break
        case "httpsGet":
          childForm.httpGet = {}
          childForm.httpGet.scheme = "HTTPS"
          if (this.form.httpGet.path) {
            childForm.httpGet.path = this.form.httpGet.path
          }
          if (this.form.httpGet.port) {
            childForm.httpGet.port = this.form.httpGet.port
          }
          break
        case "tcpSocket":
          childForm.tcpSocket = {}
          if (this.form.tcpSocket.port) {
            childForm.tcpSocket.port = this.form.tcpSocket.port
          }
          break
        case "exec":
          childForm.exec = {}
          if (this.form.exec.command) {
            childForm.exec.command = this.form.exec.command.split(",")
          }
          break
        default:
          break
      }
      switch (this.health_check_type) {
        case "Readiness Check":
          parentFrom.readinessProbe = childForm
          break
        case "Liveness Check":
          parentFrom.livenessProbe = childForm
          break
        case "Startup Check":
          parentFrom.startupProbe = childForm
          break
      }
    },
  },
  mounted() {
    if (this.healthCheckParentObj) {
      let prodeForm = {}
      switch (this.health_check_type) {
        case "Readiness Check":
          prodeForm = this.healthCheckParentObj.readinessProbe
          break
        case "Liveness Check":
          prodeForm = this.healthCheckParentObj.livenessProbe
          break
        case "Startup Check":
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
            this.form.exec.command = prodeForm.exec.command.join(",")
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