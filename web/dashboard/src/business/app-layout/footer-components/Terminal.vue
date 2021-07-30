<template>
  <div>
    <div style="position: relative;">
      <el-tabs style="min-height: 35px;background-color: #1f2224" type="border-card" v-model="currentTab" @tab-click="handleClick" closable @tab-remove="removeTab">
        <el-tab-pane v-for="item in terminalTabs" :key="item.name" :name="item.title">
          <span slot="label"><i class="el-icon-date"></i>{{item.title}}</span>
          <!-- <iframe :src="item.url"></iframe> -->
        </el-tab-pane>
      </el-tabs>
      <div class="tabButton">
        <el-button size="mini" @click="addTabs" icon="el-icon-plus" circle />
        <el-button v-if="terminalTabs.length !== 0" size="mini" @click="shrinkTabs" icon="el-icon-arrow-down" circle />
      </div>
    </div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      currentTab: "",
      terminalTabs: [],
      expand: false,
      tabIndex: 1,
    }
  },
  methods: {
    addTabs() {
      let newTabName = "Terminal " + this.tabIndex
      this.terminalTabs.push({
        title: newTabName,
        name: this.tabIndex,
        url: "www.douyu.com",
        token: "",
      })
      this.tabIndex++
      this.currentTab = newTabName
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "400")
    },
    removeTab(targetName) {
      for (let i = 0; i < this.terminalTabs.length; i++) {
        if (this.terminalTabs[i].title === targetName) {
          this.terminalTabs.splice(i, 1)
        }
      }
      if (this.terminalTabs.length === 0) {
        this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "35")
        this.expand = false
      } else {
        this.currentTab = this.terminalTabs[this.terminalTabs.length - 1].title
      }
    },
    shrinkTabs() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "35")
      this.expand = false
    },
    handleClick() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "400")
      this.expand = true
    },
    // 这里将连接的 tablist 存到 localstorage 里面，刷新后，保证仍然存在
    beforeunloadHandler() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "35")
      this.expand = false
    },
  },
  mounted() {
    window.addEventListener("beforeunload", (e) => this.beforeunloadHandler(e))
  },
  destroyed() {
    window.removeEventListener("beforeunload", (e) => this.beforeunloadHandler(e))
  },
}
</script>

<style scoped>
.tabButton {
  position: absolute;
  right: 20px;
  top: 5px;
}
.divBorder {
  border-top: 8px solid #2f3236;
  border-bottom: 2px solid #2f3236;
}
</style>