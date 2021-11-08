<template>
  <div>
    <h3>Deployment Processing</h3>
    <div v-loading="loading">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item :label="$t('commons.table.status')" prop="spec.paused">
              <el-radio-group v-model="form.spec.paused" @change="changeStatus">
                <el-radio :label="true">{{$t("business.workload.pause")}}</el-radio>
                <el-radio :label="false">{{$t("business.workload.resume")}}</el-radio>
              </el-radio-group>
              <div><span style="font-size: 12px;" v-if="form.spec.paused">{{$t('business.workload.pause_help')}}</span></div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.min_ready_time')" prop="spec.minReadySeconds">
              <span>{{form.spec.minReadySeconds || 0}}</span>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.progress_deadline')" prop="spec.progressDeadlineSeconds">
              <span>{{form.spec.progressDeadlineSeconds || 0}}</span>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </div>
  </div>
</template>

<script>
import { updateWorkLoad } from "@/api/workloads"

export default {
  name: "KoDetailPause",
  props: {
    cluster: String,
    yamlInfo: Object,
    resourceType: String,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec) {
            this.form = newYamlInfo
            if (newYamlInfo.spec.paused === undefined) {
              this.form.spec.paused = false
            }
            this.loading = false
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      loading: true,
      statu_list: [
        { label: this.$t("business.workload.pause"), value: true },
        { label: this.$t("business.workload.resume"), value: false },
      ],
      form: {
        spec: {
          paused: false,
          minReadySeconds: 0,
          progressDeadlineSeconds: 600,
        },
      },
    }
  },
  methods: {
    initForm() {
      if (this.yamlInfo.spec) {
        this.form = this.yamlInfo
        this.loading = false
      }
    },
    changeStatus() {
      if (this.form.spec.paused) {
        this.$confirm(this.$t("business.workload.suspend") + "  Deployments ?", this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        })
          .then(() => {
            this.updateForm()
          })
          .catch(() => {
            this.form.spec.paused = false
          })
      } else {
        this.updateForm()
      }
    },
    updateForm() {
      this.loading = true
      updateWorkLoad(this.cluster, this.resourceType, this.form.metadata.namespace, this.form.metadata.name, this.form)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.$emit("refresh")
        })
        .finally(() => {
          this.loading = false
        })
    },
  },
  created() {
    this.initForm()
  },
}
</script>

<style scoped>
</style>
