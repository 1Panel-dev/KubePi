<template>
  <div style="background-color: #1f2224" v-title :data-title="'控制台  ' + terminal.namespace + '/' + terminal.pod + '/' + terminal.container">
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
  </div>
</template>

<script>
import { getWorkLoadByName } from "@/api/workloads"
export default {
  name: "Terminal",
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
    }
  },
  methods: {
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