<template>
  <div style="margin-top: 20px">
    <ko-card title="Update Policy">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Strategy" prop="strategy.type">
              <ko-form-item radioLayout="vertical" itemType="radio" v-model="form.strategy.type" :radios="strategy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="form.strategy.type === 'RollingUpdate'">
          <el-col :span="12">
            <el-form-item label="Max Surge" prop="strategy.rollingUpdate.maxSurge">
              <ko-form-item deviderName="%" itemType="input" v-model="form.strategy.rollingUpdate.maxSurge" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Max Unavailable" prop="strategy.rollingUpdate.maxUnavailable">
              <ko-form-item deviderName="%" itemType="input" v-model="form.strategy.rollingUpdate.maxUnavailable" />
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
            <el-form-item label="Progress Deadline" prop="progressDeadlineSeconds">
              <ko-form-item deviderName="Seconds" itemType="input" v-model="form.progressDeadlineSeconds" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Termination Grace Period" prop="terminationGracePeriodSeconds">
              <ko-form-item deviderName="Seconds" itemType="input" v-model="form.terminationGracePeriodSeconds" />
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
  name: "KoUpgradePolicy",
  components: { KoFormItem, KoCard },
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