<template>
  <div v-loading="loading">
    <h3>{{$t('business.workload.toleration')}}</h3>
    <el-form label-position="top">
      <el-collapse v-model="activeNames">
        <el-collapse-item style="margin-top: 10px" :name="index" v-for="(item, index) in form.spec.tolerations" :key="index">
          <div v-if="item.key" slot="title" class="spanInFormStyle"><span :title="item.key">{{item.key}}</span></div>
          <el-row :gutter="20">
            <el-col v-if="item.key" :span="12">
              <el-form-item style="margin-left: 20px;" label="key">
                <div class="spanInFormStyle"><span :title="item.key">{{item.key}}</span></div>
              </el-form-item>
            </el-col>
            <el-col v-if="item.operator" :span="12">
              <el-form-item style="margin-left: 20px;" label="operator">
                <span>{{item.operator}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.effect" :span="12">
              <el-form-item style="margin-left: 20px;" label="effect">
                <span>{{item.effect}}</span>
              </el-form-item>
            </el-col>
            <el-col v-if="item.tolerationSeconds" :span="12">
              <el-form-item style="margin-left: 20px;" :label="$t('business.workload.toleration_seconds')">
                <span>{{item.tolerationSeconds}}</span>
              </el-form-item>
            </el-col>
          </el-row>
        </el-collapse-item>
      </el-collapse>
    </el-form>
    <br>
  </div>
</template>

<script>
export default {
  name: "KoDetailToleration",
  props: {
    yamlInfo: Object,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec?.tolerations) {
            this.form = newYamlInfo
            this.loading = false
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      loading: true,
      activeNames: 0,
      form: {
        spec: {
          tolerations: [],
        },
      },
    }
  },
}
</script>

<style scoped>
</style>
