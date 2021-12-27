<template>
  <div>
    <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="'CPU ' + $t('business.workload.reservation')" prop="requests.cpu">
            <ko-form-item placeholder="e.g. 1000" itemType="number" deviderName="mCPUs" v-model.number="form.requests.cpu" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.memory') + $t('business.workload.reservation')" prop="requests.memory">
            <ko-form-item placeholder="e.g. 128" itemType="number" deviderName="Mi" v-model.number="form.requests.memory" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="'CPU ' + $t('business.workload.limit')" prop="limits.cpu">
            <ko-form-item placeholder="e.g. 1000" itemType="number" deviderName="mCPUs" v-model.number="form.limits.cpu" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.memory') + $t('business.workload.limit')" prop="limits.memory">
            <ko-form-item placeholder="e.g. 128" itemType="number" deviderName="Mi" v-model.number="form.limits.memory" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"

export default {
  name: "KoResources",
  components: { KoFormItem },
  props: {
    resourceParentObj: Object,
    isReadOnly: Boolean,
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
      parentFrom.resources = { requests: {}, limits: {} }
      parentFrom.resources.requests.memory = this.form.requests.memory ? this.form.requests.memory + "Mi" : undefined
      parentFrom.resources.requests.cpu = this.form.requests.cpu ? this.form.requests.cpu + "m" : undefined
      parentFrom.resources.limits.memory = this.form.limits.memory ? this.form.limits.memory + "Mi" : undefined
      parentFrom.resources.limits.cpu = this.form.limits.cpu ? this.form.limits.cpu + "m" : undefined
    },
  },
  mounted() {
    if (this.resourceParentObj) {
      if (this.resourceParentObj.resources) {
        if (this.resourceParentObj.resources.requests) {
          if (this.resourceParentObj.resources.requests.memory) {
            this.form.requests.memory = memoryUnitConvert(this.resourceParentObj.resources.requests.memory)
          }
          if (this.resourceParentObj.resources.requests.cpu) {
            this.form.requests.cpu = cpuUnitConvert(this.resourceParentObj.resources.requests.cpu)
          }
        }
        if (this.resourceParentObj.resources.limits) {
          if (this.resourceParentObj.resources.limits.memory) {
            this.form.limits.memory = memoryUnitConvert(this.resourceParentObj.resources.limits.memory)
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