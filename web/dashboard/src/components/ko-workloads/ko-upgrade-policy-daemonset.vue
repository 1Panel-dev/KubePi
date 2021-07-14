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
              <ko-form-item deviderName="Seconds" itemType="input" v-model="form.minReadySeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Revision History Limit" prop="revisionHistoryLimit">
              <ko-form-item deviderName="Revision" itemType="input" v-model="form.revisionHistoryLimit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Termination Grace Period" prop="template.spec.terminationGracePeriodSeconds">
              <ko-form-item deviderName="Seconds" itemType="input" v-model="form.template.spec.terminationGracePeriodSeconds" />
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
            maxUnavailable: "0",
            maxUnavailableUnit: "%",
          },
        },
        template: {
          spec: {
            terminationGracePeriodSeconds: 30,
          },
        },
        minReadySeconds: 0,
        revisionHistoryLimit: 10,
      },
    }
  },
  methods: {
    transformation(parentFrom) {
      parentFrom.updateStrategy = {}
      switch (this.form.updateStrategy.type) {
        case "Recreate":
          parentFrom.updateStrategy.type = "Recreate"
          break
        case "RollingUpdate":
          parentFrom.updateStrategy.type = "RollingUpdate"
          parentFrom.updateStrategy.rollingUpdate = {}
          if (this.form.updateStrategy.rollingUpdate.maxUnavailable) {
            if (this.form.updateStrategy.rollingUpdate.maxUnavailableUnit === "%") {
              this.form.updateStrategy.rollingUpdate.maxUnavailable += "%"
            }
            parentFrom.updateStrategy.rollingUpdate.maxUnavailable = this.form.updateStrategy.rollingUpdate.maxUnavailable
          }
          break
      }
      if (this.form.minReadySeconds) {
        parentFrom.minReadySeconds = this.form.minReadySeconds
      }
      if (this.form.revisionHistoryLimit) {
        parentFrom.revisionHistoryLimit = this.form.revisionHistoryLimit
      }
      if (this.form.template.spec.terminationGracePeriodSeconds) {
        parentFrom.template = {}
        parentFrom.template.spec = {}
        parentFrom.template.spec.terminationGracePeriodSeconds = this.form.template.spec.terminationGracePeriodSeconds
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