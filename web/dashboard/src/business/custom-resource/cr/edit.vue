<template>
  <layout-content :header="editable=='1'?$t('commons.button.edit'):$t('commons.button.view_yaml')" :back-to="{name: 'CustomResourceDefinitions'}" v-loading="loading">
    <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
    <div class="bottom-button" v-if="editable=='1'">
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
import {getResource, updateResource} from "@/api/customresourcedefinitions"

export default {
  name: "CustomResourceEdit",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    names: String,
    version: String,
    group: String,
    editable: String
  },
  data () {
    return {
      loading: false,
      cluster: "",
      item: {},
      yaml: {}
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getResource(this.cluster, this.version, this.group, this.names, this.namespace, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "CustomResourceDefinitions" })
    },
    onSubmit () {
      this.loading = true
      const data = this.$refs.yaml_editor.getValue()
      updateResource(this.cluster, this.version, this.group, this.names, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "CustomResourceDefinitions" })
      }).finally(() => {
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.namespace = this.$route.query.namespace
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
