<template>
  <div>
    <el-form>
      <el-row :gutter="24">
        <el-col :span="8">
          <el-form-item label="PATH" required>
            <el-input clearable placeholder="eg: /nfs-share" v-model="spec.nfs.path"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="Server" required>
            <el-input clearable placeholder="eg: 172.16.10.100" v-model="spec.nfs.server"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="ReadOnly">
            <el-radio v-model="spec.nfs.readOnly" :label="true">Yes</el-radio>
            <el-radio v-model="spec.nfs.readOnly" :label="false">No</el-radio>
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
