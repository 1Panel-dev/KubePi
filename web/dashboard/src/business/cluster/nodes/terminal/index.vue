<template>
  <div style="background-color: #1f2224" >
    <el-row>
      <div style="float: right;margin-top: 15px;margin-bottom: 5px;margin-right: 30px">
        <span style="font-size: 20px; color: white">node path / is mounted to shell /host</span>
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
import KoFormItem from "@/components/ko-form-item/index"
import { isLogin } from "@/api/auth"
export default {
  name: "NodeTerminal",
  components: { KoFormItem },
  data() {
    return {
      height: "",
      isRefresh: false,
      terminal: {
        cluster: "",
        nodeName:""
      },
    }
  },
  methods: {
   
    getHeight() {
      this.height = document.body.clientHeight - 51 + "px"
    },
    getTerminalUrl() {
       return `${process.env.VUE_APP_TERMINAL_PATH}/node_shell?cluster=${this.terminal.cluster}&nodeName=${this.terminal.nodeName}`
    },
    pullingSession() {
      this.timer = setInterval(() => {
        isLogin().then(data => {
          this.items = data.data;
        }).catch(() => {
          this.$message.error(this.$t('commons.login.expires'))
          clearInterval(this.timer)
        })
      }, 5000)
    },
  },
  created() {
    this.terminal = {
      cluster: this.$route.query.cluster,
      nodeName: this.$route.query.nodeName
    }
    this.terminal.url = this.getTerminalUrl()
    this.getHeight()
    this.pullingSession()
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
::v-deep .scrollbar {
  .el-scrollbar__thumb {
    background-color: black;
  }
  .el-table__body-wrapper::-webkit-scrollbar {
    width: 3px;
  }
}
</style>
