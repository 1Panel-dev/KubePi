<template>
  <div>
    <div style="position: relative;">
      <el-row>
        <el-tabs class="TabsWithoutContant" style="min-height: 35px;background-color: #1f2224" type="border-card" v-model="currentTab" @tab-click="handleClick" closable @tab-remove="removeTab">
          <el-tab-pane v-for="item in terminalTabs" :label="item.name" :key="item.key" :name="item.key">
            <span slot="label"><i class="iconfont iconline-terminalzhongduan" style="color: white;border: 4"></i>  {{item.name}}</span>
          </el-tab-pane>
        </el-tabs>
        <div class="terminalOption">
          <span class="spanClass">{{$t('business.workload.container')}}</span>
          <el-select class="interval" @change="changeConditions(currentTerminal)" size="mini" v-model="currentTerminal.container">
            <el-option v-for="c in currentTerminal.containers" :key="c" :label="c" :value="c" />
          </el-select>
        </div>
        <div v-if="currentTerminal.type ==='logs'" style="background-color: #000000; display:inline">
          <div class="terminalOption">
            <span class="spanClass">{{$t('business.pod.lines')}}</span>
            <el-select class="interval" @change="changeConditions(currentTerminal)" size="mini" v-model="tailLines">
              <el-option v-for="l in tailLinesOptions" :key="l.label" :label="l.label" :value="l.value" />
            </el-select>
          </div>
          <div style="margin-top: 15px; margin-bottom: 10px; float: left">
            <span class="spanClass">{{$t('business.pod.watch')}}</span>
            <el-switch class="interval" @change="changeConditions(currentTerminal)" v-model="follow" />
          </div>
        </div>
        <div class="tabButton">
          <el-button v-if="expand" size="mini" @click="handleResize('shrink')" icon="el-icon-arrow-down" circle />
          <el-button v-if="!expand" size="mini" @click="handleResize('expand')" icon="el-icon-arrow-up" circle />
        </div>
      </el-row>
      <el-row>
        <div v-for="item in terminalTabs" :key="item.key">
          <div class="example" :style="{height: terminalHeight}">
            <el-scrollbar style="height:100%">
              <iframe v-show="item.show" :src="item.url" style="width: 100%;height: 100%;min-height: 800px;border: 0"></iframe>
            </el-scrollbar>
          </div>
        </div>
      </el-row>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      currentTab: "",
      currentTerminal: {},
      expand: false,
      tabIndex: 1,

      follow: false,
      tailLines: 20,
      tailLinesOptions: [
        { label: this.$t("business.pod.last_20_lines"), value: 20 },
        { label: this.$t("business.pod.last_100_lines"), value: 100 },
        { label: this.$t("business.pod.last_200_lines"), value: 200 },
        { label: this.$t("business.pod.last_500_lines"), value: 500 },
      ],
    }
  },
  computed: {
    terminalTabs: {
      get() {
        const terminals = this.$store.state.terminal.terminals || []
        if (terminals.length == 0) {
          this.handleResize("close")
        } else {
          this.handleResize("expand")
        }
        return this.$store.state.terminal.terminals
      },
      set() {
        this.$store.commit("terminal/TERMINALS", this.terminalTabs)
      },
    },
    terminalHeight() {
      return this.$store.state.terminal.terminalHeight - 90 + "px"
    },
  },
  watch: {
    terminalTabs: function (val) {
      if (val && val.length > 0) {
        for (const item of val) {
          item.show = false
        }
        val[val.length - 1].show = true
        val[val.length - 1].url = this.getTerminalUrl(val[val.length - 1])
        this.currentTab = val[val.length - 1].key
        this.currentTerminal = val[val.length - 1]
      }
    },
  },
  methods: {
    getTerminalUrl(item) {
      if (item.type == "terminal") {
        return `/terminal/app?cluster=${item.cluster}&pod=${item.pod}&namespace=${item.namespace}&container=${item.container}`
      } else {
        return `/terminal/logging?cluster=${item.cluster}&pod=${item.pod}&namespace=${item.namespace}&container=${item.container}&tailLines=${this.tailLines}&follow=${this.follow}`
      }
    },
    changeConditions(item) {
      item.url = this.getTerminalUrl(item)
    },
    removeTab(targetName) {
      for (let i = 0; i < this.terminalTabs.length; i++) {
        if (this.terminalTabs[i].key === targetName) {
          this.terminalTabs.splice(i, 1)
        }
      }
      if (this.terminalTabs.length === 0) {
        this.handleResize("close")
      } else {
        this.currentTab = this.terminalTabs[this.terminalTabs.length - 1].key
      }
    },
    handleResize(operation) {
      switch (operation) {
        case "expand":
          this.$store.commit("terminal/CHANGE_BOTTOM_HEIGHT", "400")
          this.$store.commit("terminal/CHANGE_TERMINAL_HEIGHT", "400")
          this.expand = true
          break
        case "shrink":
          this.$store.commit("terminal/CHANGE_BOTTOM_HEIGHT", "40")
          this.expand = false
          break
        case "close":
          this.$store.commit("terminal/CHANGE_BOTTOM_HEIGHT", "0")
          this.expand = false
          break
      }
    },
    handleClick() {
      this.handleResize("expand")
      for (const tab of this.terminalTabs) {
        if (tab.key == this.currentTab) {
          this.currentTerminal = tab
          tab.show = true
        } else {
          tab.show = false
        }
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.tabButton {
  position: absolute;
  right: 20px;
  top: 5px;
}
.terminalOption {
  float: left;
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
.example {
  ul {
    height: 20px;
  }
  /deep/ .el-scrollbar__wrap {
    height: 100%;
    overflow-x: hidden;
  }
}
</style>