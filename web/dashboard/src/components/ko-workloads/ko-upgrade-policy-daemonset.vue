<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.upgrade_policy')">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.strategy')" prop="strategy.type">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.strategy.type" :radios="strategy_list" />
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
            <el-form-item :label="$t('business.workload.max_unavaliable')" prop="strategy.rollingUpdate.maxUnavailable">
              <el-input type="number" v-model.number="form.strategy.rollingUpdate.maxUnavailable">
                <el-select slot="append" style="width: 80px" v-model="form.strategy.rollingUpdate.maxUnavailableUnit">
                  <el-option v-for="(item, index) in devider_list" :key="index" :label="item" :value="item" />
                </el-select>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.min_ready_time')" prop="minReadySeconds">
              <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.minReadySeconds" />
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
      devider_list: ["Pods", "%"],
      form: {
        strategy: {
          type: "Recreate",
          rollingUpdate: {
            maxUnavailable: null,
            maxUnavailableUnit: "%",
          },
        },
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
            restartPolicy: "Always",
          },
        },
        minReadySeconds: null,
        revisionHistoryLimit: null,
      },
    }
  },
  methods: {
    transformation(grandFrom, parentFrom) {
      grandFrom.strategy = {}
      switch (this.form.strategy.type) {
        case "Recreate":
          grandFrom.strategy.type = "Recreate"
          break
        case "RollingUpdate":
          grandFrom.strategy.type = "RollingUpdate"
          grandFrom.strategy.rollingUpdate = {}
          if (this.form.strategy.rollingUpdate.maxUnavailable) {
            if (this.form.strategy.rollingUpdate.maxUnavailableUnit === "%") {
              this.form.strategy.rollingUpdate.maxUnavailable += "%"
            }
            grandFrom.strategy.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable
          }
          break
      }
      if (this.form.minReadySeconds) {
        grandFrom.minReadySeconds = this.form.minReadySeconds
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
          this.form.strategy.type = this.upgradePolicyParentObj.updateStrategy.type
        }
        if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate) {
          if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable) {
            if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable.toString().indexOf("%") !== -1) {
              this.form.strategy.rollingUpdate.maxUnavailableUnit = "%"
              this.form.strategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable.replace("%", ""))
            } else {
              this.form.strategy.rollingUpdate.maxUnavailableUnit = "Pods"
              this.form.strategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable)
            }
          }
        }
      }
      if (this.upgradePolicyParentObj.minReadySeconds) {
        this.form.minReadySeconds = this.upgradePolicyParentObj.minReadySeconds
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