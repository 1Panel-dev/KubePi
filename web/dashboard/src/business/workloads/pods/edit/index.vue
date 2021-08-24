<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'Pods'}" v-loading="loading">
    <yaml-editor :value="item" :is-edit="true" ref="yaml_editor"></yaml-editor>
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
import {getPodByName, updatePod} from "@/api/pods"

export default {
  name: "PodEdit",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    namespace: String
  },
  data () {
    return {
      loading: false,
      item: {},
      cluster: ""
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getPodByName(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
      }).finally(() => {
        this.loading = false
      })
    },
    onSubmit () {
      this.loading = true
      const data = this.$refs.yaml_editor.getValue()
      updatePod(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Pods" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "Pods" })
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
