<template>
  <div class="complex-table">
    <div class="complex-table__header" v-if="$slots.header || header">
      <slot name="header">{{ header }}</slot>
    </div>

    <div class="complex-table__toolbar" v-if="$slots.toolbar">
      <div>
        <slot name="toolbar">
        </slot>
      </div>
    </div>

    <div class="complex-table__body">
      <fu-table v-on="$listeners" v-bind="$attrs" :columns="columns" :local-key="localKey"
                @selection-change="handleSelectionChange">
        <slot></slot>
      </fu-table>
    </div>

    <div class="complex-table__pagination" v-if="$slots.pagination || paginationConfig">
      <slot name="pagination">
<!--          <ko-page  :pagination-config="paginationConfig" :next-token.sync="paginationConfig.nextToken"-->
<!--                   :current-page.sync="paginationConfig.currentPage" :page-size.sync="paginationConfig.pageSize"-->
<!--                   :items.sync="paginationConfig.items" :remain-count.sync="paginationConfig.remainCount"-->
<!--                   @change="search"></ko-page>-->
        <k8s-page :pagination-config="paginationConfig" :next-token.sync="paginationConfig.nextToken"
                  :page-size.sync="paginationConfig.pageSize" :items.sync="paginationConfig.items"
                  :remain-count.sync="paginationConfig.remainCount" @change="search">
        </k8s-page>
      </slot>
    </div>
  </div>
</template>

<script>


import K8sPage from "@/components/k8s-page"
export default {
  name: "ComplexTable",
  components: { K8sPage  },
  props: {
    columns: {
      type: Array,
      default: () => []
    },
    localKey: String, // 如果需要记住选择的列，则这里添加一个唯一的Key
    header: String,
    searchConfig: Object,
    paginationConfig: Object,
    selects: Array
  },
  data () {
    return {
      pageShow: false
    }
  },
  methods: {
    search () {
      this.$emit("search")
    },
    handleSelectionChange (val) {
      this.$emit("update:selects", val)
    },
  },
}
</script>

<style lang="scss">
  @import "~@/styles/common/mixins.scss";
  @import "~@/styles/common/variables.scss";

  .complex-table {
    .complex-table__header {
      @include flex-row(flex-start, center);
      line-height: 60px;
      font-size: 18px;
    }

    .complex-table__toolbar {
      @include flex-row(flex-end, center);
    }

    .complex-table__pagination {
      margin-top: 20px;
      @include flex-row(flex-end);
    }
  }

</style>
