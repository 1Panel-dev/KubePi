<template>
  <el-dropdown trigger="click" @command="handleProjectSwitch">
    <span class="el-dropdown-link">
              <i class="iconfont iconnamesapce" style="color: #3884c5;margin-right: 3px" :icon="['fas', 'globe']"/>
      <span>{{ getNamespaceName }}</span>
      <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown" style="max-height: 300px;overflow:auto">
      <el-dropdown-item command="">All Namespaces</el-dropdown-item>
      <el-dropdown-item disabled divided></el-dropdown-item>
      <el-dropdown-item v-for="(value,key) in namespaceOptions" :key="key" :command="value">{{ value }}
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>

import {getNamespaces} from '@/api/auth'
import Bus from "@/utils/bus"

export default {
  name: "ProjectSwitch",
  data() {
    return {
      namespaceOptions: [],
    }
  },
  computed: {
    getNamespaceName() {
      const p = sessionStorage.getItem("namespace")
      return (p === null) ? "All Namespaces" : p
    }
  },
  methods: {
    handleProjectSwitch(command) {
      if (!command) {
        sessionStorage.removeItem("namespace")
      } else {
        sessionStorage.setItem("namespace", command)
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
