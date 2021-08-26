<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'ServiceAccounts'}"
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
import {getServiceAccount} from "@/api/serviceaccounts"

export default {
  name: "ServiceAccountDetail",
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
      this.$router.push({ name: "ServiceAccounts" })
    },
    getDetail () {
      this.loading = true
      getServiceAccount(this.cluster, this.namespace, this.name).then(res => {
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
