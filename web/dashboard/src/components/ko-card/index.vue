<template>
  <div>
    <el-collapse v-bind="$attrs" v-on="$listeners" v-model="activeName">
      <el-collapse-item name="1">
        <template slot="title">
          <div style="margin-left: 5px">
            {{title}}
          </div>
        </template>
        <div style="margin: 5px">
          <slot></slot>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script>
export default {
  name: "KoCard",
  props: {
    title: String,
  },
  data() {
    return {
      activeName:"1"
    }
  },
  methods: {
    changeExpand() {
      this.isExpand = !this.isExpand
    },
  },
}
</script>

<style>
    .el-collapse {
        border: 0;
    }
</style>