<template>
    <div>
      <el-select v-model="namespace" filterable @change="handleProjectSwitch">
        <el-option-group>
          <el-option label="All Namespaces" value=""></el-option>
        </el-option-group>
        <el-option-group>
          <el-option v-for="(value,key) in namespaceOptions" :key="key" :label="value" :value="value"></el-option>
        </el-option-group>
      </el-select>
    </div>
</template>

<script>

import {getNamespaces} from '@/api/auth'
import Bus from "@/utils/bus"

export default {
  name: "ProjectSwitch",
  data() {
    return {
      namespace: "",
      namespaceOptions: [],
    }
  },
  methods: {
    handleProjectSwitch() {
      if (!this.namespace) {
        sessionStorage.removeItem("namespace")
      } else {
        sessionStorage.setItem("namespace", this.namespace)
      }
      location.reload()
    },
    getNamespaceList() {
      const cluster = this.$store.getters.cluster
      getNamespaces(cluster).then((data) => {
        this.namespaceOptions = data.data
      })
    }
  },
  created() {
    this.getNamespaceList()
    let p = sessionStorage.getItem("namespace")
    this.namespace = p ? p : ""
  },
  mounted: function () {
    Bus.$on('refresh',() => {
      this.getNamespaceList()
    })
  }
}

</script>

<style scoped>

</style>
