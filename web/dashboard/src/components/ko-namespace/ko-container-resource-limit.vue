<template>
  <div class="tab-content">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="CPU Reservation" placeholder="e.g. 1000" clearable itemType="number" deviderName="mCPUs"
                      v-model="form.requestsCpu"></ko-form-item>
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Reservation" placeholder="e.g. 128" clearable itemType="number" deviderName="MiB"
                      v-model="form.requestsMemory"/>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="CPU Limit" placeholder="e.g. 128" clearable itemType="number" deviderName="mCPUs"
                      v-model="form.limitsCpu"/>
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Limit" placeholder="e.g. 128" clearable itemType="number" deviderName="MiB"
                      v-model="form.limitsMemory"/>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoContainerResourceLimit",
  components: { KoFormItem },
  props: {
    resourceLimit: Object
  },
  data () {
    return {
      form: {
        limitsCpu: "",
        limitsMemory: "",
        requestsCpu: "",
        requestsMemory: ""
      }
    }
  },
  watch: {
    form(val) {
      this.$emit("update:resourceLimit", val)
    }
  },
  mounted () {
    this.form = this.resourceLimit
  }
}
</script>

<style scoped>
    .tab-content {
        margin-top: 20px;
        margin-bottom: 20px;
    }
</style>