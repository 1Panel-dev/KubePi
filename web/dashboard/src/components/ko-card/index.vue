<template>
  <div style="margin-top: 20px">
    <el-card v-bind="$attrs" v-on="$listeners">
      <div slot="header" class="clearfix">
        <span>{{title}}</span>
        <span v-if="description" style="font-size: 12px;color: gray">   ({{description}})</span>
        <el-button v-if="isExpand" @click="changeExpand" icon="el-icon-minus" circle style="float: right; padding: 3px" />
        <el-button v-if="!isExpand" @click="changeExpand" icon="el-icon-plus" circle style="float: right; padding: 3px" />
      </div>
      <div v-if="isExpand">
        <slot></slot>
      </div>
    </el-card>
  </div>
</template>

<script>
export default {
  name: "KoCard",
  props: {
    title: String,
    description: String,
  },
  data() {
    return {
      isExpand: true,
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
</style>