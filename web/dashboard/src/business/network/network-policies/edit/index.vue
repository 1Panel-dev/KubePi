<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'NetworkPolicies'}"
                  v-loading="loading">
    <yaml-editor ref="yaml_editor" :is-edit="true" :value="form"></yaml-editor>
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
import {getNetworkPolicy, updateNetworkPolicy} from "@/api/networkpolicy"

export default {
  name: "NetworkPolicyCreate",
  components: { YamlEditor, LayoutContent },
  props: {
    namespace: String,
    name: String
  },
  data () {
    return {
      loading: false,
      form: {
        apiVersion: "networking.k8s.io/v1",
        kind: "NetworkPolicy",
        metadata: {
          name: "",
          namespace: ""
        },
        spec: {}
      },
      cluster: ""
    }
  },
  methods: {
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      updateNetworkPolicy(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "NetworkPolicies" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "NetworkPolicies" })
    },
    getDetail () {
      this.loading = true
      getNetworkPolicy(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
      }).finally(() => {
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
