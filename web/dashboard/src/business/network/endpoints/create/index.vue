<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Endpoints'}"
                  v-loading="loading">
    <yaml-editor ref="yaml_editor" :value="form"></yaml-editor>
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
import {createEndPoint} from "@/api/endpoints"

export default {
  name: "EndpointCreate",
  components: { YamlEditor, LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "Endpoints",
        metadata: {
          name: "",
          namespace: ""
        },
      },
      cluster: ""
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Endpoints" })
    },
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      createEndPoint(this.cluster, data.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Endpoints" })
      }).finally(() => {
        this.loading = false
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
