<template>
  <div style="margin-top: 20px">
    <ko-card title="Keys">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Public Key" prop="publicKey" required>
              <ko-form-item itemType="textarea"  v-model="form.publicKey" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Private Key" prop="privateKey" required>
              <ko-form-item itemType="textarea"  v-model="form.privateKey" @change.native="transform"></ko-form-item>
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
      }
    }
  },
  methods: {
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
    if (this.dataObj) {
      const { Base64 } = require("js-base64")
      this.form = {
        privateKey: Base64.decode(this.dataObj["ssh-privatekey"]),
        publicKey: Base64.decode(this.dataObj["ssh-publickey"])
      }
    }
  }
}
</script>

<style scoped>

</style>