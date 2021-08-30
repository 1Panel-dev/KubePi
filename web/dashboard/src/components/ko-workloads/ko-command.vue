<template>
  <div style="margin-top: 20px">
    <ko-card :key="reFresh" :title="$t('business.workload.command')">
      <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.entry_point')" prop="command">
              <ko-form-item placeholder="e.g. /bin/sh" itemType="input" v-model="form.command" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.arguments')" prop="args">
              <ko-form-item placeholder="e.g. /usr/sbin/httpd -f httpd.conf" itemType="textarea" v-model="form.args" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.workload.working_dir')" prop="workingDir">
              <ko-form-item placeholder="e.g. /myapp" itemType="input" v-model="form.workingDir" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.workload.stdin')" prop="stdin">
              <ko-form-item itemType="radio" v-model="form.stdin" :radios="stdin_list" />
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item label="TTY" prop="tty">
              <el-checkbox :disabled="form.stdin === 'No'" v-model="form.tty">TTY</el-checkbox>
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
  name: "KoCommand",
  components: { KoFormItem, KoCard },
  props: {
    commandParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      form: {
        args: "",
        command: "",
        workingDir: "",
        stdin: null,
        tty: false,
      },
      reFresh: false,
      stdin_list: [
        { label: this.$t("business.workload.no"), value: "No" },
        { label: this.$t("business.workload.once"), value: "Ones" },
        { label: this.$t("business.workload.yes"), value: "Yes" },
      ],
    }
  },
  methods: {
    transformation(parentFrom) {
      parentFrom.command = this.form.command ? this.form.command.split(",") : undefined
      parentFrom.args = this.form.args ? this.form.args.split(",") : undefined
      parentFrom.workingDir = this.form.workingDir || undefined
      if (this.form.stdin) {
        switch (this.form.stdin) {
          case this.$t("business.workload.no"):
            parentFrom.stdin = false
            break
          case this.$t("business.workload.yes"):
            parentFrom.stdin = true
            break
          case this.$t("business.workload.once"):
            parentFrom.stdin = true
            parentFrom.stdinOnce = true
            break
        }
      } else {
        delete parentFrom.stdin
        delete parentFrom.stdinOnce
      }
      parentFrom.tty = this.form.tty || false
    },
  },
  mounted() {
    if (this.commandParentObj) {
      if (this.commandParentObj.command) {
        this.form.command = this.commandParentObj.command.join(",")
      }
      if (this.commandParentObj.args) {
        this.form.args = this.commandParentObj.args.join(",")
      }
      if (this.commandParentObj.workingDir) {
        this.form.workingDir = this.commandParentObj.workingDir
      }
      if (this.commandParentObj.stdinOnce != undefined) {
        this.form.stdinOnce = true
        this.form.stdin = "Ones"
      }
      if (this.commandParentObj.stdin != undefined) {
        this.form.stdin = this.commandParentObj.stdin ? "Yes" : "No"
      }
      if (this.commandParentObj.tty) {
        this.form.tty = this.commandParentObj.tty
      }
    }
  },
}
</script>