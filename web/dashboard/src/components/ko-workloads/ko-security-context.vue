<template>
  <div style="margin-top: 20px">
    <ko-card title="Security Context">
      <el-row :gutter="20">
        <el-col :span="12">
          <ko-form-item @change="privilegedChanged" labelName="Privileged" itemType="radio" v-model="form.privileged" :radios="privileged_list" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Privilege Escalation" itemType="radio" v-model="form.allowPrivilegeEscalation" :radios="privileged_escalation_list" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <ko-form-item labelName="Run as Non-Root" itemType="radio" v-model="form.runAsNonRoot" :radios="non_root_list" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Read-Only Root Filesystem" itemType="radio" v-model="form.readOnlyRootFilesystem" :radios="ready_only_root_files_list" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 30px">
        <el-col :span="12">
          <ko-form-item labelName="Run as User ID" itemType="input" v-model="form.runAsUser" />
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 20px">
        <el-col :span="12">
          <ko-form-item labelName="Add Capabilities" multiple itemType="select" v-model="form.capabilities.add" :selections="capability_list" />
        </el-col>
        <el-col :span="12">
          <ko-form-item labelName="Drop Capabilities" multiple itemType="select" v-model="form.capabilities.drop" :selections="capability_list" />
        </el-col>
      </el-row>
    </ko-card>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoSecurity",
  components: { KoFormItem, KoCard },
  props: {
    securityContextParentObj: Object,
  },
  data() {
    return {
      privileged_list: [
        { label: "No", value: false },
        { label: "Yes: container has full access to the host", value: true },
      ],
      privileged_escalation_list: [
        { label: "No", value: false, disabledOption: false },
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
        runAsNonRoot: "",
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
    privilegedChanged(value) {
      this.privileged_escalation_list[0].disabledOption = value
    },
    transformation(parentFrom) {
      if (!parentFrom.securityContext) {
        parentFrom.securityContext = {}
      }
      if (this.form.privileged) {
        parentFrom.securityContext.privileged = this.form.privileged
      }
      if (this.form.allowPrivilegeEscalation) {
        parentFrom.securityContext.allowPrivilegeEscalation = this.form.allowPrivilegeEscalation
      }
      if (this.form.runAsNonRoot) {
        parentFrom.securityContext.runAsNonRoot = this.form.runAsNonRoot
      }
      if (this.form.readOnlyRootFilesystem) {
        parentFrom.securityContext.readOnlyRootFilesystem = this.form.readOnlyRootFilesystem
      }
      if (this.form.runAsUser) {
        parentFrom.securityContext.runAsUser = this.form.runAsUser
      }
      if (this.form.capabilities.add.length !== 0 || this.form.capabilities.drop.length !== 0) {
        parentFrom.securityContext.capabilities = {}
        if (this.form.capabilities.add.length !== 0) {
          parentFrom.securityContext.capabilities.add = this.form.capabilities.add
        }
        if (this.form.capabilities.drop.length !== 0) {
          parentFrom.securityContext.capabilities.drop = this.form.capabilities.drop
        }
      }
    },
  },
  mounted() {
    if (this.securityContextParentObj) {
      if (this.securityContextParentObj.securityContext) {
        if (this.securityContextParentObj.securityContext.privileged !== undefined) {
          this.form.privileged = this.securityContextParentObj.securityContext.privileged
        }
        if (this.securityContextParentObj.securityContext.allowPrivilegeEscalation !== undefined) {
          this.form.allowPrivilegeEscalation = this.securityContextParentObj.securityContext.allowPrivilegeEscalation
        }
        if (this.securityContextParentObj.securityContext.runAsNonRoot !== undefined) {
          this.form.runAsNonRoot = this.securityContextParentObj.securityContext.runAsNonRoot
        }
        if (this.securityContextParentObj.securityContext.readOnlyRootFilesystem !== undefined) {
          this.form.readOnlyRootFilesystem = this.securityContextParentObj.securityContext.readOnlyRootFilesystem
        }
        if (this.securityContextParentObj.securityContext.runAsUser) {
          this.form.runAsUser = this.securityContextParentObj.securityContext.runAsUser
        }
        if (this.securityContextParentObj.securityContext.capabilities) {
          if (this.securityContextParentObj.securityContext.capabilities.add) {
            this.form.capabilities.add = this.securityContextParentObj.securityContext.capabilities.add
          }
          if (this.securityContextParentObj.securityContext.capabilities.drop) {
            this.form.capabilities.drop = this.securityContextParentObj.securityContext.capabilities.drop
          }
        }
      }
    }
  },
}
</script>