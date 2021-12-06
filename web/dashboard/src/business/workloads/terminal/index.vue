<template>
  <div style="background-color: #1f2224" v-title :data-title="$t('business.pod.controller') + terminal.namespace + '/' + terminal.pod + '/' + terminal.container">
    <el-row>
      <div class="terminalOption" v-if="terminal.type ==='terminal'">
        <el-radio-group size="mini" @change="changeConditions()" v-model="shell">
          <el-radio-button label="bash"></el-radio-button>
          <el-radio-button label="sh"></el-radio-button>
        </el-radio-group>
      </div>
      <div class="terminalOption">
        <span class="spanClass">{{$t('business.workload.container')}}</span>
        <el-select class="interval" @change="changeConditions()" size="mini" v-model="terminal.container">
          <el-option v-for="c in containers" :key="c" :label="c" :value="c" />
        </el-select>
      </div>
      <div v-if="terminal.type ==='log'" style="background-color: #000000; display:inline">
        <div class="terminalOption">
          <span class="spanClass">{{$t('business.pod.lines')}}</span>
          <el-select class="interval" @change="changeConditions()" size="mini" v-model="tailLines">
            <el-option v-for="l in tailLinesOptions" :key="l.label" :label="l.label" :value="l.value" />
          </el-select>
        </div>
        <div style="margin-top: 15px; margin-bottom: 10px; float: left">
          <span class="spanClass">{{$t('business.pod.watch')}}</span>
          <el-switch class="interval" @change="changeConditions()" v-model="follow" />
        </div>
        <div class="terminalOption">
          <el-button style="margin-left: 20px;" size="mini" @click="dialogDownloadVisible = true">{{$t('business.pod.download_logs')}}</el-button>
        </div>
      </div>
      <div style="float: right;margin-top: 15px;margin-bottom: 5px;margin-right: 30px">
        <span style="font-size: 20px; color: white">{{terminal.namespace}}/{{terminal.pod}}/{{terminal.container}}</span>
      </div>

    </el-row>
    <el-row>
      <div>
        <iframe :key="isRefresh" :src="terminal.url" :style="{'height': height}" style="width: 100%;border: 0"></iframe>
      </div>
    </el-row>

    <el-dialog :title="$t('business.pod.download_logs') + ': ' + terminal.namespace + '/' + terminal.pod + '/' + terminal.container" width="70%" :close-on-click-modal="false" :visible.sync="dialogDownloadVisible">
      <el-form label-position="top" style="margin-left: 40px">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="Namespace">
              <span>{{terminal.namespace}}</span>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Pod">
              <div class="spanInFormStyle"><span :title="terminal.pod">{{terminal.pod}}</span></div>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="Container">
              <div class="spanInFormStyle"><span :title="terminal.container">{{terminal.container}}</span></div>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item :label="$t('business.pod.start_time')">
              <el-date-picker style="width:100%" v-model="form.limitDate" type="datetime" :picker-options="pickerOptions" align="right" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item :label="$t('business.pod.limit_byte')">
              <ko-form-item itemType="number" v-model="form.limitBytes" deviderName="MB" />
            </el-form-item>
          </el-col>
        </el-row>
        <div><span> {{$t('business.pod.download_log_help', [datetimeFormat(form.limitDate), form.limitBytes])}}</span></div>
      </el-form>

      <div slot="footer" class="dialog-footer">
        <el-button size="small" @click="dialogDownloadVisible = false">{{ $t("commons.button.cancel") }}</el-button>
        <el-button size="small" @click="onSubmitDown">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getWorkLoadByName } from "@/api/workloads"
