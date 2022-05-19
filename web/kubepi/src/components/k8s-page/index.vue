<template>
  <div style="margin-top: 10px; margin-right: 50px">
    <el-button :disabled="ko_currentPage===1 || ko_btnLoading" size="mini" class="el-pagination-button" icon="el-icon-arrow-left" @click="getPageAhead()" />
    <span class="el-pager-li" >{{ko_currentPage}}</span>
    <el-button :disabled="!ko_nextToken || ko_btnLoading" size="mini" class="el-pagination-button" icon="el-icon-arrow-right" @click="getPageNext()" />
  </div>
</template>
          
<script>
export default {
  name: "K8sPage",
  props: {
    nextToken: String,
    currentPage: Number,
  },
  watch: {
    nextToken: {
      handler(newName) {
        this.ko_nextToken = newName
      },
      immediate: true,
    },
    currentPage: {
      handler(newName) {
        this.ko_currentPage = newName ? newName : this.ko_currentPage
      },
      immediate: true,
    },
  },
  data() {
    return {
      ko_aheadToken: [],
      ko_continueToken: "",
      ko_nextToken: "",
      ko_currentPage: 1,
      ko_btnLoading: false,
    }
  },
  methods: {
    getPageAhead() {
      this.ko_btnLoading = true
      this.ko_currentPage--
      this.ko_continueToken = this.ko_aheadToken[this.ko_aheadToken.length - 1]
      this.ko_aheadToken.pop()
      this.$emit("pageChange", {token: this.ko_continueToken, page: this.ko_currentPage})
      this.ko_btnLoading = false
    },
    getPageNext() {
      this.ko_btnLoading = true
      this.ko_currentPage++
      this.ko_aheadToken.push(this.ko_continueToken)
      this.ko_continueToken = this.ko_nextToken
      this.$emit("pageChange", {token: this.ko_continueToken, page: this.ko_currentPage})
      this.ko_btnLoading = false
    },
  },
}
</script>
          
<style>
.el-pagination-button {
    color: #61616a;
    background-color: #1a1a1e;
    border: none;
    padding: 0 6px;

    display: inline-block;
    font-size: 14px;
    min-width: 35.5px;
    height: 28px;
    line-height: 28px;
    vertical-align: top;
    box-sizing: border-box;
}

.el-pager-li {
    padding: 0 4px;
    background: #1a1a1e;
    font-size: 13px;
    min-width: 35.5px;
    height: 28px;
    line-height: 28px;
    box-sizing: border-box;
}
</style>