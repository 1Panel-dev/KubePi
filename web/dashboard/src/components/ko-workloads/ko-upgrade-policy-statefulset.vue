<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.upgrade_policy')">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col>
            <el-form-item :label="$t('business.workload.strategy')" prop="updateStrategy.type">
              <ko-form-item itemType="radio" v-model="form.updateStrategy.type" :radios="strategy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.pod_manage_policy')" prop="podManagementPolicy">
              <ko-form-item itemType="radio" v-model="form.podManagementPolicy" :radios="pod_management_policy_list" />
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
            <el-form-item :label="$t('business.workload.revision_history_limit')" prop="revisionHistoryLimit">
              <ko-form-item :deviderName="$t('business.workload.revision')" itemType="number" v-model.number="form.revisionHistoryLimit" />
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
  name: "KoUpgradePolicyDaemonset",
  components: { KoFormItem, KoCard },
  props: {
    upgradePolicyParentObj: Object,
  },
  data() {
    return {
      strategy_list: [
        { label: this.$t("business.workload.rolling_update"), value: "RollingUpdate" },
        { label: this.$t("business.workload.recreate"), value: "Recreate" },
      ],
      restart_policy_list: [
        { label: "Always", value: "Always" },
        { label: "OnFailure", value: "OnFailure" },
        { label: "Never", value: "Never" },
      ],
      pod_management_policy_list: [
        { label: "OrderedReady", value: "OrderedReady" },
        { label: "Parallel", value: "Parallel" },
      ],
      form: {
        updateStrategy: {
          type: "Recreate",
        },
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
            restartPolicy: "Always",
          },
        },
        podManagementPolicy: "OrderedReady",
        revisionHistoryLimit: null,
      },
    }
  },
  methods: {
    transformation(grandFrom, parentFrom) {
      grandFrom.updateStrategy = {}
      grandFrom.updateStrategy.type = this.form.updateStrategy.type
      if (this.form.podManagementPolicy) {
        grandFrom.podManagementPolicy = this.form.podManagementPolicy
      }
      if (this.form.revisionHistoryLimit) {
        grandFrom.revisionHistoryLimit = this.form.revisionHistoryLimit
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
      if (this.upgradePolicyParentObj.updateStrategy) {
        if (this.upgradePolicyParentObj.updateStrategy.type) {
          this.form.updateStrategy.type = this.upgradePolicyParentObj.updateStrategy.type
        }
      }
      if (this.upgradePolicyParentObj.podManagementPolicy) {
        this.form.podManagementPolicy = this.upgradePolicyParentObj.podManagementPolicy
      }
      if (this.upgradePolicyParentObj.revisionHistoryLimit) {
        this.form.revisionHistoryLimit = this.upgradePolicyParentObj.revisionHistoryLimit
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