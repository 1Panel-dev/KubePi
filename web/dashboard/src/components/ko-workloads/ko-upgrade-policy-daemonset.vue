<template>
  <div style="margin-top: 20px">
    <ko-card title="Update Policy">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Strategy" prop="updateStrategy.type">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.updateStrategy.type" :radios="update_strategy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Max Unavailable" prop="updateStrategy.rollingUpdate.maxUnavailable">
              <el-input type="number" v-model.number="form.updateStrategy.rollingUpdate.maxUnavailable">
                <el-select slot="append" style="width: 80px" v-model="form.updateStrategy.rollingUpdate.maxUnavailableUnit">
                  <el-option v-for="(item, index) in devider_list" :key="index" :label="item" :value="item" />
                </el-select>
              </el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Minimum Ready Time" prop="minReadySeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model.number="form.minReadySeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Revision History Limit" prop="revisionHistoryLimit">
              <ko-form-item deviderName="Revision" itemType="number" v-model.number="form.revisionHistoryLimit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Termination Grace Period" prop="template.spec.terminationGracePeriodSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model.number="form.template.spec.terminationGracePeriodSeconds" />
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
      update_strategy_list: [
        { label: "Rolling Update: Create new pods, until max surge is reached, before deleting old pods. Don't stop more pods than max unavailable.", value: "RollingUpdate" },
        { label: "On Delete: New pods are only created when old pods are manually deleted.", value: "Recreate" },
      ],
      devider_list: ["Pods", "%"],
      form: {
        updateStrategy: {
          type: "Recreate",
          rollingUpdate: {
            maxUnavailable: null,
            maxUnavailableUnit: "%",
          },
        },
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
          },
        },
        minReadySeconds: null,
        revisionHistoryLimit: null,
      },
    }
  },
  methods: {
    transformation(grandFrom, parentFrom) {
      grandFrom.updateStrategy = {}
      switch (this.form.updateStrategy.type) {
        case "Recreate":
          grandFrom.updateStrategy.type = "Recreate"
          break
        case "RollingUpdate":
          grandFrom.updateStrategy.type = "RollingUpdate"
          grandFrom.updateStrategy.rollingUpdate = {}
          if (this.form.updateStrategy.rollingUpdate.maxUnavailable) {
            if (this.form.updateStrategy.rollingUpdate.maxUnavailableUnit === "%") {
              this.form.updateStrategy.rollingUpdate.maxUnavailable += "%"
            }
            grandFrom.updateStrategy.rollingUpdate.maxUnavailable = this.form.updateStrategy.rollingUpdate.maxUnavailable
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
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      if (this.upgradePolicyParentObj.updateStrategy) {
        if (this.upgradePolicyParentObj.updateStrategy.type) {
          this.form.updateStrategy.type = this.upgradePolicyParentObj.updateStrategy.type
        }
        if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate) {
          if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable) {
            if (this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable.toString().indexOf("%") !== -1) {
              this.form.updateStrategy.rollingUpdate.maxUnavailableUnit = "%"
              this.form.updateStrategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable.replace("%", ""))
            } else {
              this.form.updateStrategy.rollingUpdate.maxUnavailableUnit = "Pods"
              this.form.updateStrategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.updateStrategy.rollingUpdate.maxUnavailable)
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
        }
      }
    }
  },
}
</script>