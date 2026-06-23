<template>
  <div class="login-background">
    <div class="login-container" v-if="mfaPage">
      <div class="login-visual" aria-hidden="true">
        <div class="login-image"></div>
      </div>
      <div class="login-panel login-panel--mfa" v-loading="loading">
        <div class="login-panel-inner">
          <div class="login-heading">
            <div class="login-title">
              {{ systemName }}
            </div>
            <div class="login-border"></div>
          </div>
          <div class="mfa-content" v-if="mfaInit">
            <div class="login-welcome">
              <span>{{ $t("commons.login.mfa_helper") }}</span>
            </div>
            <div class="mfa-qr-wrap">
              <img class="mfa-qr" :src="otp.qrImage">
            </div>
            <div class="mfa-secret">
              <span>Secret:{{ otp.secret }}</span>
            </div>
            <el-form class="login-form login-form--mfa">
              <el-form-item class="login">
                <el-input v-model="mfaCredential.code"></el-input>
              </el-form-item>
              <div v-if="mfaError" class="login-form-error">{{ mfaError }}</div>
            </el-form>
            <div class="login-btn">
              <el-button type="primary" class="submit" @click="bindMfa()" size="default">
                {{ $t("commons.button.bind") }}
              </el-button>
            </div>
          </div>
          <div class="mfa-content" v-else>
            <div class="login-welcome">
              <span>{{ $t("commons.login.mfa_login_helper") }}</span>
            </div>
            <el-form class="login-form login-form--mfa">
              <el-form-item class="login">
                <el-input v-model="mfaCredential.code"></el-input>
              </el-form-item>
              <div v-if="mfaError" class="login-form-error">{{ mfaError }}</div>
            </el-form>
            <div class="login-btn">
              <el-button type="primary" class="submit" @click="mfaLogin()" size="default">
                {{ $t("commons.button.login") }}
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="login-container" v-if="!mfaPage">
      <div class="login-visual" aria-hidden="true">
        <div class="login-image"></div>
      </div>
      <div class="login-panel" v-loading="loading">
        <div class="login-panel-inner">
          <el-form v-if="!forcePasswordPage" :model="form" :rules="rules" ref="form" size="default" class="login-card-form">
            <div class="login-title">
              {{ systemName }}
            </div>
            <div class="login-border"></div>
            <div class="login-welcome">
              {{ $t("commons.login.welcome") }}
            </div>
            <div class="login-form">
              <el-form-item prop="username" class="login">
                <el-input v-model="form.username" :placeholder="$t('commons.login.username_or_email')"
                          prefix-icon="el-icon-user" clearable autofocus/>
              </el-form-item>
              <el-form-item prop="password" class="login">
                <el-input type="password" v-model="form.password" :placeholder="$t('commons.login.password')"
                          prefix-icon="el-icon-lock" show-password maxlength="30"/>
              </el-form-item>
              <div v-if="loginError" class="login-form-error">{{ loginError }}</div>
            </div>
            <div class="login-btn">
              <el-button type="primary" class="submit" @click="submit('form')" size="default">
                {{ $t("commons.button.login") }}
              </el-button>
            </div>
          </el-form>
          <el-form
            v-else
            :model="forcePasswordForm"
            :rules="forcePasswordRules"
            ref="forcePasswordForm"
            size="default"
            class="login-card-form login-card-form--force-password">
            <div class="login-title">
              {{ $t("commons.login.force_change_password_title") }}
            </div>
            <div class="login-border"></div>
            <div class="login-welcome">
              {{ $t("commons.login.force_change_password_tip") }}
            </div>
            <div class="login-form">
              <el-form-item prop="oldPassword" class="login">
                <el-input type="password" show-password v-model="forcePasswordForm.oldPassword"
                          prefix-icon="el-icon-lock"
                          @input="handleForcePasswordInput(['oldPassword', 'newPassword'])"
                          @blur="validateForcePasswordField('oldPassword')"
                          :placeholder="$t('business.user.old_password')"/>
              </el-form-item>
              <el-form-item prop="newPassword" class="login">
                <el-input type="password" show-password v-model="forcePasswordForm.newPassword"
                          prefix-icon="el-icon-lock"
                          @input="handleForcePasswordInput(['newPassword', 'confirmPassword'])"
                          @blur="validateForcePasswordField('newPassword')"
                          :placeholder="$t('business.user.new_password')"/>
              </el-form-item>
              <el-form-item prop="confirmPassword" class="login">
                <el-input type="password" show-password v-model="forcePasswordForm.confirmPassword"
                          prefix-icon="el-icon-lock"
                          @input="handleForcePasswordInput('confirmPassword')"
                          @blur="validateForcePasswordField('confirmPassword')"
                          :placeholder="$t('business.user.confirm_password')"/>
              </el-form-item>
              <div v-if="forcePasswordError" class="login-form-error">{{ forcePasswordError }}</div>
            </div>
            <div class="login-btn force-password-btn">
              <el-button type="primary" class="submit" @click="changeDefaultPassword" size="default">
                {{ $t("commons.button.confirm") }}
              </el-button>
            </div>
          </el-form>
        </div>
      </div>
    </div>
  </div>
