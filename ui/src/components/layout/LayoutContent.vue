<template>
  <div class="content-container">
    <div class="content-container__header" v-if="$slots.header || header">
      <slot name="header">
        <back-button :path="backPath" :name="backName" :to="backTo" v-if="showBack"></back-button>
        {{ header }}
      </slot>
    </div>
    <div class="content-container__toolbar" v-if="$slots.toolbar">
      <slot name="toolbar"></slot>
    </div>
    <slot></slot>
  </div>
</template>

<script>
import BackButton from "@/components/back-button/back-button"
export default {
  name: "LayoutContent",
  components: { BackButton },
  props: {
    header: String,
    description: String,
    backPath: String,
    backName: String,
    backTo: Object
  },
  computed: {
    showBack({backPath, backName, backTo}) {
      return backPath || backName || backTo
    }
  }
}
</script>

<style lang="scss">
  @import "~@/styles/common/mixins.scss";
  @import "~@/styles/common/variables";

  .content-container {
    transition: 0.3s;
    color: white;
    background-color: #222629;
    overflow: auto;
    height: 100%;
    padding: 20px;
    border-radius: 4px;
    box-shadow: 0 1px 4px 0 rgba(0, 0, 0, .14);
    box-sizing: border-box;

    .content-container__header {
      line-height: 60px;
      font-size: 18px;
      color: #FFFFFF;
    }

    .content-container__toolbar {
      @include flex-row(space-between, center);
      margin-bottom: 10px;
    }
  }
</style>
