<template>
  <div v-loading="loading">
    <h3>{{$t('business.workload.general')}}</h3>
    <el-form label-position="top">
      <el-row :gutter="20">
        <el-col v-if="form.spec.restartPolicy" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.restart_strategy')">
            <span>{{form.spec.restartPolicy}}</span>
          </el-form-item>
        </el-col>
        <el-col v-if="form.spec.terminationGracePeriodSeconds" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.termination_grace_period')">
            <span>{{form.spec.terminationGracePeriodSeconds}}</span>
          </el-form-item>
        </el-col>

        <el-col v-if="form.spec.dnsPolicy" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.dns_policy')">
            <span>{{form.spec.dnsPolicy}}</span>
          </el-form-item>
        </el-col>
        <el-col v-if="hostNetwork" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.host_network')">
            <span>{{hostNetwork}}</span>
          </el-form-item>
        </el-col>

        <el-col v-if="form.spec.serviceAccount" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.service_account')">
            <div class="spanInFormStyle"><span :title="form.spec.serviceAccount">{{form.spec.serviceAccount}}</span></div>
          </el-form-item>
        </el-col>
        <el-col v-if="form.spec.serviceAccountName" :span="12">
          <el-form-item style="margin-left: 20px;" :label="$t('business.workload.service_account_name')">
            <div class="spanInFormStyle"><span :title="form.spec.serviceAccountName">{{form.spec.serviceAccountName}}</span></div>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
    <br>
  </div>
</template>

<script>
export default {
  name: "KoDetailGeneral",
  props: {
    yamlInfo: Object,
  },
  watch: {
    yamlInfo: {
      handler(newYamlInfo) {
        if (newYamlInfo) {
          if (newYamlInfo.spec) {
            this.form = newYamlInfo
            if (newYamlInfo.spec.hostNetwork == undefined) {
              this.hostNetwork = this.$t("business.workload.no")
            } else {
              this.hostNetwork = this.hostNetwork ? this.$t("business.workload.yes") : this.$t("business.workload.no")
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
      form: {
        spec: {},
      },
      hostNetwork: "",
    }
  },
}
</script>

<style scoped>
</style>
