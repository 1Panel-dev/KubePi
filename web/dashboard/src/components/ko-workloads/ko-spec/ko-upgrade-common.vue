<template>
  <div>
    <el-form label-position="top" ref="form" :model="form">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.strategy')" prop="strategy.type">
            <ko-form-item v-if="isDeployment()" itemType="radio" v-model="form.strategy.type" :radios="deployment_strategy_list" />
            <ko-form-item v-else itemType="radio" v-model="form.strategy.type" :radios="strategy_list" />
          </el-form-item>
        </el-col>
        <el-col :span="12" v-if="isStatefulset()">
          <el-form-item :label="$t('business.workload.pod_manage_policy')" prop="podManagementPolicy">
            <ko-form-item itemType="radio" v-model="form.podManagementPolicy" :radios="pod_management_policy_list" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20" v-if="form.strategy.type === 'RollingUpdate'">
        <el-col :span="12" v-if="isDeployment()">
          <el-form-item :label="$t('business.workload.max_surge')" prop="strategy.rollingUpdate.maxSurge">
            <el-input onKeypress="return (/[\d]/.test(String.fromCharCode(event.keyCode)))" v-model.number="form.strategy.rollingUpdate.maxSurge">
              <el-select slot="append" style="width: 80px" v-model="form.strategy.rollingUpdate.maxSurgeUnit">
                <el-option v-for="(item, index) in devider_list" :key="index" :label="item.label" :value="item.value" />
              </el-select>
            </el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12" v-if="isDeployment() || isDaemonset()">
          <el-form-item :label="$t('business.workload.max_unavaliable')" prop="strategy.rollingUpdate.maxUnavailable">
            <el-input onKeypress="return (/[\d]/.test(String.fromCharCode(event.keyCode)))" clearable v-model.number="form.strategy.rollingUpdate.maxUnavailable">
              <el-select slot="append" style="width: 80px" v-model="form.strategy.rollingUpdate.maxUnavailableUnit">
                <el-option v-for="(item, index) in devider_list" :key="index" :label="item.label" :value="item.value" />
              </el-select>
            </el-input>
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12" v-if="isDeployment() || isDaemonset()">
          <el-form-item :label="$t('business.workload.min_ready_time')" prop="minReadySeconds">
            <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.minReadySeconds" />
          </el-form-item>
        </el-col>
        <el-col :span="12" v-if="isDeployment()">
          <el-form-item :label="$t('business.workload.progress_deadline')" prop="progressDeadlineSeconds">
            <ko-form-item :deviderName="$t('business.workload.seconds')" itemType="number" v-model.number="form.progressDeadlineSeconds" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.revision_history_limit')" prop="revisionHistoryLimit">
            <ko-form-item :deviderName="$t('business.workload.revision')" itemType="number" v-model.number="form.revisionHistoryLimit" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoUpgradePolicy",
  components: { KoFormItem },
  props: {
    upgradePolicyParentObj: Object,
    resourceType: String,
  },
  data() {
    return {
      deployment_strategy_list: [
        { label: this.$t("business.workload.rolling_update"), value: "RollingUpdate" },
        { label: this.$t("business.workload.recreate"), value: "Recreate" },
      ],
      strategy_list: [
        { label: this.$t("business.workload.rolling_update"), value: "RollingUpdate" },
        { label: this.$t("commons.button.delete"), value: "OnDelete" },
      ],
      pod_management_policy_list: [
        { label: "OrderedReady", value: "OrderedReady" },
        { label: "Parallel", value: "Parallel" },
      ],
      devider_list: [
        { label: "Pods", value: "" },
        { label: "%", value: "%" },
      ],
      form: {
        strategy: {
          type: "RollingUpdate",
          rollingUpdate: {
            maxUnavailable: null,
            maxUnavailableUnit: "%",
            maxSurge: null,
            maxSurgeUnit: "%",
          },
        },
        minReadySeconds: null,
        progressDeadlineSeconds: null,
        revisionHistoryLimit: null,
      },
    }
  },
  methods: {
    isDeployment() {
      return this.resourceType === "deployments"
    },
    isStatefulset() {
      return this.resourceType === "statefulsets"
    },
    isDaemonset() {
      return this.resourceType === "daemonsets"
    },
    transformation(grandFrom) {
      let itemForm = {}
      itemForm.type = this.form.strategy.type || undefined
      switch (this.resourceType) {
        case "deployments":
          {
            if (this.form.strategy.type === "RollingUpdate") {
              itemForm.rollingUpdate = {}
              if (this.form.strategy.rollingUpdate.maxUnavailableUnit === "%") {
                itemForm.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable ? this.form.strategy.rollingUpdate.maxUnavailable + "%" : undefined
              } else {
                itemForm.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable ? this.form.strategy.rollingUpdate.maxUnavailable : undefined
              }
              if (this.form.strategy.rollingUpdate.maxSurgeUnit === "%") {
                itemForm.rollingUpdate.maxSurge = this.form.strategy.rollingUpdate.maxSurge ? this.form.strategy.rollingUpdate.maxSurge + "%" : undefined
              } else {
                itemForm.rollingUpdate.maxSurge = this.form.strategy.rollingUpdate.maxSurge ? this.form.strategy.rollingUpdate.maxSurge : undefined
              }
            }
          }
          grandFrom.minReadySeconds = this.form.minReadySeconds || undefined
          grandFrom.revisionHistoryLimit = this.form.revisionHistoryLimit || undefined
          grandFrom.progressDeadlineSeconds = this.form.progressDeadlineSeconds || undefined
          grandFrom.strategy = itemForm
          break
        case "statefulsets": {
          grandFrom.podManagementPolicy = this.form.podManagementPolicy || undefined
          grandFrom.revisionHistoryLimit = this.form.revisionHistoryLimit || undefined
          grandFrom.updateStrategy = itemForm
          break
        }
        case "daemonsets": {
          if (this.form.strategy.type === "RollingUpdate") {
            itemForm.rollingUpdate = {}
            if (this.form.strategy.rollingUpdate.maxUnavailableUnit === "%") {
              itemForm.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable ? this.form.strategy.rollingUpdate.maxUnavailable + "%" : undefined
            } else {
              itemForm.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable ? this.form.strategy.rollingUpdate.maxUnavailable : undefined
            }
          }
          grandFrom.minReadySeconds = this.form.minReadySeconds || undefined
          grandFrom.revisionHistoryLimit = this.form.revisionHistoryLimit || undefined
          grandFrom.updateStrategy = itemForm
          break
        }
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      switch (this.resourceType) {
        case "deployments": {
          if (this.upgradePolicyParentObj.strategy) {
            if (this.upgradePolicyParentObj.strategy.type) {
              this.form.strategy.type = this.upgradePolicyParentObj.strategy.type
            }
            if (this.upgradePolicyParentObj.strategy.rollingUpdate) {
              if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable) {
                if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable.toString().indexOf("%") !== -1) {
                  this.form.strategy.rollingUpdate.maxUnavailableUnit = "%"
                  this.form.strategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable.replace("%", ""))
                } else {
                  this.form.strategy.rollingUpdate.maxUnavailableUnit = ""
                  this.form.strategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable)
                }
              }
              if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge) {
                if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge.toString().indexOf("%") !== -1) {
                  this.form.strategy.rollingUpdate.maxSurgeUnit = "%"
                  this.form.strategy.rollingUpdate.maxSurge = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge.replace("%", ""))
                } else {
                  this.form.strategy.rollingUpdate.maxSurgeUnit = ""
                  this.form.strategy.rollingUpdate.maxSurge = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge)
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
          if (this.upgradePolicyParentObj.progressDeadlineSeconds) {
            this.form.progressDeadlineSeconds = this.upgradePolicyParentObj.progressDeadlineSeconds
          }
          break
        }
        case "statefulsets": {
          if (this.upgradePolicyParentObj.updateStrategy) {
            if (this.upgradePolicyParentObj.updateStrategy.type) {
              this.form.strategy.type = this.upgradePolicyParentObj.updateStrategy.type
            }
          }
          if (this.upgradePolicyParentObj.podManagementPolicy) {
            this.form.podManagementPolicy = this.upgradePolicyParentObj.podManagementPolicy
          }
          if (this.upgradePolicyParentObj.revisionHistoryLimit) {
            this.form.revisionHistoryLimit = this.upgradePolicyParentObj.revisionHistoryLimit
          }
          break
        }
        case "daemonsets": {
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
                  this.form.strategy.rollingUpdate.maxUnavailableUnit = ""
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
          break
        }
      }
    }
  },
}
</script>