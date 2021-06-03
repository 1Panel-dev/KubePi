<template>
  <el-dropdown  trigger="click">
    <span class="el-dropdown-link">
        <i class="el-icon-user-solid" style="margin-right: 3px"></i>
        <span>用户</span>
        <i class="el-icon-arrow-down el-icon--right"></i>
    </span>
    <el-dropdown-menu  slot="dropdown">
      <el-dropdown-item>
        <span>基本信息</span>
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
export default {
  name: "UserSetting"
}
</script>

<style scoped>

</style>