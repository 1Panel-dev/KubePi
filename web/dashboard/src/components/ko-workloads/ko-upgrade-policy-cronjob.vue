<template>
  <div style="margin-top: 20px">
    <ko-card title="Update Policy">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Concurrency" prop="concurrencyPolicy">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.concurrencyPolicy" :radios="concurrency_policy_list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Suspend" prop="suspend">
              <ko-form-item itemType="radio" v-model="form.suspend" :radios="suspend_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Completions" prop="jobTemplate.spec.completions">
              <ko-form-item deviderName="Times" itemType="number" v-model.number="form.jobTemplate.spec.completions" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Parallelism" prop="jobTemplate.spec.parallelism">
              <ko-form-item deviderName="Times" itemType="number" v-model.number="form.jobTemplate.spec.parallelism" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Back Off Limit" prop="jobTemplate.spec.backoffLimit">
              <ko-form-item deviderName="Times" itemType="number" v-model.number="form.jobTemplate.spec.backoffLimit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Active Deadline" prop="jobTemplate.spec.activeDeadlineSeconds">
              <ko-form-item deviderName="Times" itemType="number" v-model.number="form.jobTemplate.spec.activeDeadlineSeconds" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Successful Job History Limit" prop="successfulJobsHistoryLimit">
              <ko-form-item itemType="number" v-model.number="form.successfulJobsHistoryLimit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Failed Job History Limit" prop="failedJobsHistoryLimit">
              <ko-form-item itemType="number" v-model.number="form.failedJobsHistoryLimit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Starting Deadline Seconds" prop="startingDeadlineSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model.number="form.startingDeadlineSeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Pod Active Deadline" prop="jobTemplate.spec.template.spec.terminationGracePeriodSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model.number="form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds" />
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
        { label: "Allow CronJobs to run Concurrently", value: "Allow" },
        { label: "Skip next if current run hasn't finished", value: "Forbid" },
        { label: "Replace run if current run hasn't finished", value: "Replace" },
      ],
      suspend_list: [
        { label: "Yes", value: true },
        { label: "No", value: false },
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
                terminationGracePeriodSeconds: 30,
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
    transformation(parentFrom) {
      if (this.form.jobTemplate.spec.activeDeadlineSeconds) {
        parentFrom.jobTemplate.spec.activeDeadlineSeconds = this.form.jobTemplate.spec.activeDeadlineSeconds
      }
      if (this.form.jobTemplate.spec.backoffLimit) {
        parentFrom.jobTemplate.spec.backoffLimit = this.form.jobTemplate.spec.backoffLimit
      }
      if (this.form.jobTemplate.spec.completions) {
        parentFrom.jobTemplate.spec.completions = this.form.jobTemplate.spec.completions
      }
      if (this.form.jobTemplate.spec.parallelism) {
        parentFrom.jobTemplate.spec.parallelism = this.form.jobTemplate.spec.parallelism
      }
      if (this.form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds) {
        parentFrom.jobTemplate.spec.template.spec.terminationGracePeriodSeconds = this.form.jobTemplate.spec.template.spec.terminationGracePeriodSeconds
      }

      if (this.form.concurrencyPolicy) {
        parentFrom.concurrencyPolicy = this.form.concurrencyPolicy
      }
      if (this.form.successfulJobsHistoryLimit) {
        parentFrom.successfulJobsHistoryLimit = this.form.successfulJobsHistoryLimit
      }
      if (this.form.failedJobsHistoryLimit) {
        parentFrom.failedJobsHistoryLimit = this.form.failedJobsHistoryLimit
      }
      if (this.form.startingDeadlineSeconds) {
        parentFrom.startingDeadlineSeconds = this.form.startingDeadlineSeconds
      }
      if (this.form.suspend !== undefined) {
        parentFrom.suspend = this.form.suspend
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