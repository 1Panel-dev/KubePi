<template>
  <div>
    <el-form label-position="top" ref="form" :model="form" :disabled="isReadOnly">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_non_root')" prop="runAsNonRoot">
            <ko-form-item itemType="radio" v-model="form.runAsNonRoot" :radios="non_root_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_user')" prop="runAsUser">
            <ko-form-item itemType="number" placeholder="runAsUser" v-model.number="form.runAsUser" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.run_as_group')" prop="runAsGroup">
            <ko-form-item itemType="number" placeholder="runAsGroup" v-model.number="form.runAsGroup" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="fsGroup" prop="fsGroup">
            <ko-form-item itemType="number" placeholder="fsGroup" v-model.number="form.fsGroup" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span=12>
          <span>{{$t('business.workload.supplemental_groups')}}</span>
          <table style="width: 100%;margin-top:10px" class="tab-table">
            <tr v-for="(row,index) in form.supplementalGroups" :key="index">
              <td width="96%">
                <ko-form-item placeholder="supplementalGroups" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleSupplementalGroupsDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleSupplementalGroupsAdd">{{$t('business.workload.add')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
        <el-col :span=12>
          <span>Sysctls</span>
          <table style="width: 100%;margin-top:10px" class="tab-table">
            <tr v-for="(row,index) in form.sysctls" :key="index">
              <td width="48%">
                <ko-form-item placeholder="e.g. key" itemType="input" v-model="row.name" />
              </td>
              <td width="48%">
                <ko-form-item placeholder="e.g. value" itemType="input" v-model="row.value" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleSysctlsDelete(index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleSysctlsAdd">{{$t('business.workload.add')}} {{$t('business.workload.dns_resolver')}}</el-button>
              </td>
            </tr>
          </table>
        </el-col>
      </el-row>

      <div style="margin-top:20px"><span style="font-size: 16px;">{{$t('business.workload.seLinuxOptions')}}</span></div>
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
      non_root_list: [
        { label: this.$t("business.workload.no"), value: false },
        { label: this.$t("business.workload.non_root"), value: true },
      ],
      form: {
        allowPrivilegeEscalation: null,
        runAsNonRoot: null,
        readOnlyRootFilesystem: null,
        fsGroup: "",
        runAsUser: "",
        runAsGroup: "",
        supplementalGroups: [],
        sysctls: [],
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
    handleSupplementalGroupsDelete(index) {
      this.form.supplementalGroups.splice(index, 1)
    },
    handleSysctlsDelete(index) {
      this.form.sysctls.splice(index, 1)
    },
    handleSupplementalGroupsAdd() {
      var item = {
        value: "",
      }
      this.form.supplementalGroups.push(item)
    },
    handleSysctlsAdd() {
      var item = {
        name: "",
        value: "",
      }
      this.form.sysctls.push(item)
    },

    transformation(parentFrom) {
      if (!parentFrom.securityContext) {
        parentFrom.securityContext = {}
      }
      parentFrom.securityContext.runAsNonRoot = this.form.runAsNonRoot || undefined

      parentFrom.securityContext.runAsUser = this.form.runAsUser || undefined
      parentFrom.securityContext.runAsGroup = this.form.runAsGroup || undefined
      parentFrom.securityContext.fsGroup = this.form.fsGroup || undefined

      let supplementalGroups = []
      for (const s of this.form.supplementalGroups) {
        supplementalGroups.push(s.value)
      }
      parentFrom.securityContext.supplementalGroups = supplementalGroups.length !== 0 ? supplementalGroups : undefined
      parentFrom.securityContext.sysctls = this.form.sysctls.length !== 0 ? this.form.sysctls : undefined

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
        if (this.securityContextParentObj.securityContext.runAsNonRoot !== undefined) {
          this.form.runAsNonRoot = this.securityContextParentObj.securityContext.runAsNonRoot
        }

        if (this.securityContextParentObj.securityContext.runAsUser) {
          this.form.runAsUser = this.securityContextParentObj.securityContext.runAsUser
        }
        if (this.securityContextParentObj.securityContext.runAsGroup) {
          this.form.runAsGroup = this.securityContextParentObj.securityContext.runAsGroup
        }
        if (this.securityContextParentObj.securityContext.fsGroup) {
          this.form.fsGroup = this.securityContextParentObj.securityContext.fsGroup
        }
        if (this.securityContextParentObj.securityContext.supplementalGroups) {
          this.form.supplementalGroups = []
          for (const s of this.securityContextParentObj.securityContext.supplementalGroups) {
            this.form.supplementalGroups.push({ value: s })
          }
        }
        if (this.securityContextParentObj.securityContext.sysctls) {
          this.form.sysctls = this.securityContextParentObj.securityContext.sysctls
        }

        if (this.securityContextParentObj.securityContext.seLinuxOptions) {
          this.form.seLinuxOptions = this.securityContextParentObj.securityContext.seLinuxOptions
        }
      }
    }
  },
}
</script>