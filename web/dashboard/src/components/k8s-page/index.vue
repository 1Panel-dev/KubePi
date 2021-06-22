<template>
  <div style="float: right">
    <el-pagination
            @size-change="handleSizeChange"
            :page-size="pageSize"
            :page-sizes="pageSizes"
            layout="sizes,slot">
      <el-button icon="el-icon-arrow-left" @click="pageAhead()" :disabled="current === 1"></el-button>
      <el-button>{{ current }}</el-button>
      <el-button icon="el-icon-arrow-right" @click="pageNext()" :disabled="nextToken===''"></el-button>
    </el-pagination>
  </div>
</template>

<script>
export default {
  name: "K8sPage",
  props: {
    pageSize: {
      type: Number,
      default: 10
    },
    paginationConfig: Object,
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
    pageSizes: {
      type: Array,
      default: function () {
        return [5, 10, 20, 50, 100]
      }
    },
  },
  data () {
    return {
      current: 1,
      preToken: [""],
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
      this.$emit("change")
    },
    pageAhead () {
      this.current--
      this.$emit("update:nextToken", this.preToken[this.preToken.length - 2])
      this.preToken.pop()
      this.$emit("change")
    },
  }
}
</script>

<style scoped>

</style>