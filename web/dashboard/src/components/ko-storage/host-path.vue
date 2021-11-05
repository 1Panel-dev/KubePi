<template>
  <div>
    <el-form>
      <el-row :gutter="24">
        <el-col :span="8">
          <el-form-item label="PATH" required>
            <el-input clearable placeholder="eg: /data" v-model="spec.hostPath.path"></el-input>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="Path on the Node" required>
            <el-select v-model="spec.hostPath.type">
              <el-option v-for="(h, index) in hostPathTypes"
                         :key="index"
                         :label="h.label"
                         :value="h.value">
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "HostPath",
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
      hostPath: {
        path: "",
        type: "DirectoryOrCreate"
      },
      hostPathTypes: [{
        value: "DirectoryOrCreate",
        label: this.$t("business.storage.DirectoryOrCreateLabel")
      }, {
        value: "Directory",
        label: this.$t("business.storage.DirectoryLabel")
      }, {
        value: "FileOrCreate",
        label: this.$t("business.storage.FileOrCreateLabel")
      }, {
        value: "File",
        label: this.$t("business.storage.FileLabel")
      }, {
        value: "Socket",
        label: this.$t("business.storage.SocketLabel")
      }, {
        value: "CharDevice",
        label: this.$t("business.storage.CharDeviceLabel")
      }, {
        value: "BlockDevice",
        label: this.$t("business.storage.BlockDeviceLabel")
      }],
    }
  },
  methods: {},
  created () {
    if (!this.spec.hostPath) {
      this.spec.hostPath = this.hostPath
      this.$emit("update:spec",this.spec)
    }
  }
}
</script>

<style scoped>

</style>
