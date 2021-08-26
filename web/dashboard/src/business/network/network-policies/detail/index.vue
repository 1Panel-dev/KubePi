<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'NetworkPolicies'}"
                  v-loading="loading">
    <yaml-editor ref="yaml_editor" :value="form" :read-only="true"></yaml-editor>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.back_detail") }}</el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getNetworkPolicy} from "@/api/networkpolicy"

export default {
  name: "NetworkPolicyCreate",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      loading: false,
      form: {},
      cluster: ""
    }
  },
  methods: {
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
