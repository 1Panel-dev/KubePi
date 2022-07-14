<template>
  <layout-content :header="$t('business.cluster.edit_cluster')" :back-to="{ name: 'Clusters' }">
    <el-row>
      <el-col :span="4"><br /></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :rules="rules" :model="form" label-width="150px" label-position="left">
            <el-form-item :label="$t('commons.table.name')" prop="name">
              {{ form.name }}
            </el-form-item>

            <el-divider content-position="center">{{ $t('business.cluster.connect_setting') }}</el-divider>
            <el-form-item :label="$t('business.cluster.authenticate_mode')">
              <el-radio v-model="form.mode" label="bearer" @change="onAuthenticationModeChange">Bearer Token
              </el-radio>
              <el-radio v-model="form.mode" label="certificate" @change="onAuthenticationModeChange">{{ $t('business.cluster.certificate') }}
              </el-radio>
              <el-radio v-model="form.mode" label="configFile" @change="onAuthenticationModeChange">{{ $t('business.cluster.config_file') }}
              </el-radio>
            </el-form-item>
            <div v-if="form.mode!=='configFile'">
              <el-form-item label="API Server" prop="apiServer">
                <el-input v-model="form.apiServer" placeholder="eg: https://127.0.0.1:8443" clearable></el-input>
              </el-form-item>
            </div>
            <div v-if="form.mode==='bearer'">
              <el-form-item label="Bearer Token" prop="token">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.token" clearable></el-input>
              </el-form-item>
            </div>
            <div v-if="form.mode==='certificate'">
              <el-form-item label="Certificate" prop="certData">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.certData" clearable></el-input>
              </el-form-item>
              <el-form-item label="Certificate Key" prop="keyData">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.keyData" clearable></el-input>
              </el-form-item>
            </div>
            <div v-if="form.mode==='configFile'">
              <el-form-item :label="$t('business.cluster.config_content')">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.configFileContent"></el-input>
              </el-form-item>
              <el-form-item>
                <el-upload :before-upload="readFile" action="" style="display: inline-block;margin-left: 5px">
                  <el-button>{{ $t('commons.button.upload') }}</el-button>
                </el-upload>
              </el-form-item>
            </div>
            <el-form-item>
              <div style="float: right">
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" @click="onConfirm()" :disabled="isSubmitGoing">{{
                    $t("commons.button.confirm")
                  }}
                </el-button>
              </div>
            </el-form-item>
          </el-form>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import { getCluster, updateCluster } from "@/api/clusters"
import Rule from "@/utils/rules"

export default {
  name: "ClusterEdit",
  props: ["name"],
  components: { LayoutContent },
  data() {
    return {
      form: {
        name: "",
        apiServer: "",
        mode: "bearer",
        token: "",
        certData: "",
        keyData: "",
        configFileContent: "",
        withLabel: false,
      },
      rules: {
        apiServer: [Rule.RequiredRule],
        token: [Rule.RequiredRule],
        certData: [Rule.RequiredRule],
        keyData: [Rule.RequiredRule],
        configFileContent: [Rule.RequiredRule],
      },
      loading: false,
      isSubmitGoing: false,
    }
  },
  methods: {
    readFile(file) {
      const reader = new FileReader()
      reader.readAsText(file)
      reader.onload = () => {
        this.form.configFileContent = reader.result
      }
      return false
    },
    onAuthenticationModeChange() {
      this.form.certificateData = ""
      this.form.certificateKeyData = ""
      this.form.token = ""
      this.form.configFileContent = ""
    },
    onCancel() {
      this.$router.push({ name: "Clusters" })
    },
    loadClusterInfo() {
      getCluster(this.name).then((res) => {
        this.form.name = res.data.name
        this.form.mode = res.data.spec.authentication.mode || ""
        this.form.apiServer = res.data.spec.connect.forward.apiServer || ""
      })
    },
    onConfirm() {
      if (this.isSubmitGoing) {
        return
      }
      this.$refs["form"].validate((valid) => {
        if (valid) {
          this.isSubmitGoing = true
          updateCluster(this.name, this.form)
            .then(() => {
              this.loading = false
              this.isSubmitGoing = false
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success"),
              })
              this.$router.push({ name: "Clusters" })
            })
            .finally(() => {
              this.loading = false
              this.isSubmitGoing = false
            })
        }
      })
    },
  },
  mounted() {
    this.loadClusterInfo()
  },
}
</script>

<style scoped>
</style>
