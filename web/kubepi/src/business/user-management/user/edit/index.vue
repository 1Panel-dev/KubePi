<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Users' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :rules="rules" :model="form" label-width="150px" label-position="left">

            <el-form-item :label="$t('business.user.username')" prop="name">
              <el-input v-model="form.name" disabled></el-input>
            </el-form-item>

            <el-form-item :label="$t('business.user.nickname')" prop="nickname">
              <el-input v-model="form.nickname" :disabled="form.type ==='LDAP'"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.email')" prop="email">
              <el-input v-model="form.email" :disabled="form.type ==='LDAP'"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.role')" prop="roles">
              <el-select v-model="form.roles" filterable
                         multiple
                         style="width: 100%"
                         :placeholder="$t('commons.form.select_placeholder')">
                <el-option
                        v-for="item in roleOptions "
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                </el-option>
              </el-select>
            </el-form-item>


            <el-form-item :label="$t('business.user.password')" v-if="form.type !=='LDAP'">
              <el-link @click="openedChangePassword">{{ $t("business.user.change_password") }}</el-link>
            </el-form-item>


            <el-form-item>
              <div style="float: right">
                <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" @click="onConfirm()">{{ $t("commons.button.confirm") }}
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
import {updateUser, getUser} from "@/api/users"
import {listRoles} from "@/api/roles"
import Rules from "@/utils/rules"


export default {
  name: "UserEdit",
  props: ["name"],
  components: { LayoutContent },
  data () {
    var validatePass = (rule, value, callback) => {
      if (value === "") {
        callback(new Error(this.$t("business.user.please_input_password")))
      } else {
        if (this.passwordChangeFrom.newPassword !== "") {
          this.$refs.form.validateField("checkPass")
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
      loading: false,
      isSubmitGoing: false,
      user: {},
      roleOptions: [],
      rules: {
        name: [
          Rules.RequiredRule
        ],
        email: [
          Rules.RequiredRule,
          Rules.EmailRule
        ],
        nickname: [
          Rules.RequiredRule,
        ],
        roles: [
          Rules.RequiredRule,
        ],
      },
      changePasswordOpened: false,
      passwordChangeRules: {
        newPassword: [
          Rules.RequiredRule,
          Rules.PasswordRule,
          { validator: validatePass, trigger: "blur" },
        ],
        confirmPassword: [
          Rules.RequiredRule,
          Rules.PasswordRule,
          { validator: validatePass2, trigger: "blur" }
        ],
      },
      passwordChangeFrom: {
        newPassword: "",
        confirmPassword: ""
      },
      form: {
        name: "",
        nickname: "",
        email: "",
        roles: [],
        type: "LOCAL"
      },
    }
  },
  methods: {
    onChangePasswordConfirm () {
      let isFormReady = false
      this.$refs["passwordChangeFrom"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      updateUser(this.name, {
        "password": this.passwordChangeFrom.newPassword
      }).then(() => {
        this.$message.success(this.$t("commons.msg.update_success"))
        this.changePasswordOpened = false
      })
    },
    openedChangePassword () {
      this.changePasswordOpened = true
    },
    onConfirm () {
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

      this.user.name = this.form.name
      this.user.nickName = this.form.nickname
      this.user.email = this.form.email
      this.user.roles = this.form.roles
      updateUser(this.name, this.user).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success")
        })
        this.$router.push({ name: "Users" })
      }).finally(() => {
        this.isSubmitGoing = false
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "Users" })
    },
    onCreated () {
      this.loading = true
      getUser(this.name).then(data => {
        this.form.name = data.data.name
        this.form.nickname = data.data.nickName
        this.form.email = data.data.email
        this.form.roles = data.data.roles
        this.form.type = data.data.type
        this.user = data.data
        listRoles().then(d => {
          d.data.forEach(r => {
            this.roleOptions.push({
              label: r.name,
              value: r.name,
            })
          })
          this.loading = false
        })
      })
    }
  },
  created () {
    this.onCreated()
  }
}
</script>

<style scoped>
</style>
