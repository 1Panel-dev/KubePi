<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.upgrade_policy')">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.completions')" prop="completions">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.completions" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.restart_strategy')" prop="template.spec.restartPolicy">
              <ko-form-item itemType="radio" v-model="form.template.spec.restartPolicy" :radios="restart_policy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.parallelism')" prop="parallelism">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.parallelism" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.back_off_limit')" prop="backoffLimit">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.backoffLimit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.active_deadline')" prop="activeDeadlineSeconds">
              <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.activeDeadlineSeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.termination_grace_period')" prop="template.spec.terminationGracePeriodSeconds">
              <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.template.spec.terminationGracePeriodSeconds" />
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
      restart_policy_list: [
        { label: "OnFailure", value: "OnFailure" },
        { label: "Never", value: "Never" },
      ],
      form: {
        activeDeadlineSeconds: null,
        backoffLimit: null,
        completions: null,
        parallelism: null,
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
            restartPolicy: "Always",
          },
        },
      },
    }
  },
  methods: {
    transformation(grandFrom, parentFrom) {
      if (this.form.activeDeadlineSeconds) {
        grandFrom.activeDeadlineSeconds = this.form.activeDeadlineSeconds
      }
      if (this.form.backoffLimit) {
        grandFrom.backoffLimit = this.form.backoffLimit
      }
      if (this.form.completions) {
        grandFrom.completions = this.form.completions
      }
      if (this.form.parallelism) {
        grandFrom.parallelism = this.form.parallelism
      }
      if (this.form.template.spec.terminationGracePeriodSeconds) {
        parentFrom.terminationGracePeriodSeconds = this.form.template.spec.terminationGracePeriodSeconds
      }
      if (this.form.template.spec.restartPolicy) {
        parentFrom.restartPolicy = this.form.template.spec.restartPolicy
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
      if (this.upgradePolicyParentObj.template) {
        if (this.upgradePolicyParentObj.template.spec) {
          if (this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds) {
            this.form.template.spec.terminationGracePeriodSeconds = this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds
          }
          if (this.upgradePolicyParentObj.template.spec.restartPolicy) {
            this.form.template.spec.restartPolicy = this.upgradePolicyParentObj.template.spec.restartPolicy
          }
        }
      }
    }
  },
}
</script>