<template>
  <layout-content :header="'LDAP'">
    <el-col :span="2"><br/></el-col>
    <el-col :span="15">
      <div class="grid-content bg-purple-light">
        <el-form ref="form" v-loading="loading" label-position="left" :rules="rules" :model="form" label-width="150px">
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_address')" prop="address" >
            <el-input v-model="form.address"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_port')" prop="port" >
            <el-input v-model="form.port" type="number"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_username')" prop="username">
            <el-input v-model="form.username"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_password')" prop="password">
            <el-input type="password" v-model="form.password"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_filter_dn')" prop="dn">
            <el-input v-model="form.dn"></el-input>
          </el-form-item>
          <el-form-item style="width: 100%" :label="$t('business.user.ldap_filter_rule')" prop="filter">
            <el-input v-model="form.filter"></el-input>
          </el-form-item>
          <el-form-item>
            <div style="font-size: 12px;color: #4E5051;">
              <svg class="icon" aria-hidden="true">
                <use xlink:href="#icontishi11"></use>
              </svg>
              {{ $t("business.user.ldap_helper") }}
            </div>
          </el-form-item>

          <div style="float: right">
            <el-form-item>
              <el-button @click="sync" :disabled="isSubmitGoing" v-has-permissions="{resource:'ldap',verb:'create'}">{{
                  $t("commons.button.sync")
                }}
              </el-button>
              <el-button type="primary" @click="onSubmit" :disabled="isSubmitGoing" v-has-permissions="{resource:'ldap',verb:'create'}">{{ $t("commons.button.confirm") }}
              </el-button>
            </el-form-item>
          </div>

        </el-form>
      </div>
    </el-col>
    <el-col :span="4"><br/></el-col>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import Rule from "@/utils/rules"
import {createLdap, getLdap, syncLdap, updateLdap} from "@/api/ldap"

export default {
  name: "LDAP",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      form: {},
      loading: false,
      rules: {
        address: [Rule.RequiredRule],
        port: [Rule.RequiredRule],
        username: [Rule.RequiredRule],
        password: [Rule.RequiredRule],
        dn: [Rule.RequiredRule],
        filter: [Rule.RequiredRule],
      },
      isSubmitGoing: false
    }
  },
  methods: {
    sync () {
      if (this.form.uuid === undefined || this.form.uuid === "") {
        this.$message({
          type: "warning",
          message: this.$t("business.user.ldap_sync_error")
        })
        return
      }
      this.isSubmitGoing = true
      syncLdap(this.form.uuid, {}).then(() => {
        this.$message({
          type: "success",
          message: this.$t("business.user.ldap_sync")
        })
      }).finally(() => {
        this.isSubmitGoing = false
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
          console.log(this.form)
        }
      }).finally(() => {
        this.loading = false
      })
    }
  },
  created () {
    this.list()
  }
}
</script>

<style scoped>

</style>
