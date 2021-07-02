<template>
  <div style="margin-top: 20px">
    <ko-card title="Data">
      <el-form label-position="top" ref="form" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Registry Domain Name" prop="registry">
              <ko-form-item itemType="input" v-model="form.registry" placeholder="e.g. some.docker.io" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Username" prop="username">
              <ko-form-item itemType="input" v-model="form.username" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Password" prop="password">
              <ko-form-item itemType="input" type="password" v-model="form.password" @change.native="transform"></ko-form-item>
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
  name: "KoSecretDockerData",
  components: { KoFormItem, KoCard },
  props: {
    dataObj: Object
  },
  data () {
    return {
      form: {
        registry: "",
        username: "",
        password: ""
      },
    }
  },
  methods: {
    // transformation (parentFrom) {
    //   if (this.form.registry) {
    //     const result = {
    //       auths: {
    //         [this.form.registry]: {
    //           username: this.form.username,
    //           password: this.form.password
    //         }
    //       }
    //     }
    //     const { Base64 } = require("js-base64")
    //     parentFrom.data[".dockerconfigjson"] = Base64.encode(JSON.stringify(result))
    //   }
    // },
    transform() {
      if (this.form.registry) {
        const result = {
          auths: {
            [this.form.registry]: {
              username: this.form.username,
              password: this.form.password
            }
          }
        }
        const { Base64 } = require("js-base64")
        const data = {
          [".dockerconfigjson"]:  Base64.encode(JSON.stringify(result))
        }
        this.$emit("update:dataObj" ,data)
      }
    }
  },
  mounted () {
    if (this.dataObj && this.dataObj[".dockerconfigjson"]) {
      const { Base64 } = require("js-base64")
      const value = Base64.decode(this.dataObj[".dockerconfigjson"])
      const auths = JSON.parse(value)
      for (const key in auths) {
        if (Object.prototype.hasOwnProperty.call(auths, key)) {
          this.form = {
            registry: key,
            username: auths[key].username,
            password: auths[key].password
          }
        }
      }
    }
  }
}
</script>

<style scoped>

</style>