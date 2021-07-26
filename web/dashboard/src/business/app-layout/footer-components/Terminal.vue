<template>
  <div>
    <div style="position: relative;">
      <el-tabs style="min-height: 50px;backgroud: white" v-model="currentTab" @tab-click="handleClick" closable @tab-remove="removeTab">
        <el-tab-pane v-for="item in terminalTabs" :key="item.name" :label="item.title" :name="item.title">
          <el-card :style="`height: ${expand ? '330px' : 0 }`" style="width: 100%;box-sizing: border-box;" />
          <!-- <iframe :src="item.url"></iframe> -->
        </el-tab-pane>
      </el-tabs>
      <div class="tabButton">
        <el-button size="mini" @click="addTabs" icon="el-icon-plus" circle />
        <!-- <el-button v-if="terminalTabs.length !== 0" size="mini" @click="expandTabs" icon="el-icon-full-screen" circle /> -->
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
      })
      this.tabIndex++
      this.currentTab = newTabName
    },
    removeTab(targetName) {
      for (let i = 0; i < this.terminalTabs.length; i++) {
        if (this.terminalTabs[i].title === targetName) {
          this.terminalTabs.splice(i, 1)
        }
      }
      if (this.terminalTabs.length === 0) {
        this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "50")
        this.expand = false
      } else {
        this.currentTab = this.terminalTabs[this.terminalTabs.length - 1].title
      }
    },
    shrinkTabs() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "50")
      this.expand = false
      console.log(this.terminalTabs)
    },
    expandTabs() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "600")
      this.expand = true
    },
    handleClick() {
      this.$store.commit("app/CHANGE_BOTTOM_HEIGHT", "400")
      this.expand = true
    },
  },
}
</script>

<style scoped>
.tabButton {
  position: absolute;
  right: 20px;
  top: 5px;
  color: #e6a23c;
  font-weight: 600;
  font-size: 14px;
}
</style>