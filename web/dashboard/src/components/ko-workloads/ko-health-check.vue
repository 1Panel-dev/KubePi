<template>
  <div style="margin-top: 20px">
    <div>
      <label>{{health_check_type}}
        <el-tooltip class="item" effect="dark" :content="health_check_helper" placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>
      </label>
    </div>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-row style="margin-bottom: 20px;">
          <ko-form-item labelName="Type" itemType="select" v-model="check_type" :selections="type_list" />
        </el-row>
        <div v-if="check_type === 'httpGet' || check_type == 'httpsGet'">
          <el-row style="margin-bottom: 20px;">
            <ko-form-item labelName="Check Port" itemType="number" v-model.number="form.httpGet.port" />
          </el-row>
          <el-row>
            <ko-form-item labelName="Request Path" itemType="input" v-model="form.httpGet.path" />
          </el-row>
        </div>
        <el-row v-if="check_type === 'tcpSocket'">
          <ko-form-item labelName="Check Port" itemType="number" v-model.number="form.tcpSocket.port" />
        </el-row>
        <el-row v-if="check_type === 'exec'">
          <ko-form-item labelName="Command to run" placeholder="e.g. cat /tmp/health" itemType="input" v-model="form.exec.command" />
        </el-row>
      </el-col>
      <el-col :span="12" v-if="check_type !== 'None' && check_type !== ''">
        <el-row :gutter="10" style="margin-bottom: 20px;">
          <el-col :span="8">
            <ko-form-item labelName="Check Interval" placeholder="Default: 10" deviderName="sec" itemType="input" v-model="form.periodSeconds" />
          </el-col>
          <el-col :span="8">
            <ko-form-item labelName="Initial Delay" placeholder="Default: 0" deviderName="sec" itemType="input" v-model="form.initialDelaySeconds" />
          </el-col>
          <el-col :span="8">
            <ko-form-item labelName="Timeout" placeholder="Default: 3" deviderName="sec" itemType="input" v-model="form.timeoutSeconds" />
          </el-col>
        </el-row>
        <el-row :gutter="10" style="margin-bottom: 20px;">
          <el-col :span="12">
            <ko-form-item labelName="Seccess Threshold" placeholder="Default: 1" itemType="input" v-model="form.successThreshold" />
          </el-col>
          <el-col :span="12">
            <ko-form-item labelName="Failure Threshold" placeholder="Default: 3" itemType="input" v-model="form.failureThreshold" />
          </el-col>
        </el-row>
        <el-row style="margin-bottom: 20px;">
          <el-button @click="handleAdd">Add</el-button>
          <el-table v-if="form.httpHeaders.length !== 0" :data="form.httpHeaders">
            <el-table-column min-width="40" label="Key">
              <template v-slot:default="{row}">
                <ko-form-item :withoutLabel="true" placeholder="e.g. foo" itemType="input" v-model="row.key" />
              </template>
            </el-table-column>
            <el-table-column min-width="40" label="Value">
              <template v-slot:default="{row}">
                <ko-form-item :withoutLabel="true" placeholder="e.g. bar" itemType="input" v-model="row.value" />
              </template>
            </el-table-column>
            <el-table-column  width="60px">
              <template v-slot:default="{row}">
                <el-button type="text" style="font-size: 10px" @click="handleDelete(row)">{{ $t("commons.button.delete") }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-row>
      </el-col>
    </el-row>
    <el-divider></el-divider>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoHealthCheck",
  props: {
    health_check_type: String,
    health_check_helper: String,
    healthCheckParentObj: Object,
  },
  components: { KoFormItem },
  data() {
    return {
      type_list: [
        { label: "None", value: "None" },
        { label: "HTTP request returns a successful status (200-399)", value: "httpGet" },
        { label: "HTTPS request returns a successful status", value: "httpsGet" },
        { label: "TCP connection opens successfully", value: "tcpSocket" },
        { label: "Command run inside the container exits with status 0", value: "exec" },
      ],
      check_type: "",
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
    handleDelete(row) {
      for (let i = 0; i < this.form.httpHeaders.length; i++) {
        if (this.form.httpHeaders[i] === row) {
          this.form.httpHeaders.splice(i, 1)
        }
      }
    },
    handleAdd() {
      var item = {
        key: "",
        value: "",
      }
      this.form.httpHeaders.unshift(item)
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
      if (this.form.httpHeaders) {
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
            childForm.exec.command = this.form.exec.command
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
        }
        if (prodeForm.tcpSocket) {
          this.check_type = "tcpSocket"
          if (prodeForm.tcpSocket.port) {
            this.form.tcpSocket.port = prodeForm.tcpSocket.port
          }
        }
        if (prodeForm.exec) {
          this.check_type = "exec"
          if (prodeForm.exec.command) {
            this.form.exec.command = prodeForm.exec.command
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