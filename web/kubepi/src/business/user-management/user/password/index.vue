<template>
    <layout-content :header="$t('business.user.change_password')" :back-to="{ name: 'Users' }">
      <el-row>
        <el-col :span="4"><br/></el-col>
        <el-col :span="10">
          <div class="grid-content bg-purple-light">
            <el-form :rules="passwordChangeRules" ref="passwordChangeFrom" :model="passwordChangeFrom" label-width="150px" label-position="left">
                <el-form-item :label="$t('business.user.old_password')" prop="oldPassword">
                    <el-input type="password" show-password v-model="passwordChangeFrom.oldPassword"></el-input>
                </el-form-item>

                <el-form-item :label="$t('business.user.new_password')" prop="newPassword">
                    <el-input type="password" show-password v-model="passwordChangeFrom.newPassword"></el-input>
                </el-form-item>

                <el-form-item :label="$t('business.user.confirm_password')" prop="confirmPassword">
                    <el-input type="password" show-password v-model="passwordChangeFrom.confirmPassword"></el-input>
                </el-form-item>
                <el-form-item>
                    <div style="float: right">
                        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                        <el-button type="primary" @click="onChangePasswordConfirm">{{ $t("commons.button.confirm") }}
                        </el-button>
                    </div>
                </el-form-item>
          </el-form>
          </div>
        </el-col>
        <el-col :span="4"><br/></el-col>
      </el-row>
  
    </layout-content>
  </template>
  
  <script>
  import LayoutContent from "@/components/layout/LayoutContent"
  import {getCurrentUser} from "@/api/auth"
  import {updateUser} from "@/api/users"
  import Rules from "@/utils/rules"
  
  export default {
    name: "UserPassword",
    props: ["name"],
    components: { LayoutContent },
    data () {
      var validatePass = (rule, value, callback) => {
        if (value === "") {
          callback(new Error(this.$t("business.user.please_input_password")))
        } else {
          if (this.passwordChangeFrom.newPassword !== "") {
            this.$refs.passwordChangeFrom.validateField("checkPass")
          }
          callback()
        }
      }
      var validatePass2 = (rule, value, callback) => {
        if (value === "") {
          callback(new Error(this.$t("business.user.please_input_password")))
        } else if (value !== this.passwordChangeFrom.newPassword) {
          callback(new Error(this.$t("business.user.password_not_equal")))
        } else {
          callback()
        }
      }
      return {
        passwordChangeFrom: {
            oldPassword: "",
            newPassword: "",
            confirmPassword: ""
        },
        passwordChangeRules: {
            oldPassword: [
                Rules.RequiredRule
            ],
            newPassword: [
                Rules.RequiredRule,
                Rules.PasswordRule,
                {validator: validatePass, trigger: 'blur'},
            ],
            confirmPassword: [
                Rules.RequiredRule,
                Rules.PasswordRule,
                {validator: validatePass2, trigger: 'blur'}
            ]
        }
      }
    },
    methods: {
      onChangePasswordConfirm () {
        this.$refs["passwordChangeFrom"].validate((valid) => {
          if (valid) {
            updateUser(this.name, {
                "oldPassword": this.passwordChangeFrom.oldPassword,
                "password": this.passwordChangeFrom.newPassword
            }).then(() => {
                this.$message.success(this.$t("commons.msg.update_success"))
                this.onCancel()
            })
          }
        })
      },
      onCancel () {
        this.$router.push({ name: "Users" })
      }
    },
    created() {
      if (!this.name) {
        getCurrentUser().then(data => {
          this.name = data.data.name
        })
      }
    }
  }
  </script>
  
  <style scoped>
  </style>
  