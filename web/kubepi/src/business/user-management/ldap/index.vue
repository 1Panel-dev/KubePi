<template>
  <layout-content :header="'LDAP'">
    <el-col :span="2"><br/></el-col>
    <el-col :span="15">
      <div class="grid-content bg-purple-light">
        <el-form ref="form" v-loading="loading" label-position="left" :rules="rules" :model="form" label-width="150px">
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_address')" prop="address">
            <el-input v-model="form.address" :placeholder="'172.16.10.1'"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_port')" prop="port">
            <el-input v-model="form.port" :placeholder="'389'" type="number"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_tls')" prop="tls">
            <el-checkbox v-model="form.tls">{{ $t("commons.bool.true") }}</el-checkbox>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_username')" prop="username">
            <el-input v-model="form.username"
                      :placeholder="'cn=Manager,ou=users,dc=ko,dc=com or  Manager@ko.com'"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_password')" prop="password">
            <el-input type="password" v-model="form.password"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_filter_dn')" prop="dn">
            <el-input v-model="form.dn" :placeholder="'ou=users,dc=ko,dc=com'"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_filter_rule')" prop="filter">
            <el-input v-model="form.filter" :placeholder="'(&(objectClass=organizationalPerson))'"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_mapping')" prop="mapping">
            <codemirror ref="editor" class="yaml-editor" v-model="form.mapping" :options="options"></codemirror>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.time_limit')" prop="timeLimit">
            <el-input v-model.number="form.timeLimit" :placeholder="'30'" type="number"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.size_limit')" prop="sizeLimit">
            <el-input v-model.number="form.sizeLimit" :placeholder="'1000'" type="number"></el-input>
          </el-form-item>
          <el-form-item>
            <div style="font-size: 12px;color: #4E5051;">
              {{ $t("business.user.ldap_mapping_helper") }}
              <br>
              {{ $t("business.user.ldap_helper") }}
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
              <el-button @click="openLoginPage" :disabled="isSubmitGoing"
                         v-has-permissions="{resource:'ldap',verb:'create'}">{{
                  $t("business.user.test_login")
                }}
              </el-button>
              <el-button @click="openImportPage" :disabled="isSubmitGoing"
                         v-has-permissions="{resource:'ldap',verb:'create'}">{{
                  $t("business.user.import_user")
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
    <el-col :span="4"><br/></el-col>
    <el-dialog :visible.sync="loginPageOpen" :title="$t('business.user.test_login')" width="30%">
      <el-form ref="loginForm" v-loading="loginLoading" label-position="left" :rules="rules" :model="loginForm"
               label-width="20%">
        <el-form-item style="width: 100%" :label="$t('business.user.username')" prop="username">
          <el-input v-model="loginForm.username"></el-input>
        </el-form-item>
        <el-form-item style="width: 100%" :label="$t('business.user.password')" prop="password">
          <el-input v-model="loginForm.password" type="password"></el-input>
        </el-form-item>
        <div style="height: 30px">
          <div style="float: right">
            <el-button @click="loginPageOpen=false" :disabled="loginLoading"
                       v-has-permissions="{resource:'ldap',verb:'create'}">
              {{ $t("commons.button.cancel") }}
            </el-button>
            <el-button type="primary" @click="loginTest" :disabled="loginLoading"
                       v-has-permissions="{resource:'ldap',verb:'create'}">{{ $t("commons.button.confirm") }}
            </el-button>
          </div>
        </div>
      </el-form>
    </el-dialog>
    <el-dialog :visible.sync="importUserPageOpen" :title="$t('business.user.import_user')" style="height: 900px">
      <span>{{ $t("business.user.ldap_helper") }}</span>
      <div style="text-align: right;margin-bottom: 10px">
        <el-input v-model="searchName" suffix-icon="el-icon-search" style="width: 30%" size="mini" clearable @change="handleSearch" />
      </div>
      <el-table
              :data="tableUsers.slice((pageConfig.currentPage-1)*pageConfig.pageSize,pageConfig.currentPage*pageConfig.pageSize)"
              :loading="connectLoading"
              @selection-change="handleSelectionChange">
        <el-table-column type="selection" fix :selectable="importAvailable"></el-table-column>
        <el-table-column :label="$t('commons.table.name')" prop="name" min-width="100" fix>
        </el-table-column>
        <el-table-column :label="$t('business.user.nickname')" prop="nickName" min-width="100" fix>
        </el-table-column>
        <el-table-column :label="$t('business.user.email')" prop="email" min-width="100" fix>
        </el-table-column>
        <el-table-column :label="$t('business.user.import_enable')" prop="available" min-width="100" fix>
          <template v-slot:default="{row}">
            {{ row.available ? $t("commons.bool.true") : $t("commons.bool.false") }}
          </template>
        </el-table-column>
      </el-table>
      <div style="height: 30px;margin-top: 10px;margin-bottom: 20px;text-align: right">
        <el-pagination
                :current-page.sync="pageConfig.currentPage"
                :page-size="pageConfig.pageSize"
                layout="prev, pager, next"
                :total="tableUsers.length">
        </el-pagination>
        <el-button @click="importUserPageOpen=false" :disabled="importLoading"
                   v-has-permissions="{resource:'ldap',verb:'create'}">
          {{ $t("commons.button.cancel") }}
        </el-button>
        <el-button type="primary" @click="userImport" :disabled="importLoading"
                   v-has-permissions="{resource:'ldap',verb:'create'}">{{ $t("commons.button.confirm") }}
        </el-button>
      </div>
    </el-dialog>
  </layout-content>
</template>

<script>
import "codemirror/lib/codemirror.css"
import "codemirror/theme/ayu-dark.css"
import "codemirror/mode/javascript/javascript"
import LayoutContent from "@/components/layout/LayoutContent"
import Rule from "@/utils/rules"
import {createLdap, getLdap, importUser, syncLdap, testConnect, testLogin, updateLdap} from "@/api/ldap"

export default {
  name: "LDAP",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      form: {
        mapping: "{\n" +
          "   \"Name\":\"cn\",\n" +
          "   \"NickName\":\"cn\",\n" +
          "   \"Email\":\"mail\"\n" +
          "}",
        sizeLimit: 1000,
        timeLimit: 30,
      },
      loading: false,
      rules: {
        address: [Rule.RequiredRule],
        port: [Rule.RequiredRule],
        username: [Rule.RequiredRule],
        password: [Rule.RequiredRule],
        dn: [Rule.RequiredRule],
        filter: [Rule.RequiredRule],
        mapping: [Rule.RequiredRule],
        timeLimit: [Rule.RequiredRule],
        sizeLimit: [Rule.RequiredRule],
      },
      isSubmitGoing: false,
      options: {
        value: "",
        mode: "application/json",
        theme: "ayu-dark",
        lineNumbers: true,
        tabSize: 2,
        lineWrapping: true,
        gutters: ["CodeMirror-lint-markers"],
      },
      selects: [],
      users: [],
      loginPageOpen: false,
      loginLoading: false,
      loginForm: {},
      importUserPageOpen: false,
      importLoading: false,
      connectLoading: false,
      oldConfig: {},
      pageConfig: {
        currentPage: 1,
        pageSize: 10,
      },
      tableUsers: [],
      searchName: ""
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
          message: this.$t("business.user.test_result", { count: res.data })
        })
      }).finally(() => {
        this.loading = false
        this.connectLoading = false
      })
    },
    openLoginPage () {
      this.loginPageOpen = true
      this.loginForm = {}
    },
    openImportPage () {
      this.searchName = ""
      this.loading = true
      syncLdap({}).then(res => {
        this.users = res.data
        this.tableUsers = this.users
        this.importUserPageOpen = true
      }).finally(() => {
        this.loading = false
      })
    },
    importAvailable (row) {
      return row.available
    },
    remake () {
      this.form = JSON.parse(JSON.stringify(this.oldConfig))
    },
    loginTest () {
      let isFormReady = false
      this.$refs["loginForm"].validate((valid) => {
        if (valid) {
          isFormReady = true
        }
      })
      if (!isFormReady) {
        return
      }
      this.loginLoading = true
      testLogin(this.loginForm).then(() => {
        this.$message({
          type: "success",
          message: this.$t("business.user.login_success")
        })
      }).finally(() => {
        this.loginLoading = false
      })
    },
    userImport () {
      if (this.selects.length === 0) {
        return
      }
      let req = {
        users: this.selects
      }
      this.importLoading = true
      importUser(req).then(res => {
        if (res.success) {
          this.$message({
            type: "success",
            message: this.$t("business.user.import_user_success")
          })
          this.importUserPageOpen = false
        } else {
          this.$message({
            type: "failed",
            message: this.$t("business.user.import_user_failed", { user: res.failures })
          })
        }
      }).finally(() => {
        this.importLoading = false
      })
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
        createLdap(this.form).then(() => {
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
        updateLdap(this.form).then(() => {
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
      getLdap().then((res) => {
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
