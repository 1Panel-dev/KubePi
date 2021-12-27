<template>
  <div>
    <el-form label-position="top" :model="form" ref="form">
      <el-row :gutter="20" style="margin-left: 5px">
        <el-col :span="6">
          <el-form-item :label="$t('business.workload.type')" prop="type">
            <ko-form-item itemType="select2" v-model="form.type" :selections="type_list" />
          </el-form-item>
        </el-col>
      </el-row>

      <div v-if="form.type === 'Container' || form.type === 'Pod'">
        <el-row :gutter="20" style="margin-left: 5px">
          <el-col :span="6">
            <el-form-item label="Min CPU" prop="cpuMin">
              <ko-form-item itemType="number" deviderName="mCPUs" v-model="form.cpuMin" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="Max CPU" prop="memoryRequest">
              <ko-form-item itemType="number" deviderName="mCPUs" v-model="form.cpuMax" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" style="margin-left: 5px">
          <el-col :span="6">
            <el-form-item label="Min Memory" prop="memoryMin">
              <ko-form-item itemType="number" deviderName="Mi" v-model="form.memoryMin" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="Max Memory" prop="memoryMax">
              <ko-form-item itemType="number" deviderName="Mi" v-model="form.memoryMax" />
            </el-form-item>
          </el-col>
        </el-row>
      </div>

      <el-row :gutter="20" style="margin-left: 5px" v-if="form.type === 'PersistentVolumeClaim'">
        <el-col :span="6">
          <el-form-item label="Min Storage" prop="storageMin">
            <ko-form-item itemType="number" deviderName="Mi" v-model="form.storageMin" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="Max Storage" prop="storageMax">
            <ko-form-item itemType="number" deviderName="Mi" v-model="form.storageMax" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import { cpuUnitConvert, memoryUnitConvert } from "@/utils/unitConvert"

export default {
  name: "KoLimitRange",
  components: { KoFormItem },
  props: {
    limitRangeObj: Object,
  },
  watch: {
    limitRangeObj: {
      handler(newVal) {
        if (newVal) {
          this.form.type = newVal.type
          if (newVal.type === "PersistentVolumeClaim") {
            this.form.storageMin = newVal.min?.storage ? memoryUnitConvert(newVal.min.storage) : ""
            this.form.storageMax = newVal.max?.storage ? memoryUnitConvert(newVal.max.storage) : ""
          } else {
            this.form.cpuMin = newVal.min?.cpu ? cpuUnitConvert(newVal.min.cpu) : ""
            this.form.cpuMax = newVal.max?.cpu ? cpuUnitConvert(newVal.max.cpu) : ""
            this.form.memoryMin = newVal.min?.memory ? memoryUnitConvert(newVal.min.memory) : ""
            this.form.memoryMax = newVal.max?.memory ? memoryUnitConvert(newVal.max.memory) : ""
          }
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      form: {
        type: "",
        cpuMin: "",
        cpuMax: "",
        memoryMin: "",
        memoryMax: "",
        storageMin: "",
        storageMax: "",
      },
      type_list: ["Pod", "Container", "PersistentVolumeClaim"],
    }
  },
  methods: {
    isPVC() {
      return this.form.type === "PersistentVolumeClaim"
    },
    transformation(spec) {
      if (this.form.type === "") {
        return
      }
      if (!this.isPVC()) {
        spec.limits = [
          {
            type: this.form.type,
            max: {
              cpu: this.form.cpuMax ? this.form.cpuMax + "m" : undefined,
              memory: this.form.memoryMax ? this.form.memoryMax + "Mi" : undefined,
            },
            min: {
              cpu: this.form.cpuMin ? this.form.cpuMin + "m" : undefined,
              memory: this.form.memoryMin ? this.form.memoryMin + "Mi" : undefined,
            },
          },
        ]
      } else {
        spec.limits = [
          {
            type: this.form.type,
            max: {
              storage: this.form.storageMax ? this.form.storageMax + "Mi" : undefined,
            },
            min: {
              storage: this.form.storageMin ? this.form.storageMin + "Mi" : undefined,
            },
          },
        ]
      }
    },
  },
}
</script>

<style scoped>
</style>
