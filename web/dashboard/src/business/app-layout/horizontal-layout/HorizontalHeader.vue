<template>
  <div class="horizontal-header">
    <div class="header-left">
      <sidebar-toggle-button/>
      <home></home>
    </div>
    <div class="header-right">
      <language-switch></language-switch>
      <project-switch></project-switch>
      <user-setting style="margin-left: 20px"></user-setting>
    </div>
  </div>
</template>

<script>
import SidebarToggleButton from "@/components/layout/sidebar/SidebarToggleButton"
import UserSetting from "@/business/app-layout/header-components/UserSetting"
import Home from "@/business/app-layout/header-components/Home"
import ProjectSwitch from "@/business/app-layout/header-components/ProjectSwitch";
import LanguageSwitch from "@/business/app-layout/header-components/LanguageSwitch"

export default {
  name: "HorizontalHeader",
  components: { LanguageSwitch, ProjectSwitch, Home, UserSetting, SidebarToggleButton}
}
</script>

<style lang="scss" scoped>
@import "~@/styles/common/mixins";

.horizontal-header {
  @include flex-row(flex-start, center);
  position: relative;
  height: 100%;

  .header-left {
    @include flex-row(flex-start, center);
    position: relative;
    height: 100%;
  }

  .header-right {
    @include flex-row(flex-end, center);
    flex: auto;
    height: 100%;

    .navbar-item {
      color: #2E2E2E;
      line-height: 50px;
      display: inline-block;
      padding-right: 20px;
    }

    .navbar-item + .navbar-item {
      margin-left: 20px;
    }
  }
}
</style>
