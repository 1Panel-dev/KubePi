<template>
  <div style="margin-top: 20px;margin-bottom: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="CPU Reservation" placeholder="e.g. 1000" clearable itemType="input" deviderName="mCPUs" v-model="form.requests.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Reservation" placeholder="e.g. 128" clearable itemType="input" deviderName="MiB" v-model="form.requests.memory" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
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
  props: {
    resourceObj: Object
  },
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
    transformation(parentFrom) {
      if(this.form.requests.memory) {
        parentFrom.requests.memory = this.form.requests.memory
      } else {
        delete parentFrom.requests.memory
      }
      if(this.form.requests.cpu) {
        parentFrom.requests.cpu = this.form.requests.cpu
      } else {
        delete parentFrom.requests.cpu
      }
      if(this.form.limits.memory) {
        parentFrom.limits.memory = this.form.limits.memory
      } else {
        delete parentFrom.limits.memory
      }
      if(this.form.limits.cpu) {
        parentFrom.limits.cpu = this.form.limits.cpu
      } else {
        delete parentFrom.limits.cpu
      }
    }
  },
}
</script>