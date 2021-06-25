<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Strategy" radioLayout="vertical" clearable itemType="radio" v-model="form.strategy.type" :radios="strategy_list" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px" v-if="form.strategy.type === 'RollingUpdate'">
      <el-col :span="12">
        <ko-form-item labelName="Max Surge" deviderName="%" tipContant="The maximum number of pods allowed beyond the desired scale at any given time." clearable itemType="input" v-model="form.strategy.rollingUpdate.maxSurge" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Max Unavailable" deviderName="%" tipContant="The maximum number of pods which can be unavailable at any given time." clearable itemType="input" v-model="form.strategy.rollingUpdate.maxUnavailable" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Minimum Ready Time" deviderName="Seconds" tipContant="Containers in th pods must be up for at least this long before the pod is considered available." clearable itemType="input" v-model="form.minReadySeconds" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Revision History Limit" deviderName="Revision" tipContant="The number of old ReplicaSets to retain for rollback." clearable itemType="input" v-model="form.revisionHistoryLimit" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Progress Deadline" deviderName="Seconds" tipContant="How long to wait without seeing progress before marking the deploument as stalled." clearable itemType="input" v-model="form.progressDeadlineSeconds" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Termination Grace Period" deviderName="Seconds" tipContant="The duration the pod needs to terminate successfully." clearable itemType="input" v-model="form.terminationGracePeriodSeconds" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoUpgradePolicy",
  components: { KoFormItem },
  props: {
    upgradePolicyParentObj: Object,
  },
  data() {
    return {
      strategy_list: [
        { label: "Rolling Updatate: Create new pods, then stop old", value: "RollingUpdate" },
        { label: "Recreate: Kill ALL pods, then start new", value: "Recreate" },
      ],
      form: {
        strategy: {
          type: "Recreate",
          rollingUpdate: {
            maxUnavailable: "25%",
            maxSurge: "25%",
          },
        },
        minReadySeconds: 0,
        progressDeadlineSeconds: 600,
        revisionHistoryLimit: 10,
        terminationGracePeriodSeconds: 30,
      },
    }
  },
  methods: {
    transformation(parentFrom) {
      parentFrom.strategy = {}
      switch (this.form.strategy.type) {
        case "Recreate":
          parentFrom.strategy.type = "Recreate"
          break
        case "RollingUpdate":
          parentFrom.strategy.type = "RollingUpdate"
          parentFrom.strategy.rollingUpdate = {}
          if (this.form.strategy.rollingUpdate.maxUnavailable) {
            parentFrom.strategy.rollingUpdate.maxUnavailable = this.form.strategy.rollingUpdate.maxUnavailable
          }
          if (this.form.strategy.rollingUpdate.maxSurge) {
            parentFrom.strategy.rollingUpdate.maxSurge = this.form.strategy.rollingUpdate.maxSurge
          }
          break
      }
      if (this.form.minReadySeconds) {
        parentFrom.minReadySeconds = this.form.minReadySeconds
      }
      if (this.form.progressDeadlineSeconds) {
        parentFrom.progressDeadlineSeconds = this.form.progressDeadlineSeconds
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
      if (this.upgradePolicyParentObj.strategy) {
        if (this.upgradePolicyParentObj.strategy.type) {
          if (this.upgradePolicyParentObj.strategy.type === "Recreate") {
            this.form.strategy.type = "Recreate"
          } else {
            this.form.strategy.type = "RollingUpdate"
          }
        }
        if (this.upgradePolicyParentObj.strategy.rollingUpdate) {
          if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable) {
            this.form.strategy.rollingUpdate.maxUnavailable = this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable
          }
          if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge) {
            this.form.strategy.rollingUpdate.maxSurge = this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge
          }
        }
      }
      if (this.upgradePolicyParentObj.minReadySeconds) {
        this.form.minReadySeconds = this.upgradePolicyParentObj.minReadySeconds
      }
      if (this.upgradePolicyParentObj.progressDeadlineSeconds) {
        this.form.progressDeadlineSeconds = this.upgradePolicyParentObj.progressDeadlineSeconds
      }
    }
  },
}
</script>