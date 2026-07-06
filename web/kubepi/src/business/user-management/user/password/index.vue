<template>
    <layout-content :header="headerTitle" :back-to="backTo">
      <el-row v-loading="loading">
        <el-col :span="4"><br/></el-col>
        <el-col :span="10">
          <div class="grid-content bg-purple-light">
            <el-form :rules="passwordChangeRules" :validate-on-rule-change="false" ref="passwordChangeFrom" :model="passwordChangeFrom" label-width="150px" label-position="left">
                <el-form-item v-if="isCurrentUser" :label="$t('business.user.old_password')" prop="oldPassword">
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
  import {getCurrentUser, updatePassword} from "@/api/auth"
  import {resetPassword} from "@/api/users"
  import Rules from "@/utils/rules"
  
  export default {
    name: "UserPassword",
    props: ["name"],
    components: { LayoutContent },
    data () {
      return {
        loading: false,
        currentUserName: "",
        targetName: this.name || "",
        passwordChangeFrom: {
            oldPassword: "",
            newPassword: "",
            confirmPassword: ""
        }
      }
    },
    computed: {
      isCurrentUser() {
        return this.targetName !== "" && this.targetName === this.currentUserName
      },
      headerTitle() {
        return this.isCurrentUser ? this.$t("business.user.change_password") : this.$t("business.user.reset_password")
      },
      backTo() {
        return this.isCurrentUser ? null : { name: "Users" }
      },
      passwordChangeRules() {
        const rules = {
          newPassword: [
            Rules.RequiredRule,
            Rules.PasswordRule,
            {validator: this.validateNewPassword, trigger: 'blur'},
          ],
          confirmPassword: [
            Rules.RequiredRule,
            Rules.PasswordRule,
            {validator: this.validateConfirmPassword, trigger: 'blur'}
          ]
        }
        if (this.isCurrentUser) {
          rules.oldPassword = [
            Rules.RequiredRule
          ]
        }
        return rules
      }
    },
    methods: {
      validateNewPassword(rule, value, callback) {
        if (value === "") {
          callback(new Error(this.$t("business.user.please_input_password")))
        } else {
          if (this.passwordChangeFrom.confirmPassword !== "") {
            this.$refs.passwordChangeFrom.validateField("confirmPassword")
          }
          callback()
        }
      },
      validateConfirmPassword(rule, value, callback) {
        if (value === "") {
          callback(new Error(this.$t("business.user.please_input_password")))
        } else if (value !== this.passwordChangeFrom.newPassword) {
          callback(new Error(this.$t("business.user.password_not_equal")))
        } else {
          callback()
        }
      },
      clearPasswordValidate() {
        if (this.$refs.passwordChangeFrom) {
          this.$refs.passwordChangeFrom.clearValidate()
        }
      },
      onChangePasswordConfirm () {
        this.$refs["passwordChangeFrom"].validate((valid) => {
          if (valid) {
            const request = this.isCurrentUser
              ? updatePassword({
                  "oldPassword": this.passwordChangeFrom.oldPassword,
                  "newPassword": this.passwordChangeFrom.newPassword
                })
              : resetPassword(this.targetName, this.passwordChangeFrom.newPassword)
            request.then(() => {
              this.$message.success(this.$t("commons.msg.update_success"))
              this.onCancel()
            })
          }
        })
      },
      onCancel () {
        if (this.isCurrentUser) {
          this.$router.push({ path: "/" })
          return
        }
        this.$router.push({ name: "Users" })
      }
    },
    created() {
      this.loading = true
      getCurrentUser().then(data => {
        this.currentUserName = data.data.name
        if (!this.targetName) {
          this.targetName = this.currentUserName
        }
      }).finally(() => {
        this.loading = false
        this.$nextTick(() => {
          this.clearPasswordValidate()
        })
      })
    }
  }
  </script>

  <style scoped>
  </style>
