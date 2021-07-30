<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.configuration.authentication')">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item  :label="$t('business.configuration.username')" prop="username" required>
              <ko-form-item itemType="input" v-model="form.username" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item  :label="$t('business.configuration.password')" prop="password">
              <ko-form-item itemType="password" v-model="form.password" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoSecretAuthentication",
  components: { KoFormItem, KoCard },
  props: {
    authenticationObj: Object
  },
  data () {
    return {
      form: {
        username: "",
        password: ""
      }
    }
  },
  methods: {
    transform () {
      const { Base64 } = require("js-base64")
      const auth = {
        username: Base64.encode(this.form.username),
        password:  Base64.encode(this.form.password)
      }
      this.$emit("update:authenticationObj" ,auth)
    }
  },
  mounted () {
    if (this.authenticationObj) {
      const { Base64 } = require("js-base64")
      if (this.authenticationObj.username) {
          this.form.username = Base64.decode(this.authenticationObj.username)
      }
      if (this.authenticationObj.password) {
        this.form.password = Base64.decode(this.authenticationObj.password)
      }
    }
  }
}
</script>

<style scoped>

</style>
