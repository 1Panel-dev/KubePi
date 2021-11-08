<template>
  <div>
    <h3 style="display: inline-block;">{{ $t("business.workload.deploy_strategy") }}</h3>
    <el-button v-if="!enableEdit" style="margin-left: 10px;" type="text" icon="el-icon-edit" @click="enableEdit = true">{{$t('commons.button.click_to_edit')}}</el-button>

    <div v-loading="loading">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.strategy')" prop="spec.strategy.type" required>
              <span style="margin-left: 10px;" v-if="!enableEdit">{{getLabels(form.spec.strategy.type)}}</span>
              <el-select v-else style="width:100%" v-model="form.spec.strategy.type">
                <el-option v-for="(item, index) in deployment_strategy_list" :key="index" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20" v-if="form.spec.strategy.type === 'RollingUpdate'">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.max_surge')" prop="spec.strategy.rollingUpdate.maxSurge" :rules="numberPersentRule">
              <span style="margin-left: 10px;" v-if="!enableEdit">{{form.spec.strategy.rollingUpdate.maxSurge}}</span>
              <el-input v-else v-model="form.spec.strategy.rollingUpdate.maxSurge" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.max_unavaliable')" prop="spec.strategy.rollingUpdate.maxUnavailable" :rules="numberPersentRule">
              <span style="margin-left: 10px;" v-if="!enableEdit">{{form.spec.strategy.rollingUpdate.maxUnavailable}}</span>
              <el-input v-else v-model="form.spec.strategy.rollingUpdate.maxUnavailable" />
            </el-form-item>
          </el-col>
        </el-row>
        <div style="float:right" v-if="enableEdit">
          <el-button @click="updateForm">{{$t('commons.button.submit')}}</el-button>
          <el-button @click="enableEdit = false">{{$t('commons.button.cancel')}}</el-button>
        </div>
        <br>
      </el-form>
    </div>
  </div>
</template>

<script>
import { updateWorkLoad } from "@/api/workloads"
import Rule from "@/utils/rules"

export default {
  name: "KoDetailUpdate",
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
      deployment_strategy_list: [
        { label: this.$t("business.workload.rolling_update"), value: "RollingUpdate" },
        { label: this.$t("business.workload.recreate"), value: "Recreate" },
      ],
      strategy_list: [
        { label: this.$t("business.workload.rolling_update"), value: "RollingUpdate" },
        { label: this.$t("commons.button.delete"), value: "OnDelete" },
      ],
      form: {
        spec: {
          strategy: {
            maxSurge: 0,
            maxUnavailable: 0,
          },
          type: "",
        },
      },
      enableEdit: false,
      numberPersentRule: Rule.NumberPersentRule,
    }
  },
  methods: {
    initForm() {
      if (this.yamlInfo.spec) {
        this.form = this.yamlInfo
        this.loading = false
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
          this.enableEdit = false
          this.$emit("refresh")
        })
        .finally(() => {
          this.loading = false
        })
    },
    getLabels(value) {
      switch (value) {
        case "RollingUpdate":
          return this.$t("business.workload.rolling_update")
        case "Recreate":
          return this.$t("business.workload.recreate")
        case "OnDelete":
          return this.$t("commons.button.delete")
      }
    },
  },
  created() {
    this.initForm()
  },
}
</script>

<style scoped>
</style>
