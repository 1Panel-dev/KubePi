<template>
  <div>
    <el-form>
      <el-row :gutter="24">
        <el-col :span="8">
          <el-form-item :label="$t('business.storage.path')" required>
            <el-input clearable :placeholder="$t('commons.placeholder.example_nfs_path')" v-model="spec.nfs.path"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item :label="$t('business.storage.server')" required>
            <el-input clearable :placeholder="$t('commons.placeholder.example_ip')" v-model="spec.nfs.server"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item :label="$t('business.storage.read_only')">
            <el-radio v-model="spec.nfs.readOnly" :label="true">{{ $t("business.workload.yes") }}</el-radio>
            <el-radio v-model="spec.nfs.readOnly" :label="false">{{ $t("business.workload.no") }}</el-radio>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "pvNfs",
  components: {},
  props: {
    spec: {
      type: Object,
      default: function () {
        return {}
      }
    }
  },
  data () {
    return {
      nfsObj: {
        readOnly: true,
        server: "",
        path: ""
      }
    }
  },
  methods: {},
  created () {
    if (!this.spec.nfs) {
      this.spec.nfs = this.nfsObj
      this.$emit("update:spec", this.spec)
    }
  }
}
</script>

<style scoped>

</style>
