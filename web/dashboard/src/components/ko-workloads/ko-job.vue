<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Completions" tipContant="The number of successfully finished pods the job should be run with." clearable itemType="input" deviderName="Times" v-model="form.requests.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Parallelism" tipContant="The maximum number of pods the job should run at any given time." clearable itemType="input" deviderName="Times" v-model="form.requests.memory" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Back Off Limit" tipContant="The number of retries before marking this job failed." clearable itemType="input" deviderName="Times" v-model="form.limits.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Active Deadline" tipContant="The duration that the job may be active before the system tries to terminate it." clearable itemType="input" deviderName="Seconds" v-model="form.limits.memory" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoJob",
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