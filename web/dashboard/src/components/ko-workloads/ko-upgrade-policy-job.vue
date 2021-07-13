<template>
  <div style="margin-top: 20px">
    <ko-card title="Update Policy">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-row>
              <el-form-item label="Completions" prop="completions">
                <ko-form-item deviderName="Times" itemType="number" v-model.number="form.completions" />
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item label="Parallelism" prop="parallelism">
                <ko-form-item deviderName="Times" itemType="number" v-model.number="form.parallelism" />
              </el-form-item>
            </el-row>
            <el-row>
              <el-form-item label="Back Off Limit" prop="backoffLimit">
                <ko-form-item deviderName="Times" itemType="number" v-model="form.backoffLimit" />
              </el-form-item>
            </el-row>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Active Deadline" prop="activeDeadlineSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model="form.activeDeadlineSeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Pod Active Deadline" prop="template.spec.terminationGracePeriodSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model="form.template.spec.terminationGracePeriodSeconds" />
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

export default {
  name: "KoUpgradePolicyJob",
  components: { KoFormItem, KoCard },
  props: {
    upgradePolicyParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      form: {
        activeDeadlineSeconds: null,
        backoffLimit: null,
        completions: null,
        parallelism: null,
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
          },
        },
      },
    }
  },
  methods: {
    transformation(parentFrom) {
      if (this.form.activeDeadlineSeconds) {
        parentFrom.activeDeadlineSeconds = this.form.activeDeadlineSeconds
      }
      if (this.form.backoffLimit) {
        parentFrom.backoffLimit = this.form.backoffLimit
      }
      if (this.form.completions) {
        parentFrom.completions = this.form.completions
      }
      if (this.form.parallelism) {
        parentFrom.parallelism = this.form.parallelism
      }
      if (this.form.template.spec.activeDeadlineSeconds) {
        parentFrom.template.spec.activeDeadlineSeconds = this.form.template.spec.activeDeadlineSeconds
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      if (this.upgradePolicyParentObj.activeDeadlineSeconds) {
        this.form.activeDeadlineSeconds = this.upgradePolicyParentObj.activeDeadlineSeconds
      }
      if (this.upgradePolicyParentObj.backoffLimit) {
        this.form.backoffLimit = this.upgradePolicyParentObj.backoffLimit
      }
      if (this.upgradePolicyParentObj.completions) {
        this.form.completions = this.upgradePolicyParentObj.completions
      }
      if (this.upgradePolicyParentObj.parallelism) {
        this.form.parallelism = this.upgradePolicyParentObj.parallelism
      }
      if (this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds) {
        this.form.template.spec.terminationGracePeriodSeconds = this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds
      }
    }
  },
}
</script>