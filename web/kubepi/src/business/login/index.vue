<template xmlns:el-col="http://www.w3.org/1999/html">
  <div class="login-background">
    <div class="login-container" v-if="mfaPage">
      <el-row>
        <el-col>
          <div style="text-align: center" v-if="mfaInit">
            <div>
              <span>{{ $t("commons.login.mfa_helper") }}</span>
            </div>
            <div>
              <img :src="otp.qrImage">
            </div>
            <div>
              <span>Secret:{{ otp.secret }}</span>
            </div>
          </div>
          <div style="width: 50%;margin-left: 25%" v-if="mfaInit">
            <el-form>
              <el-form-item class="login">
                <el-input v-model="mfaCredential.code"></el-input>
              </el-form-item>
            </el-form>
            <div>
              <el-button type="primary" class="submit" @click="bindMfa()" size="default">
                {{ $t("commons.button.bind") }}
              </el-button>
            </div>
          </div>
          <div style="text-align: center;width: 50%;margin-top: 15%;margin-left: 25%" v-if="!mfaInit">
            <div>
              <div>
                <span>{{ $t("commons.login.mfa_login_helper") }}</span>
              </div>
              <div>
                <span>Secret:{{ user.secret }}</span>
              </div>
            </div>
            <el-form>
              <el-form-item class="login">
                <el-input v-model="mfaCredential.code"></el-input>
              </el-form-item>
            </el-form>
            <div>
              <el-button type="primary" class="submit" @click="mfaLogin()" size="default">
                {{ $t("commons.button.login") }}
              </el-button>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>
    <div class="login-container" v-if="!mfaPage">
      <el-row type="flex" v-loading="loading">
        <el-col :span="12">
          <el-form :model="form" :rules="rules" ref="form" size="default">
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
                          autofocus/>
              </el-form-item>
              <el-form-item prop="password" class="login">
                <el-input type="password" v-model="form.password" :placeholder="$t('commons.login.password')"
                          maxlength="30"/>
              </el-form-item>
            </div>
            <div class="login-btn">
              <el-button type="primary" class="submit" @click="submit('form')" size="default">
                {{ $t("commons.button.login") }}
              </el-button>
            </div>
          </el-form>
        </el-col>
        <el-col :span="12">
          <div class="login-image"></div>
        </el-col>
      </el-row>
    </div>
  </div>
</template>
<script>

import {bind, getOtp, valid} from "@/api/mfa"

export default {
  name: "Login",
  data () {
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
      msg: "",
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
      user: {
        secret: ""
      }
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
        this.submit("form")
      }
    },
    submit (form) {
      this.$refs[form].validate((valid) => {
        if (valid) {
          this.loading = true
          this.$store.dispatch("user/login", this.form).then((res) => {
            const user = res.data
            this.user.secret = user.mfa.secret
            if (user.mfa.enable) {
              this.mfaPage = true
              if (user.mfa.secret === "") {
                getOtp().then((res) => {
                  this.otp = res.data
                  this.mfaInit = true
                })
              }
            } else {
              this.$router.push({ path: "/" })
              this.loading = false
            }
          }).catch(() => {
            this.loading = false
          })
        } else {
          return false
        }
      })
    },
    bindMfa () {
      this.mfaCredential.secret = this.otp.secret
      this.mfaCredential.userName = this.form.username
      bind(this.mfaCredential).then(() => {
        this.mfaPage = false
        location.reload()
      }).finally(() => {
      })
    },
    mfaLogin() {
      this.mfaCredential.secret = this.user.secret
      this.mfaCredential.userName = this.form.username
      valid(this.mfaCredential).then(() => {
        this.$router.push({ path: "/" })
      })
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
    background-color: #FFFFFF;
    height: 100%;
    @include login-center;
  }

  .login-container {
    min-width: 640px;
    width: 960px;
    height: 480px;
    background-color: #FFFFFF;
    box-shadow: 0 5px 5px -3px rgba(0, 0, 0, .2), 0 8px 10px 1px rgba(0, 0, 0, .14), 0 3px 14px 2px rgba(0, 0, 0, .12);
    @media only screen and (max-width: 1280px) {
      width: 900px;
      height: 380px;
    }

    .login-logo {
      margin-top: 30px;
      margin-left: 30px;
      @media only screen and (max-width: 1280px) {
        margin-top: 20px;
      }

      img {
        height: 45px;
      }
    }

    .login-title {
      margin-top: 50px;
      font-size: 32px;
      letter-spacing: 0;
      text-align: center;
      color: #999999;

      @media only screen and (max-width: 1280px) {
        margin-top: 20px;
      }
    }

    .login-border {
      height: 2px;
      margin: 20px auto 20px;
      position: relative;
      width: 80px;
      background: $--color-primary;
      @media only screen and (max-width: 1280px) {
        margin: 10px auto 10px;
      }
    }

    .login-welcome {
      margin-top: 10px;
      font-size: 14px;
      color: #999999;
      letter-spacing: 0;
      line-height: 18px;
      text-align: center;
      @media only screen and (max-width: 1280px) {
        margin-top: 20px;
      }
    }

    .login-form {
      margin-top: 30px;
      padding: 0 40px;

      @media only screen and (max-width: 1280px) {
        margin-top: 10px;
      }

      & ::v-deep .el-input__inner {
        border-radius: 0;
      }
    }

    .login-btn {
      margin-top: 40px;
      padding: 0 40px;
      @media only screen and (max-width: 1280px) {
        margin-top: 20px;
      }

      .submit {
        width: 100%;
        border-radius: 0;
      }
    }

    .login-msg {
      margin-top: 10px;
      padding: 0 40px;
      color: $--color-danger;
      text-align: center;
    }

    .login-image {
      background: url(../../assets/KubePi-login.jpg) no-repeat;
      background-size: cover;
      width: 100%;
      height: 480px;
      @media only screen and (max-width: 1280px) {
        height: 380px;
      }
    }

    .submit {
      width: 100%;
      border-radius: 0;
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

  .login /deep/ input {
    background-color: #FFFFFF !important;
    color: #000000 !important;

  }
</style>
