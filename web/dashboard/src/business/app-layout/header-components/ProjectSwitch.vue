<template>
  <el-dropdown trigger="click" @command="handleProjectSwitch">
    <span class="el-dropdown-link">
              <i class="iconfont iconnamesapce" style="color: #3884c5;margin-right: 3px" :icon="['fas', 'globe']"/>
      <span>{{ getProjectName }}</span>
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
    getProjectName() {
      const p = sessionStorage.getItem("project")
      return (p === null) ? "Global" : p
    }
  },
  methods: {
    handleProjectSwitch(command) {
      if (!command) {
        sessionStorage.removeItem("project")
      } else {
        sessionStorage.setItem("project", command)
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