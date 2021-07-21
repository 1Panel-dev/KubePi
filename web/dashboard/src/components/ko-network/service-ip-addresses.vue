<template>
  <div style="margin-top: 20px">
    <ko-card title="IP Addresses">
      <el-form ref="form" label-position="top" :model="spec">
        <el-row v-if="spec.type === 'ClusterIP'">
          <el-col :span="12">
            <ko-form-item placeholder="e.g. foo" itemType="input" v-model="spec.clusterIP"
                          @change.native="transformation" @clear="transformation"/>
          </el-col>
        </el-row>
        <h2 v-if="spec.type === 'ClusterIP'">External IPs</h2>
        <div v-for="(row, index) in spec.externalIPs" v-bind:key="index" style="margin-top: 5px">
          <span v-if="false">{{ row }}</span>
          <el-row :gutter="20">
            <el-col :span="12">
              <ko-form-item placeholder="e.g. foo" itemType="input" v-model="spec.externalIPs[index]"/>
            </el-col>
            <el-col :span="2">
              <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </el-col>
          </el-row>
        </div>
        <div style="margin-top: 5px">
          <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
        </div>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoServiceIpAddresses",
  components: { KoCard, KoFormItem },
  props: {
    specObj: Object
  },
  data () {
    return {
      spec: {
        clusterIP: "",
        externalIPs: []
      }
    }
  },
  methods: {
    transformation () {
      this.specObj.clusterIP = this.spec.clusterIP
    },
    handleDelete (index) {
      this.spec.externalIPs.splice(index, 1)
    },
    handleAdd () {
      this.spec.externalIPs.push("")
    },
  },
  created () {
    if (this.specObj) {
      if (this.specObj.type === 'ClusterIP') {
        this.spec.clusterIP = this.specObj.clusterIP ? this.specObj.clusterIP : this.specObj.clusterIP = ""
      }
      this.spec.externalIPs = this.specObj.externalIPs ? this.specObj.externalIPs : this.specObj.externalIPs = []
    }
  }
}
</script>

<style scoped>

</style>
