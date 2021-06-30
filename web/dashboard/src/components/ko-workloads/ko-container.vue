<template>
  <div style="margin-top: 20px">
    <ko-card title="Container">
      <el-form label-position="top" ref="form" :rules="rules" :model="form">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Container Name" prop="name">
              <ko-form-item itemType="input" v-model="form.name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Container Image" prop="image">
              <ko-form-item placeholder="e.g. nginx:latest" itemType="input" v-model="form.image" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Pull Policy" prop="imagePullPolicy">
              <ko-form-item itemType="select" v-model="form.imagePullPolicy" :selections="image_pull_policy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item label="Pull Secrets" prop="imagePullSecrets">
            <ko-form-item itemType="select" multiple v-model="form.imagePullSecrets" :selections="secret_list" />
          </el-form-item>
        </el-row>
      </el-form>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"
import Rule from "@/utils/rules"
import { listSecrets } from "@/api/secrets"

export default {
  name: "KoContainer",
  components: { KoFormItem, KoCard },
  props: {
    containerParentObj: Object,
  },
  data() {
    return {
      form: {
        name: "",
        image: "",
        imagePullPolicy: "",
      },
      rules: {
        image: [Rule.RequiredRule],
      },
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "ifNotPresent", value: "ifNotPresent" },
        { label: "Never", value: "Never" },
      ],
      secret_list: [],
    }
  },
  methods: {
    loadSecrets() {
      listSecrets(this.$route.params.cluster).then((res) => {
        this.secret_list = []
        for (const secret of res.items) {
          this.secret_list.push({ label: secret.metadata.name, value: secret.metadata.name })
        }
      })
    },
    checkIsValid() {
      let isValid = true
      this.$refs["form"].validate((valid) => {
        isValid = valid
      })
      return isValid
    },
    transformation(parentFrom) {
      if (this.form.name) {
        parentFrom.name = this.form.name
      }
      if (this.form.image) {
        parentFrom.image = this.form.image
      }
      if (this.form.imagePullPolicy) {
        parentFrom.imagePullPolicy = this.form.imagePullPolicy
      }
    },
  },
  mounted() {
    this.loadSecrets()
    if (this.containerParentObj) {
      if (this.containerParentObj.name) {
        this.form.name = this.containerParentObj.name
      }
      if (this.containerParentObj.image) {
        this.form.image = this.containerParentObj.image
      }
      if (this.containerParentObj.imagePullPolicy) {
        this.form.imagePullPolicy = this.containerParentObj.imagePullPolicy
      }
    }
  },
}
</script>