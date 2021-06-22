<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Strategy" clearable itemType="radio" v-model="form.strategy.type" :radios="strategy_list" />
      </el-col>
      <el-col :span="12">
        <div v-if="form.strategy.type === 'RollingUpdate'">
          <el-row>
            <ko-form-item labelName="Max Surge" deviderName="Seconds" tipContant="The maximum number of pods allowed beyond the desired scale at any given time." clearable itemType="input" v-model="form.strategy.rollingUpdate.maxSurge" />
          </el-row>
          <el-row style="margin-top: 20px">
            <ko-form-item labelName="Max Unavailable" deviderName="Seconds" tipContant="The maximum number of pods which can be unavailable at any given time." clearable itemType="input" v-model="form.strategy.rollingUpdate.maxUnavailable" />
          </el-row>
        </div>
        <el-row v-if="form.strategy.type === 'RollingUpdateStartStop' || form.strategy.type === 'RollingUpdateStopStart'">
          <ko-form-item labelName="Batch Size" deviderName="Pod" tipContant="Pods will be started and stopped this many at a time" clearable itemType="input" v-model="form.revisionHistoryLimit" />
        </el-row>
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Minimum Ready Time" deviderName="Seconds" tipContant="Containers in th pods must be up for at least this long before the pod is considered available." clearable itemType="input" v-model="form.minReadySeconds" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Progress Deadline" deviderName="Seconds" tipContant="How long to wait without seeing progress before marking the deploument as stalled." clearable itemType="input" v-model="form.progressDeadlineSeconds" />
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
        { label: "Rolling: start new pods, then stop old", value: "RollingUpdateStartStop" },
        { label: "Rolling: stop old pods, then start new", value: "RollingUpdateStopStart" },
        { label: "Kill ALL pods, then start new", value: "Recreate" },
        { label: "Custom", value: "RollingUpdate" },
      ],
      form: {
        strategy: {
          type: "Recreate",
          rollingUpdate: {
            maxUnavailable: "",
            maxSurge: "",
          },
        },
        revisionHistoryLimit: 1,
        minReadySeconds: 0,
        progressDeadlineSeconds: 600,
      },
    }
  },
  methods: {
    transformation(parentFrom) {
      switch (this.form.strategy.type) {
        case "RollingUpdateStartStop":
          parentFrom.strategy.type = "RollingUpdate"
          if (this.form.revisionHistoryLimit) {
            parentFrom.revisionHistoryLimit = this.form.revisionHistoryLimit
          }
          break
        case "RollingUpdateStopStart":
          parentFrom.strategy.type = "RollingUpdate"
          if (this.form.revisionHistoryLimit) {
            parentFrom.revisionHistoryLimit = this.form.revisionHistoryLimit
          }
          break
        case "Recreate":
          parentFrom.strategy.type = "Recreate"
          break
        case "RollingUpdate":
          parentFrom.strategy.type = "RollingUpdate"
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