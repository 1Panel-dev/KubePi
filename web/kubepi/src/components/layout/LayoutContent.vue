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
    color: #1f2937;
    background-color: #ffffff;
    overflow: auto;
    height: 100%;
    padding: 20px;
    border-radius: 4px;
    border: 1px solid #e2e8f0;
    box-shadow: 0 10px 30px rgba(15, 23, 42, 0.06);
    box-sizing: border-box;

    .content-container__header {
      line-height: 60px;
      font-size: 18px;
      color: #0f172a;
    }

    .content-container__toolbar {
      @include flex-row(space-between, center);
      margin-bottom: 10px;
    }
  }
</style>
