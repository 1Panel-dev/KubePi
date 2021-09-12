<template>
  <div>
    <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20" style="margin-top: 10px">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.working_dir')" prop="workingDir">
            <ko-form-item placeholder="e.g. /myapp" itemType="input" v-model="form.workingDir" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('business.workload.stdin')" prop="stdin">
            <ko-form-item itemType="radio" v-model="form.stdin" :radios="stdin_list" />
          </el-form-item>
        </el-col>
        <el-col :span="4">
          <el-form-item label="TTY" prop="tty">
            <el-checkbox :disabled="form.stdin === 'No'" v-model="form.tty">TTY</el-checkbox>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <table style="width: 100%" class="tab-table">
            <tr>
              <th scope="col" width="90%" align="left">
                <label>{{$t('business.workload.commands')}}</label>
              </th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row,index) in form.command" :key="index">
              <td>
                <ko-form-item placeholder="e.g. /bin/sh" itemType="textarea" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleCommandDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleCommandAdd">{{$t('business.workload.add')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
        <el-col :span="12">
          <table style="width: 100%" class="tab-table">
            <tr>
              <th scope="col" width="90%" align="left">
                <label>{{$t('business.workload.arguments')}}</label>
              </th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row,index) in form.args" :key="index">
              <td>
                <ko-form-item placeholder="e.g. /bin/sh" itemType="textarea" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleArgsDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleArgsAdd">{{$t('business.workload.add')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoCommand",
  components: { KoFormItem },
  props: {
    commandParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      form: {
        args: [],
        command: [],
        workingDir: "",
        stdin: null,
        tty: false,
      },
      reFresh: false,
      stdin_list: [
        { label: this.$t("business.workload.no"), value: "No" },
        { label: this.$t("business.workload.once"), value: "Once" },
        { label: this.$t("business.workload.yes"), value: "Yes" },
      ],
    }
  },
  methods: {
    handleCommandAdd() {
      var item = {
        value: "",
      }
      this.form.command.push(item)
    },
    handleCommandDelete(index) {
      this.form.command.splice(index, 1)
    },
    handleArgsAdd() {
      var item = {
        value: "",
      }
      this.form.args.push(item)
    },
    handleArgsDelete(index) {
      this.form.args.splice(index, 1)
    },

    transformation(parentFrom) {
      let commands = []
      for (const cmd of this.form.command) {
        commands.push(cmd.value)
      }
      let args = []
      for (const arg of this.form.args) {
        args.push(arg.value)
      }
      parentFrom.command = commands.length !== 0 ? commands : undefined
      parentFrom.args = args.length !== 0 ? args : undefined
      parentFrom.workingDir = this.form.workingDir || undefined
      if (this.form.stdin) {
        switch (this.form.stdin) {
          case "No":
            parentFrom.stdin = false
            break
          case "Yes":
            parentFrom.stdin = true
            break
          case "Once":
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
        this.form.command = []
        for (const item of this.commandParentObj.command) {
          this.form.command.push({ value: item })
        }
      }
      if (this.commandParentObj.args) {
        this.form.args = []
        for (const item of this.commandParentObj.args) {
          this.form.args.push({ value: item })
        }
      }
      if (this.commandParentObj.workingDir) {
        this.form.workingDir = this.commandParentObj.workingDir
      }
      if (this.commandParentObj.stdinOnce != undefined) {
        this.form.stdinOnce = true
        this.form.stdin = "Once"
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