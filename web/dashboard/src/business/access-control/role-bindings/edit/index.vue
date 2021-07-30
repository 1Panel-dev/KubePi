<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'RoleBindings'}" v-loading="loading">
    <yaml-editor :value="form" ref="yaml_editor"></yaml-editor>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
      <el-button v-loading="loading" @click="onSubmit" type="primary">
        {{ $t("commons.button.submit") }}
      </el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getRoleBinding, updateRoleBinding} from "@/api/rolebings"


export default {
  name: "RoleBindingEdit",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      loading: false,
      cluster: "",
      form: {}
    }
  },
  methods: {
    onSubmit () {
      this.loading = true
      const data = this.$refs.yaml_editor.getValue()
      updateRoleBinding(this.cluster, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "RoleBindings" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "RoleBindings" })
    },
    getDetail () {
      this.loading = true
      getRoleBinding(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