</template>
<script>

import { updatePassword } from "@/api/auth"
import {bind, getOtp, valid} from "@/api/mfa"
import Rules from "@/utils/rules"

export default {
  name: "Login",
  data () {
    const validateForcePassword = (rule, value, callback) => {
      if (value === this.forcePasswordForm.oldPassword) {
        callback(new Error(this.$t("commons.login.password_same")))
      } else {
        callback()
      }
    }
    const validateForceConfirmPassword = (rule, value, callback) => {
      if (!value) {
        callback()
        return
      }
      if (value !== this.forcePasswordForm.newPassword) {
        callback(new Error(this.$t("business.user.password_not_equal")))
      } else {
        callback()
      }
    }
    return {
      loading: false,
      form: {
        username: "",
        password: "",
      },
      rules: {
        username: [
          {
            required: true,
            message: this.$tm("commons.validate.input", "commons.login.username_or_email"),
            trigger: "blur"
          },
        ],
        password: [
          {
            required: true,
            message: this.$tm("commons.validate.input", "commons.login.password"),
            trigger: "blur"
          },
        ],
      },
      loginError: "",
      mfaError: "",
      forcePasswordError: "",
      redirect: undefined,
      otherQuery: {},
      systemName: this.$t("commons.login.title"),
      loadingPage: false,
      otp: {
        secret: "",
        qrImage: "",
      },
      mfaPage: false,
      mfaInit: false,
      mfaCredential: {
        userName: "",
        secret: "",
        code: "",
      },
      pendingForceChangePassword: false,
      forcePasswordPage: false,
      forcePasswordForm: {
        oldPassword: "",
        newPassword: "",
        confirmPassword: "",
      },
      forcePasswordRules: {
        oldPassword: [
          Rules.RequiredRule
        ],
        newPassword: [
          Rules.RequiredRule,
          Rules.PasswordRule,
          { validator: validateForcePassword, trigger: "blur" },
        ],
        confirmPassword: [
          {
            required: true,
            trigger: "blur",
            message: this.$t("business.user.please_input_password_again"),
          },
          { validator: validateForceConfirmPassword, trigger: "blur" },
        ]
      },
    }
  },
  watch: {
    $route: {
      handler: function (route) {
        const query = route.query
        if (query) {
          this.redirect = query.redirect
          this.otherQuery = this.getOtherQuery(query)
        }
      },
      immediate: true
    },
    "form.username"() {
      this.loginError = ""
    },
    "form.password"() {
      this.loginError = ""
    },
    "mfaCredential.code"() {
      this.mfaError = ""
    },
    "forcePasswordForm.oldPassword"() {
      this.forcePasswordError = ""
    },
    "forcePasswordForm.newPassword"() {
      this.forcePasswordError = ""
    },
    "forcePasswordForm.confirmPassword"() {
      this.forcePasswordError = ""
    }
  },
  created: function () {
    document.addEventListener("keydown", this.watchEnter)
  },

  destroyed () {
    document.removeEventListener("keydown", this.watchEnter)
  },
  methods: {
    watchEnter (e) {
      let keyCode = e.keyCode
      if (keyCode === 13) {
        if (this.forcePasswordPage) {
          this.changeDefaultPassword()
        } else if(this.mfaPage) {
          this.mfaLogin()
        }else {
          this.submit("form")
        }
      }
    },
    submit (form) {
      this.loginError = ""
      this.$refs[form].validate((valid) => {
        if (valid) {
          this.loading = true
          this.$store.dispatch("user/login", this.form).then((res) => {
            const user = res.data
            this.pendingForceChangePassword = !!user.forceChangePassword
            if (user.mfa.enable) {
              this.mfaPage = true
              this.loading = false
              if (!user.mfa.configured) {
                getOtp().then((res) => {
                  this.otp = res.data
                  this.mfaInit = true
                })
              } else {
                this.mfaInit = false
              }
            } else {
              this.completeLogin(user)
            }
          }).catch((error) => {
            this.loading = false
            this.loginError = this.formatLoginError(this.getErrorMessage(error))
          })
        } else {
          return false
        }
      })
    },
    bindMfa () {
      this.mfaError = ""
      this.mfaCredential.secret = this.otp.secret
      bind(this.mfaCredential, {silentError: true}).then(() => {
        this.mfaPage = false
        location.reload()
      }).catch((error) => {
        this.mfaError = this.getErrorMessage(error)
      }).finally(() => {
      })
    },
    mfaLogin() {
      this.mfaError = ""
      valid(this.mfaCredential, {silentError: true}).then(() => {
        this.completeLogin({ forceChangePassword: this.pendingForceChangePassword })
      }).catch((error) => {
        this.mfaError = this.getErrorMessage(error)
      })
    },
    completeLogin(user) {
      this.loading = false
      if (user.forceChangePassword) {
        this.mfaPage = false
        this.forcePasswordForm.oldPassword = this.form.password
        this.forcePasswordForm.newPassword = ""
        this.forcePasswordForm.confirmPassword = ""
        this.forcePasswordPage = true
        this.$nextTick(() => {
          if (this.$refs.forcePasswordForm) {
            this.$refs.forcePasswordForm.clearValidate()
          }
        })
        return
      }
      this.$router.push({ path: "/" })
    },
    changeDefaultPassword() {
      this.forcePasswordError = ""
      this.$refs.forcePasswordForm.validate((valid) => {
        if (!valid) {
          return
        }
        updatePassword({
          oldPassword: this.forcePasswordForm.oldPassword,
          newPassword: this.forcePasswordForm.newPassword,
        }, {silentError: true}).then(() => {
          this.$message.success(this.$t("commons.msg.update_success"))
          this.forcePasswordPage = false
          this.$router.push({ path: "/" })
        }).catch((error) => {
          this.forcePasswordError = this.getErrorMessage(error)
        })
      })
    },
    handleForcePasswordInput(props) {
      this.forcePasswordError = ""
      this.clearForcePasswordValidate(props)
    },
    clearForcePasswordValidate(props) {
      if (this.$refs.forcePasswordForm) {
        this.$refs.forcePasswordForm.clearValidate(props)
      }
    },
    validateForcePasswordField(prop) {
      if (this.$refs.forcePasswordForm) {
        this.$refs.forcePasswordForm.validateField(prop)
      }
    },
    getErrorMessage(error) {
      if (!error) {
        return ""
      }
      if (typeof error === "string") {
        return error
      }
      if (error.response) {
        const data = error.response.data
        if (data && typeof data === "object") {
          return data.message || data.msg || JSON.stringify(data)
        }
        return data || error.message || ""
      }
      return error.message || ""
    },
    formatLoginError(message) {
      if (!message) {
        return ""
      }
      return message
        .replace(/^登录失败\s*[,，]\s*/, "")
        .replace(/^login failed\s*[,，]\s*/i, "")
    },

    getOtherQuery (query) {
      return Object.keys(query).reduce((acc, cur) => {
        if (cur !== "redirect") {
          acc[cur] = query[cur]
        }
        return acc
      }, {})
    },
  },
}
</script>

