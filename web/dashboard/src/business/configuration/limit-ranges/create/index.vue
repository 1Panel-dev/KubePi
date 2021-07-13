<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'LimitRanges'}"
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
import {createLimitRange} from "@/api/limitranges"

export default {
  name: "ResourceQuotaCreate",
  components: { YamlEditor, LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "LimitRange",
        metadata: {
          name: "",
          namespace: "default"
        },
        spec: {
        }
      },
      cluster: ""
    }
  },
  methods: {
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      createLimitRange(this.cluster, data.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "LimitRanges" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "LimitRanges" })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
