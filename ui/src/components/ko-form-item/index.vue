<template>
  <div class="formClass">
    <div v-if="!withoutLabel">
      <label style="color: gray; margin-left: 15px; font-size: 12px">
        {{labelName}}
        <el-tooltip v-if="tipContant" class="item" effect="dark" :content="tipContant" placement="top-start">
          <i class="el-icon-question" />
        </el-tooltip>

        <el-select v-if="itemType==='select'" style="width: 100%; margin-top: 2px" size="mini" v-bind="$attrs" v-on="$listeners">
          <el-option v-for="(item, index) in selections" :key="index" :label="item.label" :value="item.value" />
        </el-select>

        <el-input v-if="itemType==='input'" style="width:100%; margin-top: 2px" size="mini" v-bind="$attrs" v-on="$listeners"></el-input>

        <div style="display: block; margin: 10px 0 0 20px">
          <el-radio-group v-if="itemType==='radio'" v-bind="$attrs" v-on="$listeners">
            <el-radio v-for="(item, index) in radios" :key="index" :label="item.value" style="display: block; line-height: 25px">{{item.label}}</el-radio>
          </el-radio-group>
        </div>
      </label>
      <div class="deviderClass" v-if="deviderName">
        <span>{{deviderName}}</span>
      </div>
    </div>

    <div v-if="withoutLabel">
      <el-input v-if="itemType==='input'" style="width:100%; margin-top: 2px" size="mini" v-bind="$attrs" v-on="$listeners"></el-input>

      <el-select v-if="itemType==='select'" style="width: 95%; margin-top: 2px" size="mini" v-bind="$attrs" v-on="$listeners">
        <el-option v-for="(item, index) in selections" :key="index" :label="item.label" :value="item.value" />
      </el-select>

      <div style="display: block; margin: 10px 0 0 20px">
          <el-radio-group v-if="itemType==='radio'" v-bind="$attrs" v-on="$listeners">
            <el-radio v-for="(item, index) in radios" :key="index" :label="item.value" style="display: block; line-height: 25px">{{item.label}}</el-radio>
          </el-radio-group>
        </div>
    </div>
  </div>
</template>
          
<script>
export default {
  name: "KoFormItem",
  props: {
    itemType: String, // input, select, radio
    labelName: String,
    selections: Array, // 如果是 select
    deviderName: String, // 如果需要加上单位
    radios: Array, // 如果是 radio
    tipContant: String, // 如果需要tooltips
    withoutLabel: Boolean,
  },
}
</script>
          
<style>
.formClass {
  border: solid 1px rgb(228, 204, 204);
  border-radius: 5px;
  position: relative;
  display: table;
  width: 98%;
}
.formClass .el-input__inner {
  border: none !important;
  background: none !important;
  color: white;
  font-size: 12px;
  font-weight: 600;
}
.formClass .el-input__inner::placeholder {
  color: #aaaaaa;
}
.formClass .deviderClass {
  display: table-cell;
  width: 1%;
  white-space: nowrap;
  vertical-align: middle;
  border-left: 1px solid #dddddd;
  padding: 0 10px;
}
</style>>