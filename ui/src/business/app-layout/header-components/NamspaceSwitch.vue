<template>
  <el-dropdown trigger="click">
    <span class="el-dropdown-link">
              <i class="iconfont iconnamesapce" style="color: #3884c5;margin-right: 3px" :icon="['fas', 'globe']"/>
      <span>namespace</span>
      <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu slot="dropdown">
      <el-dropdown-item>
        <span>All Namespaces</span>
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
export default {
  name: "NamespaceSwitch"
}
</script>

<style scoped>

</style>