<template>
  <div style="margin-top: 20px">
    <ko-card title="Security Context">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Privileged" prop="privileged">
              <ko-form-item @change="privilegedChanged" itemType="radio" v-model="form.privileged" :radios="privileged_list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Privilege Escalation" prop="allowPrivilegeEscalation">
              <ko-form-item itemType="radio" v-model="form.allowPrivilegeEscalation" :radios="privileged_escalation_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Run as Non-Root" prop="runAsNonRoot">
              <ko-form-item itemType="radio" v-model="form.runAsNonRoot" :radios="non_root_list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Read-Only Root Filesystem" prop="readOnlyRootFilesystem">
              <ko-form-item itemType="radio" v-model="form.readOnlyRootFilesystem" :radios="ready_only_root_files_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Run as User ID" prop="runAsUser">
              <ko-form-item itemType="input" v-model="form.runAsUser" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Add Capabilities" prop="capabilities.add">
              <ko-form-item multiple itemType="select2" filterable allow-create v-model="form.capabilities.add" :selections="capability_list" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Drop Capabilities" prop="capabilities.drop">
              <ko-form-item multiple itemType="select2" filterable allow-create v-model="form.capabilities.drop" :selections="capability_list" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
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
    isReadOnly: Boolean,
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
      capability_list: ["ALL", "AUDIT_CONTROL", "AUDIT_WRITE", "BLOCK_SUSPEND", "CHOWN", "DAC_OVERRIDE", "DAC_READ_SEARCH", "FOWNER", "FSETID", "IPC_LOCK", "IPC_OWNER", "KILL", "LEASE", "LINUX_IMMUTABLE", "MAC_ADMIN", "MAC_OVERRIDE", "MKNOD", "NET_ADMIN", "NET_BIND_SERVICE", "NET_BROADCAST", "NET_RAW", "SETFCAP", "SETGID", "SETPCAP", "SETUID", "SYSLOGSYS_ADMIN", "SYS_BOOT", "SYS_CHROOT", "SYS_MODULE", "SYS_NICE", "SYS_PACCT", "SYS_PTRACE", "SYS_RAWIO", "SYS_RESOURCE", "SYS_TIME", "SYS_TTY_CONFIG", "WAKE_ALARM"],
      form: {
        privileged: false,
        allowPrivilegeEscalation: true,
        runAsNonRoot: false,
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