<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Users' }">
    <el-row v-loading="loading">
      <el-col :span="4"><br/></el-col>
      <el-col :span="10">
        <div class="grid-content bg-purple-light">
          <el-form ref="form" :model="form" :rules="rules" label-width="150px" label-position="left">

            <el-form-item :label="$t('business.user.username')" prop="name">
              <el-input v-model="form.name"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.nickname')" prop="nickname">
              <el-input v-model="form.nickname"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.email')" prop="email">
              <el-input v-model="form.email"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.password')" prop="password">
              <el-input type="password" v-model="form.password"></el-input>
            </el-form-item>


            <el-form-item :label="$t('business.user.confirm_password')" prop="confirmPassword">
              <el-input type="password" v-model="form.confirmPassword"></el-input>
            </el-form-item>

            <el-form-item :label="$t('business.user.role')" prop="roles">
              <el-select v-model="form.roles" multiple filterable
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
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {createUser} from "@/api/users"
import {listRoles} from "@/api/roles"
import Rules from "@/utils/rules"

export default {
  name: "UserCreate",
  components: {LayoutContent},
  data() {
    var validatePass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('business.user.please_input_password')));
      } else {
        if (this.form.password !== '') {
          this.$refs.form.validateField('checkPass');
        }
        callback();
      }
    };
    var validatePass2 = (rule, value, callback) => {
      if (value === '') {
        callback(new Error(this.$t('business.user.please_input_password')));
      } else if (value !== this.form.password) {
        callback(new Error(this.$t('business.user.password_not_equal')));
      } else {
        callback();
      }
    };
    return {
      loading: false,
      isSubmitGoing: false,
      roleOptions: [],
      rules: {
        name: [
          Rules.RequiredRule,
          Rules.CommonNameRule
        ],
        nickname: [
          Rules.RequiredRule
        ],
        email: [
          Rules.RequiredRule,
          Rules.EmailRule
        ],
        password: [
          Rules.PasswordRule,
          {validator: validatePass, trigger: 'blur'},
        ],
        confirmPassword: [
          Rules.PasswordRule,
          {validator: validatePass2, trigger: 'blur'}
        ],
        roles: [
          Rules.RequiredRule,
        ],
      },
      form: {
        name: "",
        nickname: "",
        email: "",
        password: "",
        confirmPassword: "",
        roles: [],
      },
    }
  },
  methods: {

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
      this.isSubmitGoing = true
      this.loading = true
      const req = {
        "apiVersion": "v1",
        "kind": "User",
        "name": this.form.name,
        "roles": this.form.roles,
        "nickName": this.form.nickname,
        "email": this.form.email,
        "authenticate": {
          "password": this.form.confirmPassword
        }
      }
      createUser(req).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.$router.push({name: "Users"})
      }).finally(
          () => {
            this.isSubmitGoing = false
            this.loading = false
          }
      )
    },

    onCancel() {
      this.$router.push({name: "Users"})
    },
  },
  created() {
    this.loading = true
    listRoles().then(d => {
      d.data.forEach(r => {
        this.roleOptions.push({
          label: r.name,
          value: r.name,
        })
      })
      this.loading = false
    })
  }
}
</script>

<style scoped>
</style>
