<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'PSPs'}"
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
import {getPSP, updatePSP} from "@/api/podsecuritypolicies"

export default {
  name: "PSPEdit",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
  },
  data () {
    return {
      loading: false,
      form: {
        metadata: {},
        spec: {}
      },
      cluster: ""
    }
  },
  methods: {
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      updatePSP(this.cluster, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "PSPs" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "PSPs" })
    },
    getDetail () {
      getPSP(this.cluster, this.name).then(res => {
        this.form = res
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
