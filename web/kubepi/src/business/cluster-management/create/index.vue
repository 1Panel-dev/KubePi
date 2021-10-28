<template>
  <layout-content :header="$t('business.cluster.import')" :back-to="{ name: 'Clusters' }">
    <el-row>
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :rules="rules" :model="form" label-width="150px" label-position="left">

            <el-form-item :label="$t('commons.table.name')" prop="name">
              <el-input v-model="form.name" clearable></el-input>
            </el-form-item>

            <el-divider content-position="center">{{ $t('business.cluster.connect_setting') }}</el-divider>
            <!-- <el-form-item :label="$t('business.cluster.connect_direction')">
              <el-radio v-model="form.direction" label="forward" @change="onDirectionChange">
                {{ $t('business.cluster.connect_forward') }}
              </el-radio>
              <el-radio v-model="form.direction" disabled label="backward" @change="onDirectionChange">
                {{ $t('business.cluster.connect_backward') }}({{ $t("business.cluster.expect") }})
              </el-radio>
            </el-form-item> -->
            <el-form-item :label="$t('business.cluster.authenticate_mode')">
              <el-radio v-model="form.authenticationMode" label="bearer"
                        @change="onAuthenticationModeChange">Bearer Token
              </el-radio>
              <el-radio v-model="form.authenticationMode" label="certificate"
                        @change="onAuthenticationModeChange">{{ $t('business.cluster.certificate') }}
              </el-radio>
              <el-radio v-model="form.authenticationMode" label="configFile"
                        @change="onAuthenticationModeChange">{{ $t('business.cluster.config_file') }}
              </el-radio>
            </el-form-item>
            <div v-if="form.direction==='forward'&&form.authenticationMode!=='configFile'">
              <el-form-item label="API Server" prop="apiServer">
                <el-input v-model="form.apiServer" placeholder="eg: https://127.0.0.1:8443" clearable></el-input>
              </el-form-item>
              <div v-if="!form.apiServerInsecure">
                <el-form-item label="Ca Certificate" prop="caDataStr">
                  <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.caDataStr" clearable></el-input>
                </el-form-item>
              </div>
            </div>
            <div v-if="form.authenticationMode==='bearer'">
              <el-form-item label="Bearer Token" prop="token">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.token" clearable></el-input>
              </el-form-item>
            </div>
            <div v-if="form.authenticationMode==='certificate'">
              <el-form-item label="Certificate" prop="certDataStr">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.certDataStr" clearable></el-input>
              </el-form-item>
              <el-form-item label="Certificate Key" prop="keyDataStr">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea" v-model="form.keyDataStr" clearable></el-input>
              </el-form-item>
            </div>
            <div v-if="form.authenticationMode==='configFile'">
              <el-form-item v-if="form.configContent" :label="$t('business.cluster.config_content')">
                <el-input :autosize="{ minRows: 5, maxRows: 10}" type="textarea"  v-model="form.configContent"></el-input>
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
import {createCluster} from "@/api/clusters"
import Rule from "@/utils/rules"
import {trim} from "@/utils/string"

export default {
  name: "ClusterCreate",
  components: {LayoutContent},
  data() {
    return {
      form: {
        name: "",
        direction: "forward",
        apiServer: "",
        apiServerInsecure: true,
        authenticationMode: "bearer",
        token: "",
        proxyEnable: false,
        proxyServer: "",
        proxyAuthEnable: false,
        proxyUsername: "",
        proxyPassword: "",
        certDataStr: "-----BEGIN CERTIFICATE-----\n" +
            "-----END CERTIFICATE-----",
        keyDataStr: "-----BEGIN RSA PRIVATE KEY-----\n" +
            "-----END RSA PRIVATE KEY-----",
        caDataStr: "",
        configContent: ""
      },
      rules: {
        name: [Rule.RequiredRule, Rule.LengthRule],
        apiServer: [Rule.RequiredRule],
        token: [Rule.RequiredRule],
        certDataStr: [Rule.RequiredRule],
        keyDataStr: [Rule.RequiredRule],
        caDataStr: [Rule.RequiredRule]
      },
      loading: false,
      isSubmitGoing: false,
    }
  },
  methods: {
    readFile(file) {
      const reader = new FileReader()
      reader.readAsText(file)
      reader.onerror = (e) => {
        console.log("error" + e)
      }
      reader.onload = () => {
        this.form.configContent = reader.result
      }
      return false
    },
    onDirectionChange() {
      if (this.form.direction === 'forward') {
        this.form.apiServer = ""
        this.form.authenticationMode = "bearer"
        this.form.certificateData = ""
        this.form.certificateKeyData = ""
        this.form.configContent = ""
        this.form.token = ""
      }
    },
    onAuthenticationModeChange() {
      this.form.certificateData = ""
      this.form.certificateKeyData = ""
      this.form.token = ""
      this.form.configContent = ""
    },
    onCancel() {
      this.$router.push({name: "Clusters"})
    },
    onConfirm() {
      if (this.isSubmitGoing) {
        return
      }

      let isFormReady = false
      this.$refs["form"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      this.loading = true
      this.isSubmitGoing = true
      const req = {
        apiVersion: "v1",
        kind: "Cluster",
        name: trim(this.form.name),
        spec: {
          connect: {
            direction: this.form.direction
          },
          authentication: {
            mode: this.form.authenticationMode
          }

        },
        keyDataStr: "",
        caDataStr: this.form.caDataStr,
        certDataStr: ""
      }
      if (this.form.direction === 'forward') {
        const forwardConfig = {}
        forwardConfig.apiServer = trim(this.form.apiServer)
        if (this.form.proxyEnable) {
          forwardConfig.proxy.username = trim(this.form.proxyUsername)
          forwardConfig.proxy.password = trim(this.form.proxyPassword)
          forwardConfig.proxy.url = trim(this.form.proxyServer)
        }
        req.spec.connect.forward = forwardConfig
      }
      switch (this.form.authenticationMode) {
        case "bearer":
          req.spec.authentication.bearerToken = trim(this.form.token)
          break
        case "certificate":
          req.certDataStr = trim(this.form.certDataStr)
          req.keyDataStr = trim(this.form.keyDataStr)
          break
        case "configFile":
          req.configContentStr = this.form.configContent
      }
      createCluster(req).then(() => {
        this.loading = false
        this.isSubmitGoing = false
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.$router.push({"name": "Clusters"})
      }).finally(() => {
        this.loading = false
        this.isSubmitGoing = false
      })
    }
  },
}
</script>

<style scoped>

</style>