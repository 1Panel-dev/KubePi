<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.upgrade_policy')">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.concurrency')" prop="concurrencyPolicy">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.concurrencyPolicy" :radios="concurrency_policy_list" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.workload.suspend')" prop="suspend">
              <ko-form-item itemType="radio" v-model="form.suspend" :radios="suspend_list" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.workload.restart_strategy')" prop="jobTemplate.spec.template.spec.restartPolicy">
              <ko-form-item itemType="radio" v-model="form.jobTemplate.spec.template.spec.restartPolicy" :radios="restart_policy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.completions')" prop="jobTemplate.spec.completions">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.jobTemplate.spec.completions" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.parallelism')" prop="jobTemplate.spec.parallelism">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.jobTemplate.spec.parallelism" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.back_off_limit')" prop="jobTemplate.spec.backoffLimit">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.jobTemplate.spec.backoffLimit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.active_deadline')" prop="jobTemplate.spec.activeDeadlineSeconds">
              <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.jobTemplate.spec.activeDeadlineSeconds" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.successful_job_history_limit')" prop="successfulJobsHistoryLimit">
              <ko-form-item itemType="number" v-model.number="form.successfulJobsHistoryLimit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.failed_job_history_limit')" prop="failedJobsHistoryLimit">
              <ko-form-item itemType="number" v-model.number="form.failedJobsHistoryLimit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.allow_run')" prop="startingDeadlineSeconds">
              <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.startingDeadlineSeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.allow_run')" prop="jobTemplate.spec.template.spec.terminationGracePeriodSeconds">
              <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds" />
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
  name: "KoUpgradePolicyCronjob",
  components: { KoFormItem, KoCard },
  props: {
    upgradePolicyParentObj: Object,
  },
  data() {
    return {
      concurrency_policy_list: [
        { label: this.$t("business.workload.allow_run"), value: "Allow" },
        { label: this.$t("business.workload.skip_if_not_finish"), value: "Forbid" },
        { label: this.$t("business.workload.replace_if_not_finish"), value: "Replace" },
      ],
      restart_policy_list: [
        { label: "OnFailure", value: "OnFailure" },
        { label: "Never", value: "Never" },
      ],
      suspend_list: [
        { label: this.$t("business.workload.yes"), value: true },
        { label: this.$t("business.workload.no"), value: false },
      ],
      form: {
        jobTemplate: {
          spec: {
            activeDeadlineSeconds: null,
            backoffLimit: null,
            completions: null,
            parallelism: null,
            template: {
              spec: {
                terminationGracePeriodSeconds: null,
                restartPolicy: "OnFailure",
              },
            },
          },
        },
        concurrencyPolicy: null,
        successfulJobsHistoryLimit: null,
        failedJobsHistoryLimit: null,
        startingDeadlineSeconds: null,
        suspend: false,
      },
    }
  },
  methods: {
    transformation(grandFrom, parentFrom) {
      if (this.form.jobTemplate.spec.activeDeadlineSeconds) {
        grandFrom.jobTemplate.spec.activeDeadlineSeconds = this.form.jobTemplate.spec.activeDeadlineSeconds
      }
      if (this.form.jobTemplate.spec.backoffLimit) {
        grandFrom.jobTemplate.spec.backoffLimit = this.form.jobTemplate.spec.backoffLimit
      }
      if (this.form.jobTemplate.spec.completions) {
        grandFrom.jobTemplate.spec.completions = this.form.jobTemplate.spec.completions
      }
      if (this.form.jobTemplate.spec.parallelism) {
        grandFrom.jobTemplate.spec.parallelism = this.form.jobTemplate.spec.parallelism
      }
      if (this.form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds) {
        parentFrom.terminationGracePeriodSeconds = this.form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds
      }
      if (this.form.jobTemplate.spec.template.spec.restartPolicy) {
        parentFrom.restartPolicy = this.form.jobTemplate.spec.template.spec.restartPolicy
      }

      if (this.form.concurrencyPolicy) {
        grandFrom.concurrencyPolicy = this.form.concurrencyPolicy
      }
      if (this.form.successfulJobsHistoryLimit) {
        grandFrom.successfulJobsHistoryLimit = this.form.successfulJobsHistoryLimit
      }
      if (this.form.failedJobsHistoryLimit) {
        grandFrom.failedJobsHistoryLimit = this.form.failedJobsHistoryLimit
      }
      if (this.form.startingDeadlineSeconds) {
        grandFrom.startingDeadlineSeconds = this.form.startingDeadlineSeconds
      }
      if (this.form.suspend !== undefined) {
        grandFrom.suspend = this.form.suspend
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      if (this.upgradePolicyParentObj.jobTemplate) {
        if (this.upgradePolicyParentObj.jobTemplate.spec) {
          if (this.upgradePolicyParentObj.jobTemplate.spec.activeDeadlineSeconds) {
            this.form.jobTemplate.spec.activeDeadlineSeconds = this.upgradePolicyParentObj.jobTemplate.spec.activeDeadlineSeconds
          }
          if (this.upgradePolicyParentObj.jobTemplate.spec.backoffLimit) {
            this.form.jobTemplate.spec.backoffLimit = this.upgradePolicyParentObj.jobTemplate.spec.backoffLimit
          }
          if (this.upgradePolicyParentObj.jobTemplate.spec.completions) {
            this.form.jobTemplate.spec.completions = this.upgradePolicyParentObj.jobTemplate.spec.completions
          }
          if (this.upgradePolicyParentObj.jobTemplate.spec.parallelism) {
            this.form.jobTemplate.spec.parallelism = this.upgradePolicyParentObj.jobTemplate.spec.parallelism
          }
          if (this.upgradePolicyParentObj.jobTemplate.spec.template) {
            if (this.upgradePolicyParentObj.jobTemplate.spec.template.spec) {
              if (this.upgradePolicyParentObj.jobTemplate.spec.template.spec.terminationGracePeriodSeconds) {
                this.form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds = this.upgradePolicyParentObj.jobTemplate.spec.template.spec.terminationGracePeriodSeconds
              }
              if (this.upgradePolicyParentObj.jobTemplate.spec.template.spec.restartPolicy) {
                this.form.jobTemplate.spec.template.spec.restartPolicy = this.upgradePolicyParentObj.jobTemplate.spec.template.spec.restartPolicy
              }
            }
          }
        }
      }
      if (this.upgradePolicyParentObj.concurrencyPolicy) {
        this.form.concurrencyPolicy = this.upgradePolicyParentObj.concurrencyPolicy
      }
      if (this.upgradePolicyParentObj.successfulJobsHistoryLimit) {
        this.form.successfulJobsHistoryLimit = this.upgradePolicyParentObj.successfulJobsHistoryLimit
      }
      if (this.upgradePolicyParentObj.failedJobsHistoryLimit) {
        this.form.failedJobsHistoryLimit = this.upgradePolicyParentObj.failedJobsHistoryLimit
      }
      if (this.upgradePolicyParentObj.startingDeadlineSeconds) {
        this.form.startingDeadlineSeconds = this.upgradePolicyParentObj.startingDeadlineSeconds
      }
      if (this.upgradePolicyParentObj.suspend !== undefined) {
        this.form.suspend = this.upgradePolicyParentObj.suspend
      }
    }
  },
}
</script>