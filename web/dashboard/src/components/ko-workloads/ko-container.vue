<template>
  <div style="margin-top: 20px">
    <ko-card title="Container">
      <el-row :gutter="20">
        <el-col :span="6">
          <ko-form-item labelName="Container Name" itemType="input" v-model="form.name" />
        </el-col>
        <el-col :span="6">
          <ko-form-item labelName="Pull Policy" itemType="select" v-model="form.imagePullPolicy" :selections="image_pull_policy_list" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Container Image" placeholder="e.g. nginx:latest" itemType="input" v-model="form.image" />
        </el-col>
      </el-row>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

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