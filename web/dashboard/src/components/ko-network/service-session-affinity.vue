<template>
  <div style="margin-top: 20px">
    <ko-card title="sessionAffinity">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-radio-group v-model="spec.sessionAffinity" @change="changeType">
            <div style="margin-top: 5px">
              <el-radio label="None">None</el-radio>
            </div>
            <div style="margin-top: 5px">
              <el-radio label="ClientIP">Client IP</el-radio>
            </div>
          </el-radio-group>
        </el-col>
        <div v-if="spec.sessionAffinity === 'ClientIP'">
          <el-col :span="11">
            <div>Session Sticky Time</div>
            <el-input v-model.number="spec.sessionAffinityConfig.clientIP.timeoutSeconds">
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
      spec: {
        sessionAffinity: "None",
        sessionAffinityConfig: {
          clientIP: {
            timeoutSeconds: 10800
          }
        }
      }
    }
  },
  methods: {
    changeType (value) {
      this.specObj.sessionAffinity = value
      if (value === "ClientIP") {
        this.spec.sessionAffinityConfig.clientIP.timeoutSeconds = 10800
      } else {
        delete this.spec.sessionAffinityConfig.clientIP.timeoutSeconds
      }
    }
  },
  created () {
    if (this.specObj) {
      this.spec.sessionAffinity = this.specObj.sessionAffinity ? this.specObj.sessionAffinity : this.specObj.sessionAffinity = "None"
      this.spec.sessionAffinityConfig = this.specObj.sessionAffinityConfig ? this.specObj.sessionAffinityConfig : this.specObj.sessionAffinityConfig = {
        clientIP: {
          timeoutSeconds: 10800
        }
      }
    }
  }
}
</script>

<style scoped>

</style>
