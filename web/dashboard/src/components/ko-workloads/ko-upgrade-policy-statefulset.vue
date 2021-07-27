<template>
  <div style="margin-top: 20px">
    <ko-card title="Update Policy">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-row>
              <el-form-item label="Strategy" prop="updateStrategy.type">
                <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.updateStrategy.type" :radios="update_strategy_list" />
              </el-form-item>
            </el-row>
            <el-form-item label="Pod Management policy" prop="podManagementPolicy">
              <ko-form-item itemType="radio" v-model="form.podManagementPolicy" :radios="pod_management_policy_list" />
            </el-form-item>
            <el-row>
            </el-row>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Revision History Limit" prop="revisionHistoryLimit">
              <ko-form-item deviderName="Revision" itemType="number" v-model.number="form.revisionHistoryLimit" />
            </el-form-item>
          </el-col>
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
        }
      }
    }
  },
}
</script>