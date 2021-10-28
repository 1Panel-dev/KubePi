<template>
  <div class="horizontal-header">
    <div class="header-left">
      <sidebar-toggle-button/>
    </div>
    <div class="header-right">
      <help />
      <language-switch style="margin-left: 20px"></language-switch>
      <user-setting style="margin-left: 20px"></user-setting>
    </div>
  </div>
</template>

<script>
import SidebarToggleButton from "@/components/layout/sidebar/SidebarToggleButton"
import UserSetting from "@/business/app-layout/header-components/UserSetting"
import LanguageSwitch from "@/business/app-layout/header-components/LanguageSwitch";
import Help from "@/business/app-layout/header-components/Help";

export default {
  name: "HorizontalHeader",
  components: {LanguageSwitch, UserSetting, SidebarToggleButton, Help}
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