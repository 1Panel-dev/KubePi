<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.configuration.data')">
      <el-form label-position="top" ref="form" :model="form" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.configuration.kind')">
              <el-radio-group v-model="registryType" @change="changeType(registryType)">
                <el-radio :label="'Custom'" ></el-radio>
                <el-radio :label="'DockerHub'"></el-radio>
                <el-radio :label="'Quay.io'"></el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20" v-if="registryType==='Custom'">
          <el-col :span="12">
            <el-form-item :label="$t('business.configuration.registry_domain_name')" prop="registry">
              <ko-form-item itemType="input" v-model="form.registry" placeholder="e.g. some.docker.io"
                            @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.configuration.username')" prop="username">
              <ko-form-item itemType="input" v-model="form.username" @change.native="transform"></ko-form-item>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.configuration.password')" prop="password">
              <ko-form-item itemType="input" type="password" v-model="form.password" @change.native="transform"
                            show-password></ko-form-item>
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
import Rules from "../../../../kubepi/src/utils/rules"

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
      registryType: "Custom",
      rules:{
        registry:[
          Rules.RequiredRule
        ]
      }
    }
  },
  methods: {
    transform () {
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
          [".dockerconfigjson"]: Base64.encode(JSON.stringify(result))
        }
        this.$emit("update:dataObj", data)
      }
    },
    changeType(type) {
      switch (type){
        case "DockerHub":
          this.form.registry = 'index.docker.io/v1/'
          break;
        case "Quay.io":
          this.form.registry = 'quay.io'
          break;
        default:
          this.form.registry = ''
          break;
      }
      this.transform()
    }
  },
  mounted () {
    if (this.dataObj && this.dataObj[".dockerconfigjson"]) {
      const { Base64 } = require("js-base64")
      const value = Base64.decode(this.dataObj[".dockerconfigjson"])
      const obj = JSON.parse(value)
      const auths = obj.auths
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