import { getPodLogsByName } from "@/api/pods"
import KoFormItem from "@/components/ko-form-item/index"
import FileSaver from "file-saver"
import { datetimeFormat } from "fit2cloud-ui/src/filters/time"
export default {
  name: "Terminal",
  components: { KoFormItem },
  data() {
    return {
      height: "",
      shell: "bash",
      isRefresh: false,
      follow: true,
      tailLines: 20,
      tailLinesOptions: [
        { label: this.$t("business.pod.last_20_lines"), value: 20 },
        { label: this.$t("business.pod.last_100_lines"), value: 100 },
        { label: this.$t("business.pod.last_200_lines"), value: 200 },
        { label: this.$t("business.pod.last_500_lines"), value: 500 },
      ],
      terminal: {
        cluster: "",
        namespace: "",
        pod: "",
        container: "",
        type: "",
        url: "",
      },
      containers: [],
      dialogDownloadVisible: false,
      form: {
        limitBytes: 1,
        limitDate: new Date(new Date().setTime(new Date().getTime() - 1800 * 1000)),
      },
      pickerOptions: {
        shortcuts: [
          {
            text: this.$t("business.pod.last_half_hour"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 1800 * 1000)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_three_hour"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 3)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_day"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_three_days"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 3)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_week"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_month"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
              picker.$emit("pick", start)
            },
          },
          {
            text: this.$t("business.pod.last_three_month"),
            onClick(picker) {
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
              picker.$emit("pick", start)
            },
          },
        ],
        disabledDate: (time) => {
          return time.getTime() > Date.now()
        },
      },
    }
  },
  methods: {
    datetimeFormat(val) {
      return datetimeFormat(val)
    },
    getHeight() {
      this.height = document.body.clientHeight - 51 + "px"
    },
    getTerminalUrl() {
      if (this.terminal.type == "terminal") {
        return `${process.env.VUE_APP_TERMINAL_PATH}/app?cluster=${this.terminal.cluster}&pod=${this.terminal.pod}&namespace=${this.terminal.namespace}&container=${this.terminal.container}&shell=${this.shell}`
      } else {
        return `${process.env.VUE_APP_TERMINAL_PATH}/logging?cluster=${this.terminal.cluster}&pod=${this.terminal.pod}&namespace=${this.terminal.namespace}&container=${this.terminal.container}&tailLines=${this.tailLines}&follow=${this.follow}`
      }
    },
    changeConditions() {
      this.terminal.url = this.getTerminalUrl()
      this.isRefresh = !this.isRefresh
    },
    loadContainters() {
      this.containers = []
      getWorkLoadByName(this.terminal.cluster, "pods", this.terminal.namespace, this.terminal.pod).then((res) => {
        for (const c of res.spec.containers) {
          this.containers.push(c.name)
        }
      })
    },
    onSubmitDown() {
      let params = {}
      params.container = this.terminal.container
      params.pretty = true
      params.timestamps = true
      if (this.form.limitDate) {
        if (Math.floor((new Date().getTime() - this.form.limitDate.getTime()) / 1000) > 0) {
          params.sinceSeconds = Math.floor((new Date().getTime() - this.form.limitDate.getTime()) / 1000)
        }
      }
      if (this.form.limitBytes) {
        params.limitBytes = this.form.limitBytes * 1024 * 1024
      }
      getPodLogsByName(this.terminal.cluster, this.terminal.namespace, this.terminal.pod, params).then((res) => {
        const blob = new Blob([res], { type: "text/json" })
        FileSaver.saveAs(blob, this.terminal.pod + "_" + this.terminal.container + ".log")
      })
    },
  },
  created() {
    this.terminal = {
      cluster: this.$route.query.cluster,
      namespace: this.$route.query.namespace,
      pod: this.$route.query.pod,
      container: this.$route.query.container,
      type: this.$route.query.type,
    }
    this.terminal.url = this.getTerminalUrl()
    this.loadContainters()

    this.getHeight()
  },
}
</script>

<style lang="scss" scoped>
.terminalOption {
  float: left;
  margin-left: 5px;
  margin-top: 10px;
  margin-bottom: 10px;
}
.spanClass {
  margin-left: 20px;
  color: white;
}
.interval {
  margin-left: 10px;
}
/deep/ .scrollbar {
  .el-scrollbar__thumb {
    background-color: black;
  }
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 3px;
  }
}
</style>