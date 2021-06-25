<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-row>
          <el-col :span="12">
            <ko-form-item labelName="Container Name" clearable itemType="input" v-model="form.name" />
          </el-col>
          <el-col :span="12">
            <ko-form-item labelName="Pull Policy" clearable itemType="select" v-model="form.imagePullPolicy" :selections="image_pull_policy_list" />
          </el-col>
        </el-row>
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Container Image" placeholder="e.g. nginx:latest" clearable itemType="input" v-model="form.image" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoGeneral",
  components: { KoFormItem },
  props: {
    generalParentObj: Object,
  },
  data() {
    return {
      form: {
        name: "",
        image: "",
        imagePullPolicy: "",
      },
      image_pull_policy_list: [
        { label: "Always", value: "Always" },
        { label: "ifNotPresent", value: "ifNotPresent" },
        { label: "Never", value: "Never" },
      ],
    }
  },
  methods: {
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
    if (this.generalParentObj) {
      if (this.generalParentObj.name) {
        this.form.name = this.generalParentObj.name
      }
      if (this.generalParentObj.image) {
        this.form.image = this.generalParentObj.image
      }
      if (this.generalParentObj.imagePullPolicy) {
        this.form.imagePullPolicy = this.generalParentObj.imagePullPolicy
      }
    }
  },
}
</script>