<style lang="scss" scoped>
  @import "../../styles/common/variables";

  @mixin login-center {
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .login-background {
    min-height: 100vh;
    height: 100%;
    padding: 32px;
    box-sizing: border-box;
    overflow: auto;
    background:
      linear-gradient(90deg, rgba(37, 99, 235, 0.06) 1px, transparent 1px),
      linear-gradient(180deg, rgba(37, 99, 235, 0.05) 1px, transparent 1px),
      #eef4fb;
    background-size: 42px 42px;
    @include login-center;
  }

  .login-container {
    position: relative;
    display: grid;
    grid-template-columns: minmax(0, 1fr) 440px;
    width: min(1120px, calc(100vw - 64px));
    min-height: 620px;
    overflow: hidden;
    background-color: var(--kp-surface);
    border: 1px solid rgba(203, 213, 225, 0.78);
    border-radius: 8px;
    box-shadow: 0 24px 60px rgba(15, 23, 42, 0.14);

    .login-visual {
      position: relative;
      min-height: 620px;
      overflow: hidden;
      background-color: #0f2f8f;

      &::after {
        content: "";
        position: absolute;
        inset: 0;
        background:
          linear-gradient(90deg, rgba(15, 23, 42, 0.08) 0%, rgba(15, 23, 42, 0.02) 48%, rgba(255, 255, 255, 0.34) 100%),
          linear-gradient(180deg, rgba(255, 255, 255, 0.08) 0%, rgba(15, 23, 42, 0.1) 100%);
        pointer-events: none;
      }
    }

    .login-image {
      position: absolute;
      inset: 0;
      width: 100%;
      height: 100%;
      background: url(../../assets/KubePi-login.jpg) no-repeat;
      background-position: center center;
      background-size: contain;
    }

    .login-panel {
      position: relative;
      z-index: 1;
      display: flex;
      align-items: center;
      min-width: 0;
      background: rgba(255, 255, 255, 0.96);
      backdrop-filter: blur(18px);
      border-left: 1px solid rgba(219, 228, 240, 0.76);
    }

    .login-panel-inner {
      width: 100%;
      padding: 56px 52px;
      box-sizing: border-box;
    }

    .login-card-form {
      width: 100%;
    }

    .login-card-form--force-password {
      & ::v-deep .el-form-item__error {
        position: static;
        width: 100%;
        padding-top: 7px;
        line-height: 18px;
        white-space: normal;
        word-break: break-word;
        overflow-wrap: anywhere;
        box-sizing: border-box;
      }
    }

    .login-title {
      font-size: 34px;
      font-weight: 600;
      letter-spacing: 0;
      line-height: 42px;
      text-align: left;
      color: var(--kp-text);
    }

    .login-border {
      height: 2px;
      margin: 18px 0 18px;
      position: relative;
      width: 64px;
      background: $--color-primary;
      border-radius: 2px;
    }

    .login-welcome {
      max-width: 320px;
      font-size: 14px;
      color: var(--kp-text-muted);
      letter-spacing: 0;
      line-height: 22px;
      text-align: left;
    }

    .login-form {
      margin-top: 34px;

      .el-form-item {
        margin-bottom: 20px;
      }

      & ::v-deep .el-input__inner {
        height: 46px;
        line-height: 46px;
        border-radius: 6px;
        border-color: #d7e1ef;
        background-color: #fbfdff;
        transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
      }

      & ::v-deep .el-input__prefix,
      & ::v-deep .el-input__suffix {
        line-height: 46px;
      }

      & ::v-deep .el-input__inner:hover {
        border-color: #bfdbfe;
      }

      & ::v-deep .el-input__inner:focus {
        border-color: var(--kp-primary);
        background-color: var(--kp-surface);
        box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.12);
      }
    }

    .login-btn {
      margin-top: 32px;

      .submit {
        width: 100%;
        height: 46px;
        border-radius: 6px;
        font-weight: 500;
        box-shadow: 0 12px 24px rgba(37, 99, 235, 0.18);
        transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
      }

      .submit:hover,
      .submit:focus {
        transform: translateY(-1px);
        box-shadow: 0 16px 28px rgba(37, 99, 235, 0.22);
      }
    }

    .login-form-error {
      margin-top: -4px;
      color: $--color-danger;
      font-size: 13px;
      line-height: 20px;
      text-align: left;
    }

    .login-msg {
      margin-top: 10px;
      padding: 0 40px;
      color: $--color-danger;
      text-align: center;
    }

    .submit {
      width: 100%;
      border-radius: 6px;
    }

    .forget-password {
      margin-top: 40px;
      padding: 0 40px;
      float: right;
      @media only screen and (max-width: 1280px) {
        margin-top: 20px;
      }
    }

  }

  .login-panel--mfa {
    .login-form--mfa {
      margin-top: 24px;
    }
  }

  .mfa-content {
    margin-top: 6px;
  }

  .mfa-qr-wrap {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 184px;
    height: 184px;
    margin-top: 24px;
    padding: 12px;
    box-sizing: border-box;
    background: var(--kp-surface);
    border: 1px solid var(--kp-border);
    border-radius: 8px;
    box-shadow: 0 10px 24px rgba(15, 23, 42, 0.08);
  }

  .mfa-qr {
    display: block;
    max-width: 100%;
    max-height: 100%;
  }

  .mfa-secret {
    margin-top: 14px;
    color: var(--kp-text-secondary);
    font-size: 13px;
    line-height: 20px;
    word-break: break-all;
  }

  .force-password-btn {
    margin-top: 26px !important;
  }

  .login ::v-deep input {
    color: var(--kp-text) !important;
  }

  @media only screen and (max-width: 1180px) {
    .login-container {
      grid-template-columns: minmax(0, 1fr) 420px;
      min-height: 560px;

      .login-visual {
        min-height: 560px;
      }

      .login-panel-inner {
        padding: 44px 42px;
      }
    }
  }

  @media only screen and (max-width: 900px) {
    .login-background {
      align-items: flex-start;
      padding: 24px;
    }

    .login-container {
      display: block;
      width: min(460px, calc(100vw - 48px));
      min-height: auto;

      .login-visual {
        display: none;
      }

      .login-panel {
        border-left: 0;
      }

      .login-panel-inner {
        padding: 36px 30px;
      }

      .login-title {
        font-size: 30px;
        line-height: 38px;
      }
    }
  }

  @media only screen and (max-height: 700px) and (min-width: 901px) {
    .login-background {
      padding-top: 24px;
      padding-bottom: 24px;
    }

    .login-container {
      min-height: 540px;

      .login-visual {
        min-height: 540px;
      }

      .login-panel-inner {
        padding-top: 38px;
        padding-bottom: 38px;
      }
    }
  }
</style>
