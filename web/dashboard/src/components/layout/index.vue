<template>
  <el-container class="layout-container">
    <slot>
      <layout-sidebar />
      <layout-main>
        <layout-header>
          <slot name="header"></slot>
        </layout-header>
        <div style="width:100%;height:100%;">
          <fu-split-pane :bottom="buttomHeight" direction="vertical" resizer-type="line" :resizer-style="{background:'#424649',height:'2px'}">
            <div slot="top">
              <layout-view />
            </div>
            <div slot="bottom">
              <layout-footer>
                <slot name="footer"></slot>
              </layout-footer>
            </div>
          </fu-split-pane>
        </div>
      </layout-main>
    </slot>
  </el-container>
</template>

<script>
import LayoutSidebar from "./LayoutSidebar"
import LayoutMain from "./LayoutMain"
import LayoutHeader from "./LayoutHeader"
import LayoutFooter from "./LayoutFooter"
import LayoutView from "./LayoutView"
import FuSplitPane from "@/components/split-pane/FuSplitPane.vue"

export default {
  name: "Layout",
  components: { LayoutView, LayoutFooter, LayoutHeader, LayoutMain, LayoutSidebar, FuSplitPane },
  computed: {
    buttomHeight() {
      return this.$store.state.app.buttomHeight
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