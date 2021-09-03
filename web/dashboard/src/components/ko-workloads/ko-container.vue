<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.container')">
      <el-form label-position="top" ref="form" :rules="rules" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.container_name')" prop="name">
              <ko-form-item itemType="input" @input="changeName" v-model="form.name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.container_image')" prop="image">
              <ko-form-item placeholder="e.g. nginx:latest" itemType="input" v-model="form.image" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.pull_policy')" prop="imagePullPolicy">
              <ko-form-item itemType="select" :noClear="true" v-model="form.imagePullPolicy" :selections="image_pull_policy_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-form-item :label="$t('business.workload.pull_secrets')" prop="imagePullSecrets">
            <ko-form-item itemType="select2" multiple v-model="form.imagePullSecrets" :selections="secret_list" />
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

export default {
  name: "KoContainer",
  components: { KoFormItem, KoCard },
  props: {
    containerParentObj: Object,
    secretList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    secretList: {
      handler(newName) {
        this.secret_list = []
        for (const s of newName) {
          this.secret_list.push(s.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      form: {
        name: "",
        image: "",
        imagePullPolicy: "",
        imagePullSecrets: "",
      },
      rules: {
        name: [Rule.CommonNameRule],
        image: [Rule.RequiredRule],
        imagePullPolicy: [Rule.SelectRule],
      },
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "IfNotPresent", value: "IfNotPresent" },
        { label: "Never", value: "Never" },
      ],
      secret_list: [],
    }
  },
  methods: {
    changeName(val) {
      this.$emit("updateContanerList", val)
    },
    checkIsValid() {
      let isValid = true
      this.$refs["form"].validate((valid) => {
        isValid = valid
      })
      return isValid
    },
    transformation(parentFrom, specForm) {
      parentFrom.name = this.form.name || undefined
      parentFrom.image = this.form.image || undefined
      parentFrom.imagePullPolicy = this.form.imagePullPolicy || undefined
      specForm.imagePullSecrets = this.form.imagePullSecrets || undefined
    },
  },
  mounted() {
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