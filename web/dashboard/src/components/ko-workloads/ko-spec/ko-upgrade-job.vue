<template>
  <div>
    <el-form label-position="top" ref="form" :model="form">
      <el-row :gutter="20" v-if="isCronJob()">
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
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.completions')" prop="completions">
            <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.completions" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.parallelism')" prop="parallelism">
            <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.parallelism" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.back_off_limit')" prop="backoffLimit">
            <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.backoffLimit" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.active_deadline')" prop="activeDeadlineSeconds">
            <ko-form-item :deviderName="$t('business.workload.times')" itemType="number" v-model.number="form.activeDeadlineSeconds" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20" v-if="isCronJob()">
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
      <el-row :gutter="20" v-if="isCronJob()">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.starting_deadline_seconds')" prop="startingDeadlineSeconds">
            <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.startingDeadlineSeconds" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoUpgradePolicyCronjob",
  components: { KoFormItem },
  props: {
    upgradePolicyParentObj: Object,
    resourceType: String,
  },
  data() {
    return {
      concurrency_policy_list: [
        { label: this.$t("business.workload.allow_run"), value: "Allow" },
        { label: this.$t("business.workload.skip_if_not_finish"), value: "Forbid" },
        { label: this.$t("business.workload.replace_if_not_finish"), value: "Replace" },
      ],
      suspend_list: [
        { label: this.$t("business.workload.yes"), value: true },
        { label: this.$t("business.workload.no"), value: false },
      ],
      form: {
        activeDeadlineSeconds: null,
        backoffLimit: null,
        completions: null,
        parallelism: null,

        concurrencyPolicy: null,
        successfulJobsHistoryLimit: null,
        failedJobsHistoryLimit: null,
        startingDeadlineSeconds: null,
        suspend: false,
      },
    }
  },
  methods: {
    isCronJob() {
      return this.resourceType === "cronjobs"
    },
    transformation(grandFrom) {
      if (this.isCronJob()) {
        grandFrom.jobTemplate.spec.activeDeadlineSeconds = this.form.activeDeadlineSeconds || undefined
        grandFrom.jobTemplate.spec.backoffLimit = this.form.backoffLimit || undefined
        grandFrom.jobTemplate.spec.completions = this.form.completions || undefined
        grandFrom.jobTemplate.spec.parallelism = this.form.parallelism || undefined

        grandFrom.concurrencyPolicy = this.form.concurrencyPolicy || undefined
        grandFrom.successfulJobsHistoryLimit = this.form.successfulJobsHistoryLimit || undefined
        grandFrom.failedJobsHistoryLimit = this.form.failedJobsHistoryLimit || undefined
        grandFrom.startingDeadlineSeconds = this.form.startingDeadlineSeconds || undefined
        grandFrom.suspend = this.form.suspend || undefined
      } else {
        grandFrom.activeDeadlineSeconds = this.form.activeDeadlineSeconds || undefined
        grandFrom.backoffLimit = this.form.backoffLimit || undefined
        grandFrom.completions = this.form.completions || undefined
        grandFrom.parallelism = this.form.parallelism || undefined
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      let itemSpec = {}
      if (this.isCronJob()) {
        if (this.upgradePolicyParentObj.jobTemplate?.spec) {
          itemSpec = this.upgradePolicyParentObj.jobTemplate.spec
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
      } else {
        itemSpec = this.upgradePolicyParentObj
      }
      if (itemSpec.activeDeadlineSeconds) {
        this.form.activeDeadlineSeconds = itemSpec.activeDeadlineSeconds
      }
      if (itemSpec.backoffLimit) {
        this.form.backoffLimit = itemSpec.backoffLimit
      }
      if (itemSpec.completions) {
        this.form.completions = itemSpec.completions
      }
      if (itemSpec.parallelism) {
        this.form.parallelism = itemSpec.parallelism
      }
    }
  },
}
</script>