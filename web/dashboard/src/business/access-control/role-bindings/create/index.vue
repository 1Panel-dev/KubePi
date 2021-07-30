<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'RoleBindings'}" v-loading="loading">
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
import {createRoleBinding} from "@/api/rolebings"

export default {
  name: "RoleBindingCreate",
  components: { YamlEditor, LayoutContent },
  data () {
    return {
      loading: false,
      cluster: "",
      form: {
        apiVersion: "rbac.authorization.k8s.io/v1",
        kind: "RoleBinding",
        metadata: {
          name: ""
        },
      },
    }
  },
  methods: {
    onSubmit () {
      this.loading = true
      const data = this.$refs.yaml_editor.getValue()
      createRoleBinding(this.cluster, data).then(() => {
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
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
