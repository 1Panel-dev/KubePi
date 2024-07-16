<template>
  <layout-content :header="'SSO'">
    <el-col :span="2"><br/></el-col>
    <el-col :span="15">
      <div class="grid-content bg-purple-light">
        <el-form ref="form" v-loading="loading" label-position="left" :rules="rules" :model="form" label-width="150px">
          <el-form-item style="width: 100%" :label="$t('business.user.sso_protocol')" prop="protocol">
            <el-select v-model="form.protocol"  filterable clearable style="width: 100%" :placeholder="$t('commons.form.select_placeholder')">
                <el-option v-for="item in protocolOptions"
                    :key="item.protocol"
                    :label="item.label"
                    :value="item.value">
                </el-option>
            </el-select>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.sso_interface_address')" prop="interfaceAddress">
            <el-input v-model="form.interfaceAddress"
                      :placeholder="'http://192.168.56.101:30020/realms/kubepi'"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.sso_client_id')" prop="clientId">
            <el-input v-model="form.clientId"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.sso_client_secret')" prop="clientSecret">
            <el-input type="password" v-model="form.clientSecret"></el-input>
          </el-form-item>
          <el-form-item>
            <div style="font-size: 12px;color: #4E5051;">
              {{ $t("business.user.sso_helper") }}
            </div>
          </el-form-item>
          <el-form-item style="width: 100%" prop="enable">
            <el-checkbox v-model="form.enable">{{ $t("commons.enable.true") }}</el-checkbox>
          </el-form-item>
          <div style="float: right">
            <el-form-item>
              <el-button @click="connectTest" :disabled="isSubmitGoing"
                         v-has-permissions="{resource:'ldap',verb:'create'}">{{
                  $t("business.user.ldap_test")
                }}
              </el-button>
              <el-button @click="remake" :disabled="isSubmitGoing"
                         v-has-permissions="{resource:'ldap',verb:'create'}">{{
                  $t("business.user.ldap_remake")
                }}
              </el-button>
              <el-button type="primary" @click="onSubmit" :disabled="isSubmitGoing"
                         v-has-permissions="{resource:'ldap',verb:'create'}">{{ $t("commons.button.confirm") }}
              </el-button>
            </el-form-item>
          </div>

        </el-form>
      </div>
    </el-col>
  </layout-content>
</template>

<script>
import "codemirror/lib/codemirror.css"
import "codemirror/theme/ayu-dark.css"
import "codemirror/mode/javascript/javascript"
import LayoutContent from "@/components/layout/LayoutContent"
import Rule from "@/utils/rules"
import {createSso, getSso, testConnect, updateSso} from "@/api/sso"
export default {
  name: "SSO",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      protocol: "",
      protocolOptions: [
        {
          value: 'openid',
          label: 'OpenID Connect',
        },
        {
          value: 'saml',
          label: 'SAML',
        },
      ],
      interfaceAddress: "",
      clientId: "",
      clientSecret: "",
      form: {
        enable: false,
      },
      loading: false,
      rules: {
        protocol: [Rule.RequiredRule],
        interfaceAddress: [Rule.RequiredRule],
        clientId: [Rule.RequiredRule],
        clientSecret: [Rule.RequiredRule],
      },
      isSubmitGoing: false,
    }
  },
  methods: {
    connectTest () {
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
      this.connectLoading = true
      testConnect(this.form).then(res => {
        this.$message({
          type: "success",
          message: this.$t("business.user.sso_test_result", { count: res.data })
        })
      }).finally(() => {
        this.loading = false
        this.connectLoading = false
      })
    },
    remake () {
      this.form = JSON.parse(JSON.stringify(this.oldConfig))
    },
    onSubmit () {
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
      this.isSubmitGoing = true
      this.loading = true
      if (this.form.uuid === undefined || this.form.uuid === "") {
        createSso(this.form).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success")
          })
          this.list()
        }).finally(() => {
          this.isSubmitGoing = false
          this.loading = false
        })
      } else {
        updateSso(this.form).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success")
          })
          this.list()
        }).finally(() => {
          this.isSubmitGoing = false
          this.loading = false
        })
      }
    },
    list () {
      this.loading = true
      getSso().then((res) => {
        if (res.data.length > 0) {
          this.form = res.data[0]
          this.oldConfig = res.data[0]
        }
      }).finally(() => {
        this.loading = false
      })
    },
    handleSelectionChange(val) {
      this.selects = val
    },
    handleSearch(){
      this.tableUsers =this.users.filter(user => {
        return user.name.indexOf(this.searchName) > -1
      })
    }
  },
  created () {
    this.list()
  }
}
</script>

<style scoped>
    .yaml-editor {
        height: 100%;
        position: relative;
    }
    .yaml-editor >>> .CodeMirror {
        height: auto;
        min-height: 100px;
    }
</style>