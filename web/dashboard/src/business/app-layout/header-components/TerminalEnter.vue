<template>
  <div style="margin-right: 20px;margin-top: 5px">
    <el-button size="mini" @click="onOpenTerminal"><i class="el-icon-arrow-right"/>Terminal</el-button>
  </div>
</template>

<script>
import {getTerminalSession} from "@/api/webkubectl";

export default {
  name: "TerminalEnter",
  data() {
    return {}
  },
  methods: {
    onOpenTerminal() {
      const clusterName = this.$store.getters.cluster
      getTerminalSession(clusterName).then(data => {
        const token = data.data.token
        const url = `/webkubectl/root?token=${token}`
        console.log(url)
        open(url, "_black")
      })
    }
  }
}
</script>

<style scoped>

</style>