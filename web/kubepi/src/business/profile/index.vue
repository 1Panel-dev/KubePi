<template>
  <layout-content :header="$t('commons.personal.profile')">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="profileForm" :rules="profileRules" :model="profileForm" label-width="150px"
                   label-position="left">

            <el-form-item :label="$t('business.user.username')" prop="name">
              {{ profileForm.name }}
            </el-form-item>

            <el-form-item :label="$t('business.user.nickname')" prop="nickname">
              <span v-if="!editModeOpened">{{ profileForm.nickname }}</span>
              <span v-if="editModeOpened"><el-input v-model="profileForm.nickname"></el-input></span>
            </el-form-item>

            <el-form-item :label="$t('business.user.email')" prop="email">
              <span v-if="!editModeOpened">{{ profileForm.email }}</span>
              <span v-if="editModeOpened"><el-input v-model="profileForm.email"></el-input></span>
            </el-form-item>

            <el-form-item :label="$t('business.user.password')">
              <el-link @click="openPasswordChange">{{ $t('business.user.change_password') }}</el-link>
            </el-form-item>

            <el-form-item>
              <div style="float: right">
                <el-button v-if="!editModeOpened" @click="onEditProfile">{{ $t("commons.button.edit") }}</el-button>
                <el-button v-if="editModeOpened" @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button v-if="editModeOpened" type="primary" @click="onConfirm()">{{
                    $t("commons.button.confirm")
                  }}
                </el-button>
              </div>
            </el-form-item>
          </el-form>

        </div>
      </el-col>
      <el-col :span="4"><br/></el-col>
    </el-row>

    <el-dialog
        :title="$t('business.user.change_password')"
        :visible.sync="changePasswordOpened"
        :close-on-click-modal="false"
        width="30%">
      <div>
        <el-form :rules="passwordChangeRules" ref="passwordChangeFrom" :model="passwordChangeFrom" label-width="150px"
                 label-position="left">
          <el-form-item :label="$t('business.user.old_password')" prop="oldPassword">
            <el-input type="password" v-model="passwordChangeFrom.oldPassword"></el-input>
          </el-form-item>
          <el-form-item :label="$t('business.user.new_password')" prop="newPassword">
            <el-input type="password" v-model="passwordChangeFrom.newPassword"></el-input>
          </el-form-item>
          <el-form-item :label="$t('business.user.confirm_password')" prop="confirmPassword">
            <el-input type="password" v-model="passwordChangeFrom.confirmPassword"></el-input>
          </el-form-item>
        </el-form>
      </div>
      <span slot="footer" class="dialog-footer">
    <el-button @click="changePasswordOpened = false">{{ $t("commons.button.cancel") }}</el-button>
    <el-button type="primary" @click="onChangePasswordConfirm">{{ $t("commons.button.confirm") }}</el-button>
  </span>
    </el-dialog>


  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getCurrentUser, updateProfile, updatePassword,} from "@/api/auth"
import Rules from "@/utils/rules"

export default {
  name: "UserProfile",
  components: {LayoutContent},
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('business.user.please_input_password')));
      } else {
        if (this.passwordChangeFrom.newPassword !== '') {
          this.$refs.passwordChangeFrom.validateField('checkPass');
        }
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('business.user.please_input_password')));
      } else if (value !== this.passwordChangeFrom.newPassword) {
        callback(new Error(this.$t('business.user.password_not_equal')));
      } else {
        callback();
      }
    };
    return {
      loading: false,
      isSubmitGoing: false,
      roleOptions: [],
      editModeOpened: false,
      changePasswordOpened: false,
      profileForm: {
        name: "",
        nickname: "",
        email: "",
      },
      profileRules: {
        nickname: [
          Rules.RequiredRule
        ],
        email: [
          Rules.RequiredRule,
          Rules.EmailRule
        ],
      },
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

    onConfirm() {
      if (this.isSubmitGoing) {
        return
      }
      let isFormReady = false
      this.$refs["profileForm"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      this.isSubmitGoing = true
      this.loading = true
      updateProfile({"nickName": this.profileForm.nickname, "email": this.profileForm.email}).then(() => {
        this.isSubmitGoing = false
        this.loading = false
        this.editModeOpened = false
        this.$message.success(this.$t('commons.msg.update_success'))
        // 修改vuex
        window.location.reload()
      })

    },
    onEditProfile() {
      this.editModeOpened = true
    },
    onCancel() {
      this.editModeOpened = false
      this.onVueCreated()
    },

    openPasswordChange() {
      this.passwordChangeFrom.newPassword = ""
      this.passwordChangeFrom.oldPassword = ""
      this.changePasswordOpened = true
    },
    onChangePasswordConfirm() {
      let isFormReady = false
      this.$refs["passwordChangeFrom"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      updatePassword({
        "oldPassword": this.passwordChangeFrom.oldPassword,
        "newPassword": this.passwordChangeFrom.newPassword
      }).then(() => {
        this.$message.success(this.$t('commons.msg.update_success'))
        this.changePasswordOpened = false
      })
    },
    onVueCreated() {
      getCurrentUser().then(data => {
        this.profileForm.name = data.data.name
        this.profileForm.nickname = data.data.nickName
        this.profileForm.email = data.data.email
      })
    }
  },


  created() {
    this.onVueCreated()
  }
}
</script>

<style scoped>
</style>
