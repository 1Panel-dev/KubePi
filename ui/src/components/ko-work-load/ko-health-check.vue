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
          <ko-form-item labelName="Type" clearable itemType="select" v-model="check_type" :selections="type_list" />
        </el-row>
        <div v-if="check_type === 'httpGet' || check_type == 'httpsGet'">
          <el-row style="margin-bottom: 20px;">
            <ko-form-item labelName="Check Port" clearable itemType="input" v-model="form.port" />
          </el-row>
          <el-row>
            <ko-form-item labelName="Request Path" clearable itemType="input" v-model="form.path" />
          </el-row>
        </div>
        <el-row v-if="check_type === 'tcpSocket'">
          <ko-form-item labelName="Check Port" clearable itemType="input" v-model="form.port" />
        </el-row>
        <el-row v-if="check_type === 'exec'">
          <ko-form-item labelName="Command to run" placeholder="e.g. cat /tmp/health" clearable itemType="input" v-model="form.command" />
        </el-row>
      </el-col>
      <el-col :span="12" v-if="check_type !== 'None' && check_type !== ''">
        <el-row :gutter="10" style="margin-bottom: 20px;">
          <el-col :span="8">
            <ko-form-item labelName="Check Interval" placeholder="Default: 10" clearable deviderName="sec" itemType="input" v-model="form.periodSeconds" />
          </el-col>
          <el-col :span="8">
            <ko-form-item labelName="Initial Delay" placeholder="Default: 0" clearable deviderName="sec" itemType="input" v-model="form.initialDelaySeconds" />
          </el-col>
          <el-col :span="8">
            <ko-form-item labelName="Timeout" placeholder="Default: 3" clearable deviderName="sec" itemType="input" v-model="form.timeoutSeconds" />
          </el-col>
        </el-row>
        <el-row :gutter="10" style="margin-bottom: 20px;">
          <el-col :span="12">
            <ko-form-item labelName="Seccess Threshold" placeholder="Default: 1" clearable itemType="input" v-model="form.successThreshold" />
          </el-col>
          <el-col :span="12">
            <ko-form-item labelName="Failure Threshold" placeholder="Default: 3" clearable itemType="input" v-model="form.failureThreshold" />
          </el-col>
        </el-row>
        <el-row style="margin-bottom: 20px;">
          <el-button @click="handleAdd">Add</el-button>
          <el-table v-if="form.httpHeaders.length !== 0" :data="form.httpHeaders">
            <el-table-column min-width="40" label="Key">
              <template v-slot:default="{row}">
                <ko-form-item placeholder="e.g. foo" clearable itemType="input" v-model="row.name" />
              </template>
            </el-table-column>
            <el-table-column min-width="40" label="Value">
              <template v-slot:default="{row}">
                <ko-form-item placeholder="e.g. bar" clearable itemType="input" v-model="row.value" />
              </template>
            </el-table-column>
            <el-table-column width="120px">
              <template v-slot:default="{row}">
                <el-button type="text" style="font-size: 20px" @click="handleDelete(row)">REMOVE</el-button>
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
        port: "",
        path: "",
        periodSeconds: 10,
        initialDelaySeconds: 0,
        timeoutSeconds: 3,
        successThreshold: 1,
        failureThreshold: 3,
        httpHeaders: [],
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
        name: "",
        value: "",
      }
      this.form.httpHeaders.unshift(item)
    },
  },
}
</script>