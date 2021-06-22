<template>
  <div style="float: right">
    <el-pagination
            @size-change="handleSizeChange"
            @prev-click="pageAhead"
            @next-click="pageNext"
            :current-page="currentPage"
            :page-sizes="pageSizes"
            :page-size="pageSize"
            :total="total"
            layout="total, sizes, prev, slot, next">
      <el-button>{{ current }}</el-button>
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
    paginationConfig: Object,
    currentPage: {
      type: Number,
      default: 1
    },
    nextToken: {
      type: String,
      default: ""
    },
    remainCount: {
      type: Number,
      default: 0
    },
    items: {
      type: Number,
      default: 0
    },
  },
  data () {
    return {
      preToken: [""],
      total: 0,
      current: 1,
    }
  },
  watch: {
    paginationConfig: {
      handler (newValue) {
        if (newValue.remainCount > 0) {
          this.total = newValue.pageSize * this.current + newValue.remainCount
        }
        if (newValue.items > 0) {
          this.total = newValue.pageSize * (this.current - 1) + newValue.items
        }
      },
      deep: true
    },
  },
  methods: {
    handleSizeChange: function (size) {
      this.$emit("update:pageSize", size)
      this.$emit("update:nextToken", "")
      this.$emit("update:remainCount", 0)
      this.$emit("update:items", 0)
      this.$emit("change")
      this.current = 1
    },
    pageNext () {
      this.current++
      this.preToken.push(this.paginationConfig.nextToken)
      this.$emit("update:currentPage", this.currentPage)
      this.$emit("change")
    },
    pageAhead () {
      this.current--
      this.$emit("update:nextToken", this.preToken[this.preToken.length - 2])
      this.preToken.pop()
      this.$emit("update:currentPage", this.currentPage)
      this.$emit("change")
    },
  },
}
</script>

<style scoped>

</style>