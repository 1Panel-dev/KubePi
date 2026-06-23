<template>
  <div>
    <el-form>
      <el-col :span="8">
        <el-form-item :label="$t('business.storage.path_on_node')" required>
          <el-input clearable :placeholder="$t('commons.placeholder.example_data_path')" v-model="spec.local.path"></el-input>
        </el-form-item>
      </el-col>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "Local",
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
      localObj: {
        path: ""
      }
    }
  },
  methods: {},
  created () {
    if (!this.spec.local) {
      this.spec.local = this.localObj
      this.$emit("update:spec", this.spec)
    }
  }
}
</script>

<style scoped>

</style>
