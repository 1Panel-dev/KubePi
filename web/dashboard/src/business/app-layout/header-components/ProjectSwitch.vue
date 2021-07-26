<template>
  <el-dropdown trigger="click" @command="handleProjectSwitch">
    <span class="el-dropdown-link">
              <i class="iconfont iconnamesapce" style="color: #3884c5;margin-right: 3px" :icon="['fas', 'globe']"/>
      <span>{{ getNamespaceName }}</span>
      <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item command="">全部</el-dropdown-item>
      <el-dropdown-item disabled divided></el-dropdown-item>
      <el-dropdown-item v-for="(value,key) in namespaceOptions" :key="key" :command="value">{{ value }}
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>

import {getNamespaces} from '@/api/auth'

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
      return (p === null) ? "全部" : p
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
    }
  },
  created() {
    const cluster = this.$store.getters.cluster
    getNamespaces(cluster).then((data) => {
      this.namespaceOptions = data.data
    })
  }
}

</script>

<style scoped>

</style>