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
              <el-input type="number" v-model.number="form.strategy.rollingUpdate.maxSurge">
                <el-select slot="append" style="width: 80px" v-model="form.strategy.rollingUpdate.maxSurgeUnit">
                  <el-option v-for="(item, index) in devider_list" :key="index" :label="item" :value="item" />
                </el-select>
              </el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Max Unavailable" prop="strategy.rollingUpdate.maxUnavailable">
              <el-input type="number" clearable v-model.number="form.strategy.rollingUpdate.maxUnavailable">
                <el-select slot="append" style="width: 80px" v-model="form.strategy.rollingUpdate.maxUnavailableUnit">
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
            <el-form-item label="Progress Deadline" prop="progressDeadlineSeconds">
              <ko-form-item deviderName="Seconds" itemType="number" v-model.number="form.progressDeadlineSeconds" />
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
      devider_list: ["Pods", "%"],
      form: {
        strategy: {
          type: "Recreate",
          rollingUpdate: {
            maxUnavailable: null,
            maxUnavailableUnit: "%",
            maxSurge: null,
            maxSurgeUnit: "%",
          },
        },
        template: {
          spec: {
            terminationGracePeriodSeconds: null,
          },
        },
        minReadySeconds: null,
        progressDeadlineSeconds: null,
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
          if (this.form.strategy.rollingUpdate.maxSurge) {
            if (this.form.strategy.rollingUpdate.maxSurgeUnit === "%") {
              this.form.strategy.rollingUpdate.maxSurge += "%"
            }
            grandFrom.strategy.rollingUpdate.maxSurge = this.form.strategy.rollingUpdate.maxSurge
          }
          break
      }
      if (this.form.minReadySeconds) {
        grandFrom.minReadySeconds = this.form.minReadySeconds
      }
      if (this.form.revisionHistoryLimit) {
        grandFrom.revisionHistoryLimit = this.form.revisionHistoryLimit
      }
      if (this.form.progressDeadlineSeconds) {
        grandFrom.progressDeadlineSeconds = this.form.progressDeadlineSeconds
      }
      if (this.form.template.spec.terminationGracePeriodSeconds) {
        parentFrom.terminationGracePeriodSeconds = this.form.template.spec.terminationGracePeriodSeconds
      }
    },
  },
  mounted() {
    if (this.upgradePolicyParentObj) {
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
              this.form.strategy.rollingUpdate.maxUnavailableUnit = "Pods"
              this.form.strategy.rollingUpdate.maxUnavailable = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxUnavailable)
            }
          }
          if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge) {
            if (this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge.toString().indexOf("%") !== -1) {
              this.form.strategy.rollingUpdate.maxSurgeUnit = "%"
              this.form.strategy.rollingUpdate.maxSurge = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge.replace("%", ""))
            } else {
              this.form.strategy.rollingUpdate.maxSurgeUnit = "Pods"
              this.form.strategy.rollingUpdate.maxSurge = Number(this.upgradePolicyParentObj.strategy.rollingUpdate.maxSurge)
            }
          }
        }
      }
      if (this.upgradePolicyParentObj.template) {
        if (this.upgradePolicyParentObj.template.spec) {
          if (this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds) {
            this.form.template.spec.terminationGracePeriodSeconds = this.upgradePolicyParentObj.template.spec.terminationGracePeriodSeconds
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
    }
  },
}
</script>