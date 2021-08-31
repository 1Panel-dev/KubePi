<template>
  <div style="margin-top: 20px">
    <ko-card title="Session Affinity">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-radio-group v-model="specObj.sessionAffinity" @change="changeType(specObj.sessionAffinity)" :key="key">
            <div style="margin-top: 5px">
              <el-radio label="None">None</el-radio>
            </div>
            <div style="margin-top: 5px">
              <el-radio label="ClientIP">Client IP</el-radio>
            </div>
          </el-radio-group>
        </el-col>
        <div v-if="specObj.sessionAffinity === 'ClientIP'">
          <el-col :span="11">
            <div>Session Sticky Time</div>
            <el-input v-model.number="specObj.sessionAffinityConfig.clientIP.timeoutSeconds" @input="$forceUpdate()">
              <template slot="append">Seconds</template>
            </el-input>
          </el-col>
        </div>
      </el-row>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoServiceSessionAffinity",
  components: { KoCard },
  props: {
    specObj: Object
  },
  data () {
    return {
      key: 0
    }
  },
  methods: {
    changeType (value) {
      this.key ++
      if (value === "ClientIP") {
        this.specObj.sessionAffinityConfig.clientIP.timeoutSeconds = 10800
      } else {
        delete this.specObj.sessionAffinityConfig.clientIP.timeoutSeconds
      }
    }
  },
  created () {
    if (!this.specObj?.sessionAffinityConfig) {
      this.specObj.sessionAffinity = "None"
      this.specObj.sessionAffinityConfig = {
        clientIP: {}
      }
    }
  }
}
</script>

<style scoped>

</style>
