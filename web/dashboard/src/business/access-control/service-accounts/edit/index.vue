<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'ServiceAccounts'}"
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
import {getServiceAccount, updateServiceAccount} from "@/api/serviceaccounts"

export default {
  name: "ServiceAccountEdit",
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
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      updateServiceAccount(this.cluster,this.namespace,this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "ServiceAccounts" })
      }).finally(() => {
        this.loading = false
      })
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
