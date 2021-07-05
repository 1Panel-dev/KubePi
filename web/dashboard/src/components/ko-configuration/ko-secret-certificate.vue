<template>
  <div style="margin-top: 20px">
    <ko-card title="Certificate" ref="form" :model="form">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Private Key" prop="crt" required>
              <ko-form-item itemType="textarea" v-model="form.crt" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Certificate" prop="key" required>
              <ko-form-item itemType="textarea" v-model="form.key" @change.native="transform"></ko-form-item>
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
  name: "KoSecretCertificate",
  components: { KoFormItem, KoCard },
  props: {
    certificateObj: Object
  },
  data () {
    return {
      form: {
        crt: "",
        key: ""
      }
    }
  },
  methods: {
    transform () {
      const { Base64 } = require("js-base64")
      const result = {
        "tls.crt": Base64.encode(this.form.crt),
        "tls.key": Base64.encode(this.form.key)
      }
      this.$emit("update:certificateObj", result)
    }
  },
  mounted () {
    if (this.certificateObj) {
      const { Base64 } = require("js-base64")
      this.form = {
        crt: Base64.decode(this.certificateObj["tls.crt"]),
        key: Base64.decode(this.certificateObj["tls.key"])
      }
    }
  }
}
</script>

<style scoped>

</style>