<template>
  <div class="horizontal-header">
    <div class="header-left">
      <sidebar-toggle-button/>
    </div>
    <div class="header-right">
      <namespace-switch></namespace-switch>
    </div>
  </div>
</template>

<script>
import SidebarToggleButton from "@/components/layout/sidebar/SidebarToggleButton"
import NamespaceSwitch from "@/business/app-layout/header-components/NamspaceSwitch"
export default {
  name: "HorizontalHeader",
  components: { NamespaceSwitch, SidebarToggleButton }
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