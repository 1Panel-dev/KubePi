<template>
  <div class="complex-table">
    <div class="complex-table__header" v-if="$slots.header || header">
      <slot name="header">{{ header }}</slot>
    </div>

    <div class="complex-table__toolbar" v-if="$slots.toolbar || searchConfig">
      <div>


        <div v-if="searchConfig">

          <!--          <el-input :placeholder="$t('commons.button.search')"-->
          <!--                    suffix-icon="el-icon-search" clearable v-model="searchConfig.keywords"-->
          <!--                    @change="search(true)">-->
          <!--          </el-input>-->
          <fu-search-bar v-bind="searchConfig" @exec="search">
            <template #complex>
              <slot name="complex"></slot>
            </template>
            <slot name="buttons"></slot>
<!--            <fu-table-column-select :columns="columns"/>-->
          </fu-search-bar>
        </div>


        <slot name="toolbar">

        </slot>
      </div>
    </div>

    <div class="complex-table__body">
      <fu-table :empty-text="$t('commons.table.empty_text')" v-on="$listeners" v-bind="$attrs" :columns="columns" :local-key="localKey"
                @selection-change="handleSelectionChange">
        <slot></slot>
      </fu-table>
    </div>

    <div class="complex-table__pagination" v-if="$slots.pagination || paginationConfig">
      <slot name="pagination">
        <fu-table-pagination :current-page.sync="paginationConfig.currentPage"
                             :page-size.sync="paginationConfig.pageSize"
                             v-bind="paginationConfig"
                             @change="search"/>
      </slot>
    </div>
  </div>
</template>

<script>

export default {
  name: "ComplexTable",
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
  data() {
    return {
      condition: {},
    }
  },
  methods: {
    search(condition, e) {
      if (condition) {
        this.condition = condition
      }
      this.$emit("search", this.condition, e)
    },
    handleSelectionChange(val) {
      // this.selects = val
      this.$emit("update:selects", val)
    }
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
    margin-bottom: 10px;
  }

  .complex-table__pagination {
    margin-top: 20px;
    @include flex-row(flex-end);
  }
}

</style>
