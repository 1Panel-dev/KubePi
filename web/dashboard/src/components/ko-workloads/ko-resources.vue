<template>
  <div style="margin-top: 20px;margin-bottom: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="CPU Reservation" placeholder="e.g. 1000" clearable itemType="number" deviderName="mCPUs" v-model="form.requests.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Reservation" placeholder="e.g. 128" clearable itemType="number" deviderName="MiB" v-model="form.requests.memory" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="CPU Limit" placeholder="e.g. 1000" clearable itemType="number" deviderName="mCPUs" v-model="form.limits.cpu" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Memory Limit" placeholder="e.g. 128" clearable itemType="number" deviderName="MiB" v-model="form.limits.memory" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import {cpuUnitConvert, memeryUnitConvert} from '@/utils/unitConvert'

export default {
  name: "KoResources",
  components: { KoFormItem },
  props: {
    resourceParentObj: Object,
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
      if (this.form.requests.memory || this.form.requests.cpu) {
        parentFrom.resources.requests = {}
        if (this.form.requests.memory) {
          parentFrom.resources.requests.memory = this.form.requests.memory + "Mi"
        }
        if (this.form.requests.cpu) {
          parentFrom.resources.requests.cpu = this.form.requests.cpu + "m"
        }
      }
      if (this.form.limits.memory || this.form.limits.cpu) {
        parentFrom.resources.limits = {}
        if (this.form.limits.memory) {
          parentFrom.resources.limits.memory = this.form.limits.memory + "Mi"
        }
        if (this.form.limits.cpu) {
          parentFrom.resources.limits.cpu = this.form.limits.cpu + "m"
        }
      }
    },
  },
  mounted() {
    if (this.resourceParentObj) {
      if (this.resourceParentObj.resources) {
        if (this.resourceParentObj.resources.requests) {
          if (this.resourceParentObj.resources.requests.memory) {
            this.form.requests.memory = memeryUnitConvert(this.resourceParentObj.resources.requests.memory)
          }
          if (this.resourceParentObj.resources.requests.cpu) {
            this.form.requests.cpu = cpuUnitConvert(this.resourceParentObj.resources.requests.cpu)
          }
        }
        if (this.resourceParentObj.resources.limits) {
          if (this.resourceParentObj.resources.limits.memory) {
            this.form.limits.memory = memeryUnitConvert(this.resourceParentObj.resources.limits.memory)
          }
          if (this.resourceParentObj.resources.limits.cpu) {
            this.form.limits.cpu = cpuUnitConvert(this.resourceParentObj.resources.limits.cpu)
          }
        }
      }
    }
  },
}
</script>