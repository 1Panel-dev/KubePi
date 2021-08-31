<template>
  <div style="margin-top: 20px">
    <ko-card :title="title">
      <div v-if="type === 'ClusterIP' || type === 'LoadBalancer' || type === 'NodePort'">
        <el-row>
          <el-col :span="12">
            <el-input placeholder="e.g. 10.1.1.1" v-model="specObj.clusterIP" @input="$forceUpdate()"></el-input>
          </el-col>
        </el-row>
        <h2>{{ $t("business.network.external_ip") }}</h2>
      </div>
      <div :key="key">
        <div v-for="(row, index) in specObj.externalIPs" style="margin-top: 5px" v-bind:key="index">
          <span v-if="false">{{ row }}</span>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-input placeholder="e.g. 10.1.1.1" v-model="specObj.externalIPs[index]"
                        @input="$forceUpdate()"></el-input>
            </el-col>
            <el-col :span="2">
              <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </el-col>
          </el-row>
        </div>

      </div>
      <div style="margin-top: 5px">
        <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
      </div>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoServiceIpAddresses",
  components: { KoCard },
  props: {
    specObj: Object,
    type: String
  },
  data () {
    return {
      title: this.$t("business.network.ip_address"),
      key: 0
    }
  },
  methods: {
    transformation () {
    },
    handleDelete (index) {
      this.specObj.externalIPs.splice(index, 1)
      this.key = Math.random()
    },
    handleAdd () {
      this.specObj.externalIPs.push("")
      this.key = Math.random()
    },
  },
  created () {
    if (this.specObj) {
      if (!this.specObj.externalIPs) {
        this.specObj.externalIPs = []
      }
      if (!this.specObj.clusterIP) {
        this.specObj.clusterIP = ""
      }
    }
  },
  watch: {
    type: function (newValue) {
      if (newValue === "ExternalName" || newValue === "Headless") {
        this.title = this.$t("business.network.external_ip")
      } else {
        this.title = this.$t("business.network.ip_address")
      }
    }
  }
}
</script>

<style scoped>

</style>
