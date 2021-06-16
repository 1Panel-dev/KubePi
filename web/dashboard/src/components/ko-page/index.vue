<template>
  <div style="float: right">
    <br>
    <el-pagination
            @size-change="handleSizeChange"
            @prev-click="pageAhead"
            @next-click="pageNext"
            :current-page="currentPage"
            :page-sizes="pageSizes"
            :page-size="pageSize"
            :total="total"
            layout="total, sizes, prev, slot, next">
      <el-button>{{ currentPage }}</el-button>
    </el-pagination>
  </div>
</template>

<script>
export default {
  name: "KoPage",
  props: {
    pageSizes: {
      type: Array,
      default: function () {
        return [5, 10, 20, 50, 100]
      }
    },
    pageSize: {
      type: Number,
      default: 10
    },
    page: Object
  },
  data () {
    return {
      nextToken: "",
      preToken: [""],
      total: 0,
      currentPage: 1
    }
  },
  watch: {
    continue: {},
    page: {
      handler () {
        if (this.page ) {
          if (this.page.remainCount > 0) {
            this.total = this.page.pageSize * this.currentPage + this.page.remainCount
          }
          if (this.page.items > 0 ){
            this.total = this.page.pageSize * (this.currentPage-1) + this.page.items
          }
        }
      },
      immediate: true,
      deep: true
    }
  },
  methods: {
    handleSizeChange: function (size) {
      const page = {
        pageSize: size,
        nextToken: "",
        remainingItemCount: 0,
      }
      this.$emit("update:page", page)
      this.$emit("change")
      this.currentPage = 1
    },
    pageNext () {
      this.currentPage++
      this.preToken.push(this.page.nextToken)
      this.$emit("change")
    },
    pageAhead () {
      this.currentPage--
      this.page.nextToken = this.preToken[this.preToken.length - 2]
      this.preToken.pop()
      this.$emit("update:page", this.page)
      this.$emit("change")
    }
  }
}
</script>

<style scoped>

</style>