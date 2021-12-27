<template>
  <div class="tab-content">
    <el-form label-position="top" ref="form" :rules="rules" :model="form">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="CPU Reservation" prop="requestsCpu">
            <ko-form-item placeholder="e.g. 1000" clearable itemType="number" deviderName="mCPUs" v-model="form.requestsCpu" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Memory Reservation" prop="requestsMemory">
            <ko-form-item placeholder="e.g. 128" clearable itemType="number" deviderName="Mi" v-model="form.requestsMemory" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="CPU Limit" prop="limitsCpu">
            <ko-form-item placeholder="e.g. 128" clearable itemType="number" deviderName="mCPUs" v-model="form.limitsCpu" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Memory Limit" prop="limitsMemory">
            <ko-form-item placeholder="e.g. 128" clearable itemType="number" deviderName="Mi" v-model="form.limitsMemory" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoContainerResourceLimit",
  components: { KoFormItem },
  props: {
    resourceLimit: Object,
  },
  data() {
    return {
      form: {
        limitsCpu: "",
        limitsMemory: "",
        requestsCpu: "",
        requestsMemory: "",
      },
      rules: [],
    }
  },
  watch: {
    form(val) {
      this.$emit("update:resourceLimit", val)
    },
  },
  mounted() {
    this.form = this.resourceLimit
  },
}
</script>

<style scoped>
.tab-content {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>