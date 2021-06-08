<template>
  <div style="margin-top: 20px">
    <el-row :gutter="20">
      <el-col :span="12">
        <ko-form-item @change="privilegedChanged" labelName="Privileged" clearable itemType="radio" v-model="form.privileged" :radios="privileged_list" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Privilege Escalation" clearable itemType="radio" v-model="form.allowPrivilegeEscalation" :radios="privileged_escalation_list" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Run as Non-Root" clearable itemType="radio" v-model="form.runAsRoot" :radios="non_root_list" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Read-Only Root Filesystem" clearable itemType="radio" v-model="form.readOnlyRootFilesystem" :radios="ready_only_root_files_list" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 30px">
      <el-col :span="12">
        <ko-form-item labelName="Run as User ID" clearable itemType="input" v-model="form.runAsUser" />
      </el-col>
    </el-row>
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <ko-form-item labelName="Add Capabilities" multiple clearable itemType="select" v-model="form.capabilities.add" :selections="capability_list" />
      </el-col>
      <el-col :span="12">
        <ko-form-item labelName="Drop Capabilities" multiple clearable itemType="select" v-model="form.capabilities.drop" :selections="capability_list" />
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
      privileged_list: [
        { label: "No", value: false },
        { label: "Yes: container has full access to the host", value: true },
      ],
      privileged_escalation_list: [
        { label: "No", value: false, disabledOption: false},
        { label: "Yes: container can gain more privileges than its parent process", value: true },
      ],
      non_root_list: [
        { label: "No", value: false },
        { label: "Yes: container must run as a non-root user", value: true },
      ],
      ready_only_root_files_list: [
        { label: "No", value: false },
        { label: "Yes: container has a read-only root filesystem", value: true },
      ],
      capability_list: [
        { label: "ALL", value: "ALL" },
        { label: "AUDIT_CONTROL", value: "AUDIT_CONTROL" },
        { label: "AUDIT_WRITE", value: "AUDIT_WRITE" },
        { label: "BLOCK_SUSPEND", value: "BLOCK_SUSPEND" },
      ],
      form: {
        privileged: "",
        allowPrivilegeEscalation: "",
        runAsRoot: "",
        readOnlyRootFilesystem: false,
        runAsUser: "",
        capabilities: {
          add: [],
          drop: [],
        },
      },
    }
  },
  methods: {
    privilegedChanged (value) {
      this.privileged_escalation_list[0].disabledOption = value
    }
  },
}
</script>