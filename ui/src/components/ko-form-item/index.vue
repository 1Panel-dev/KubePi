<template>
  <div class="formClass">
    <label style="color: gray; margin-left: 16px; font-size: 12px">
      {{labelName}}
      <el-select v-if="itemType==='select'" style="width: 100%; margin-top: 2px" v-bind="$attrs" v-on="$listeners">
        <el-option v-for="(item, index) in selections" :key="index" :label="item.label" :value="item.value" />
      </el-select>
      <el-input v-if="itemType==='input'" style="width:100%; margin-top: 2px" v-bind="$attrs" v-on="$listeners"></el-input>
    </label>
  </div>
</template>
          
<script>
export default {
  name: "KoFormItem",
  props: {
    itemType: String,
    labelName: String,
    selections: Array,
  },
}
</script>
          
<style>
.formClass {
  padding: 5px;
  border: solid 1px rgb(228, 204, 204);
  border-radius: 5px;
}
.formClass .el-input__inner {
  border: none !important;
  background: none !important;
  color: white;
  font-size: 16px;
  font-weight: 600;
}
</style>>