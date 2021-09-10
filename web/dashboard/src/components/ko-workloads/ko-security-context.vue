<template>
  <div>
    <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.privileged')" prop="privileged">
            <ko-form-item @change="privilegedChanged" itemType="radio" v-model="form.privileged" :radios="privileged_list" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.privilege_escalation')" prop="allowPrivilegeEscalation">
            <ko-form-item itemType="radio" v-model="form.allowPrivilegeEscalation" :radios="privileged_escalation_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_non_root')" prop="runAsNonRoot">
            <ko-form-item itemType="radio" v-model="form.runAsNonRoot" :radios="non_root_list" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.read_only_root_filesystem')" prop="readOnlyRootFilesystem">
            <ko-form-item itemType="radio" v-model="form.readOnlyRootFilesystem" :radios="ready_only_root_files_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_user')" prop="runAsUser">
            <ko-form-item itemType="input" placeholder="runAsUser" v-model="form.runAsUser" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_group')" prop="runAsGroup">
            <ko-form-item itemType="input" placeholder="runAsGroup" v-model="form.runAsGroup" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.proc_mount')" prop="procMount">
            <ko-form-item itemType="input" placeholder="procMount" v-model="form.procMount" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.add_capabilities')" prop="capabilities.add">
            <ko-form-item multiple itemType="select2" filterable allow-create v-model="form.capabilities.add" :selections="capability_list" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.drop_capabilities')" prop="capabilities.drop">
            <ko-form-item multiple itemType="select2" filterable allow-create v-model="form.capabilities.drop" :selections="capability_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <span style="font-size: 16px">{{$t('business.workload.seLinuxOptions')}}</span>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="level" prop="seLinuxOptions.level">
            <ko-form-item itemType="input" placeholder="level" v-model="form.seLinuxOptions.level" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="role" prop="seLinuxOptions.role">
            <ko-form-item itemType="input" placeholder="role" v-model="form.seLinuxOptions.role" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="type" prop="seLinuxOptions.type">
            <ko-form-item itemType="input" placeholder="type" v-model="form.seLinuxOptions.type" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="user" prop="seLinuxOptions.user">
            <ko-form-item itemType="input" placeholder="user" v-model="form.seLinuxOptions.user" />
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoSecurity",
  components: { KoFormItem },
  props: {
    securityContextParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      privileged_list: [
        { label: this.$t("business.workload.no"), value: false },
        { label: this.$t("business.workload.full_access"), value: true },
      ],
      privileged_escalation_list: [
        { label: this.$t("business.workload.no"), value: false, disabledOption: false },
        { label: this.$t("business.workload.gain_more_privileges"), value: true },
      ],
      non_root_list: [
        { label: this.$t("business.workload.no"), value: false },
        { label: this.$t("business.workload.non_root"), value: true },
      ],
      ready_only_root_files_list: [
        { label: this.$t("business.workload.no"), value: false },
        { label: this.$t("business.workload.filesystem_read_only"), value: true },
      ],
      capability_list: ["ALL", "AUDIT_CONTROL", "AUDIT_WRITE", "BLOCK_SUSPEND", "CHOWN", "DAC_OVERRIDE", "DAC_READ_SEARCH", "FOWNER", "FSETID", "IPC_LOCK", "IPC_OWNER", "KILL", "LEASE", "LINUX_IMMUTABLE", "MAC_ADMIN", "MAC_OVERRIDE", "MKNOD", "NET_ADMIN", "NET_BIND_SERVICE", "NET_BROADCAST", "NET_RAW", "SETFCAP", "SETGID", "SETPCAP", "SETUID", "SYSLOGSYS_ADMIN", "SYS_BOOT", "SYS_CHROOT", "SYS_MODULE", "SYS_NICE", "SYS_PACCT", "SYS_PTRACE", "SYS_RAWIO", "SYS_RESOURCE", "SYS_TIME", "SYS_TTY_CONFIG", "WAKE_ALARM"],
      form: {
        privileged: null,
        allowPrivilegeEscalation: null,
        runAsNonRoot: null,
        readOnlyRootFilesystem: null,
        runAsUser: "",
        runAsGroup: "",
        procMount: "",
        capabilities: {
          add: [],
          drop: [],
        },
        seLinuxOptions: {
          level: "",
          role: "",
          type: "",
          user: "",
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
      parentFrom.securityContext.privileged = this.form.privileged || undefined
      parentFrom.securityContext.allowPrivilegeEscalation = this.form.allowPrivilegeEscalation || undefined
      parentFrom.securityContext.runAsNonRoot = this.form.runAsNonRoot || undefined
      parentFrom.securityContext.readOnlyRootFilesystem = this.form.readOnlyRootFilesystem || undefined
      parentFrom.securityContext.runAsUser = this.form.runAsUser || undefined
      parentFrom.securityContext.runAsGroup = this.form.runAsGroup || undefined
      parentFrom.securityContext.procMount = this.form.procMount || undefined
      parentFrom.securityContext.capabilities = {
        add: this.form.capabilities.add.length !== 0 ? this.form.capabilities.add : undefined,
        drop: this.form.capabilities.drop.length !== 0 ? this.form.capabilities.drop : undefined,
      }
      parentFrom.securityContext.seLinuxOptions = {
        level: this.form.seLinuxOptions.level || undefined,
        role: this.form.seLinuxOptions.role || undefined,
        type: this.form.seLinuxOptions.type || undefined,
        user: this.form.seLinuxOptions.user || undefined,
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
        if (this.securityContextParentObj.securityContext.runAsGroup) {
          this.form.runAsGroup = this.securityContextParentObj.securityContext.runAsGroup
        }
        if (this.securityContextParentObj.securityContext.procMount) {
          this.form.procMount = this.securityContextParentObj.securityContext.procMount
        }
        if (this.securityContextParentObj.securityContext.capabilities) {
          if (this.securityContextParentObj.securityContext.capabilities.add) {
            this.form.capabilities.add = this.securityContextParentObj.securityContext.capabilities.add
          }
          if (this.securityContextParentObj.securityContext.capabilities.drop) {
            this.form.capabilities.drop = this.securityContextParentObj.securityContext.capabilities.drop
          }
        }

        if (this.securityContextParentObj.securityContext.seLinuxOptions) {
          this.form.seLinuxOptions = this.securityContextParentObj.securityContext.seLinuxOptions
        }
      }
    }
  },
}
</script>