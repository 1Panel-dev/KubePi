<template>
  <div class="complex-table">
    <div class="complex-table__header" v-if="$slots.header || header">
      <slot name="header">{{ header }}</slot>
    </div>

    <div class="complex-table__toolbar">
        <template v-if="searchConfig">
          <el-row v-if="showFullTextSwitch">
            <el-col :span="8">
              <div  style="margin-top: 5px;  float: left">
                 <span class="spanClass">{{$t('commons.search.fullTextSearch')}}</span>
                 <el-switch class="interval"  v-model="isFullTextSearch" @change="OnIsFullTextSearchChange"/>
              </div>
            </el-col>  
            <el-col :span="6">
              <el-button size="mini" @click="search(true)"><i class="el-icon-search"/>{{this.$t('commons.button.search')}}</el-button>
            </el-col> 
            <el-col :span="10">
              <el-input  :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchConfig.keywords" @change="search(true)">
              </el-input>
            </el-col>
          </el-row>  
          <el-row v-if="!showFullTextSwitch">
            <el-col :span="12">
              <el-button size="mini" @click="search(true)"><i class="el-icon-search"/>{{this.$t('commons.button.search')}}</el-button>
            </el-col>  
            <el-col :span="12">
              <el-input  :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchConfig.keywords" @change="search(true)">
              </el-input>
            </el-col>
          </el-row> 
        </template>
        <slot name="toolbar">
        </slot>
    </div>

    <div class="complex-table__body">
      <fu-table :empty-text="$t('commons.table.empty_text')" v-on="$listeners" v-bind="$attrs" :columns="columns" :local-key="localKey" @selection-change="handleSelectionChange">
        <slot></slot>
      </fu-table>
    </div>

    <div class="complex-table__pagination" v-if="$slots.pagination || paginationConfig">
      <slot name="pagination">
        <fu-table-pagination :current-page.sync="paginationConfig.currentPage" :page-size.sync="paginationConfig.pageSize" v-bind="paginationConfig" @current-change="search()" @size-change="search(true)" />
      </slot>
    </div>
  </div>
</template>

<script>
export default {
  name: "ComplexTable",
  components: {},
  props: {
    columns: {
      type: Array,
      default: () => [],
    },
    localKey: String, // 如果需要记住选择的列，则这里添加一个唯一的Key
    header: String,
    searchConfig: Object,
    paginationConfig: Object,
    selects: Array,
    showFullTextSwitch: {//是否显示全文本搜索选项
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pageShow: false,
      isFullTextSearch :false,//是否全文本搜索
      intervalId:null,//定时
    }
  },
  methods: {
    search(resetPage) {
      this.$emit("search", resetPage)
    },
    handleSelectionChange(val) {
      this.$emit("update:selects", val)
    },
    OnIsFullTextSearchChange(val){
      this.$emit("update:isFullTextSearch", val)
    },
    stopTimeTick(){
      if (!this.intervalId) {
        return;
      }
      clearInterval(this.intervalId); //清除计时器
      this.intervalId = null; //设置为null
    },
    startTimeTick(tick){
      // 计时器为空，操作
      this.intervalId = setInterval(() => {
          this.stopTimeTick()
          this.search(true);
          this.startTimeTick(tick)
      },tick*1000);
    }
  },
  watch: {
    /*监控自动刷新变量*/
   "$store.state.app.autorefresh":{
    handler:function(newVal,oldVal){
      if(!newVal ||  newVal=='-1' || newVal=='undefined' || newVal==''){
        this.stopTimeTick();
      } else {
        this.stopTimeTick();
        // 计时器为空，操作
        this.startTimeTick(parseInt(newVal));
      }
    }
   }
  },
  destroyed(){
    // 在页面销毁后，清除计时器
    this.stopTimeTick();
  }
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
    margin-bottom: 10px;
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
