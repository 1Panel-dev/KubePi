<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item labelName="Privileged" clearable itemType="radio" v-model="form.item1" :radios="privileged_list1" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Privilege Escalation" clearable itemType="radio" v-model="form.item2" :radios="privileged_list2" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Run as Non-Root" clearable itemType="radio" v-model="form.item3" :radios="privileged_list3" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Read-Only Root Filesystem" clearable itemType="radio" v-model="form.item4" :radios="privileged_list4" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 30px">
      <el-col :span="12">
        <ko-form-item labelName="Run as User ID" clearable itemType="input" v-model="form.user_id" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Add Capabilities" clearable itemType="input" v-model="form.add_ca" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Drop Capabilities" clearable itemType="input" v-model="form.drop_ca" />
      </el-col>
    </el-row>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoSecurity",
  components: { KoFormItem },
  data() {
    return {
      privileged_list1: [
        { label: "No", value: "No" },
        { label: "Yes: container has full access to the host", value: "Yes: container has full access to the host" },
      ],
      privileged_list2: [{ label: "No", value: "No" }, { label: "Yes: container can gain more privileges than its parent process" }],
      privileged_list3: [{ label: "No", value: "No" }, { label: "Yes: container must run as a non-root user" }],
      privileged_list4: [{ label: "No", value: "No" }, { label: "Yes: container has a read-only root filesystem" }],
      form: {
        item1:  "",
        item2:  "",
        item3:  "",
        item4:  "",
        user_id: "",
        add_ca: "",
        drop_ca: "",
      },
    }
  },
  methods: {},
}
</script>