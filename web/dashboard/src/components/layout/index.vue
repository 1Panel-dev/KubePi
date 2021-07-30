<template>
  <el-container class="layout-container">
    <slot>
      <layout-sidebar/>
      <layout-main>
        <layout-header>
          <slot name="header"></slot>
        </layout-header>
        <layout-view/>
        <!--        <layout-footer>-->
        <!--          <slot name="footer"></slot>-->
        <!--        </layout-footer>-->
      </layout-main>
    </slot>
  </el-container>
</template>

<script>
import LayoutSidebar from "./LayoutSidebar"
import LayoutMain from "./LayoutMain"
import LayoutHeader from "./LayoutHeader"
// import LayoutFooter from "./LayoutFooter"
import LayoutView from "./LayoutView"
// import FuSplitPane from "@/components/split-pane/FuSplitPane.vue"

export default {
  name: "Layout",
  components: {LayoutView, LayoutHeader, LayoutMain, LayoutSidebar},
  // components: { LayoutView, LayoutFooter, LayoutHeader, LayoutMain, LayoutSidebar, FuSplitPane },
  data() {
    return {
      isRefresh: false,
    }
  },
  computed: {
    buttomHeight() {
      return this.$store.state.app.buttomHeight
    },
  },
  watch: {
    buttomHeight: function () {
      this.isRefresh = !this.isRefresh
    },
  },
  methods: {},
}
</script>

<style lang="scss">
@import "~@/styles/common/variables";

.layout-container {
  min-width: 1024px;
  height: 100%;
  background-color: $layout-bg-color;
}

.main-container {
  position: relative;
}

.sidebar-container {
  position: relative;
  transition: width 0.28s;
  width: $sidebar-open-width;
  min-width: $sidebar-open-width;
  background-color: $sidebar-bg-color;
  background-image: $sidebar-bg-gradient;

  &.is-collapse {
    width: $sidebar-close-width;
    min-width: $sidebar-close-width;
  }
}

.header-container {
  height: $header-height;
  padding: 0 $header-padding;
}

.view-container {
  display: block;
  flex: auto;
  overflow: auto;
  box-sizing: border-box;
  padding: $view-padding;
}
</style>