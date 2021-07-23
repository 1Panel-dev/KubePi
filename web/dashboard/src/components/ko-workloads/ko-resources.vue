<template>
  <div style="margin-top: 20px">
    <ko-card title="Resources">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="CPU Reservation" prop="requests.cpu">
              <ko-form-item placeholder="e.g. 1000" itemType="number" deviderName="mCPUs" v-model.number="form.requests.cpu" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Memory Reservation" prop="requests.memory">
              <ko-form-item placeholder="e.g. 128" itemType="number" deviderName="MiB" v-model.number="form.requests.memory" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="CPU Limit" prop="limits.cpu">
              <ko-form-item placeholder="e.g. 1000" itemType="number" deviderName="mCPUs" v-model.number="form.limits.cpu" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Memory Limit" prop="limits.memory">
              <ko-form-item placeholder="e.g. 128" itemType="number" deviderName="MiB" v-model.number="form.limits.memory" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"
import { cpuUnitConvert, memeryUnitConvert } from "@/utils/unitConvert"

export default {
  name: "KoResources",
  components: { KoFormItem, KoCard },
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
      if (!parentFrom.resources) {
        parentFrom.resources = {}
      }
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