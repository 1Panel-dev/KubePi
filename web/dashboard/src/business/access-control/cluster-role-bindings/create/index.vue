<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ClusterRoleBindings'}" v-loading="loading">
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
import {createClusterRoleBinding} from "@/api/clusterrolebindings"

export default {
  name: "ClusterRoleBindingCreate",
  components: { YamlEditor, LayoutContent },
  data () {
    return {
      loading: false,
      cluster: "",
      form: {
        apiVersion: "rbac.authorization.k8s.io/v1",
        kind: "ClusterRoleBinding",
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
      createClusterRoleBinding(this.cluster, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "ClusterRoleBindings" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "ClusterRoleBindings" })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
