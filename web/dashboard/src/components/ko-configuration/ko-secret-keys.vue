<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.configuration.certificate')">
      <el-form label-position="top" ref="form" :model="form" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item  :label="$t('business.configuration.publicKey')" prop="publicKey" required>
              <ko-form-item itemType="textarea" v-model="form.publicKey" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item  :label="$t('business.configuration.privateKey')" prop="privateKey" required>
              <ko-form-item itemType="textarea" v-model="form.privateKey" @change.native="transform"></ko-form-item>
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
import Rule from "@/utils/rules"

export default {
  name: "KoSecretKeys",
  components: { KoFormItem, KoCard },
  props: {
    dataObj: Object
  },
  data () {
    return {
      form: {
        privateKey: "",
        publicKey: ""
      },
      rules: {
        publicKey: [Rule.RequiredRule]
      }
    }
  },
  methods: {
    checkIsValid () {
      let result = {
        success: false,
        messages: ["测试用户"]
      }
      this.$refs["form"].validate((valid) => {
        result.success = valid
      })
      return result
    },
    transform () {
      const { Base64 } = require("js-base64")
      const result = {
        "ssh-privatekey": Base64.encode(this.form.privateKey),
        "ssh-publickey": Base64.encode(this.form.publicKey)
      }
      this.$emit("update:dataObj", result)
    }
  },
  mounted () {
    if (this.dataObj && this.dataObj !== {}) {
      const { Base64 } = require("js-base64")
      if (this.dataObj["ssh-privatekey"]) {
        this.form.privateKey = Base64.decode(this.dataObj["ssh-privatekey"])
      }
      if (this.dataObj["ssh-publickey"]) {
        this.form.publicKey = Base64.decode(this.dataObj["ssh-publickey"])
      }
    }
  }
}
</script>

<style scoped>

</style>
