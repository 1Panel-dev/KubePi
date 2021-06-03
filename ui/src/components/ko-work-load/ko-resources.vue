<template>
  <div style="margin-top: 20px">
    <el-row :gutter="40">
      <el-col :span="12">
        <ko-form-item labelName="CPU Reservation" placeholder="e.g. 1000" clearable itemType="input" deviderName="mCPUs" v-model="form.requests.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Reservation" placeholder="e.g. 128" clearable itemType="input" deviderName="MiB" v-model="form.requests.memory" />
      </el-col>
    </el-row>
    <el-row :gutter="40" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="CPU Limit" placeholder="e.g. 1000" clearable itemType="input" deviderName="mCPUs" v-model="form.limits.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Limit" placeholder="e.g. 128" clearable itemType="input" deviderName="MiB" v-model="form.limits.memory" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoResources",
  components: { KoFormItem },
  data() {
    return {
      form: {
        requests: {
          memory: "",
          cpu: "",
        },
        limits: {
          memory: "",
          cpu: "",
        },
      },
    }
  },
  methods: {
  },
}
</script>