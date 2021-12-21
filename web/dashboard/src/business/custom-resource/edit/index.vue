<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'CustomResourceDefinitions'}" v-loading="loading">
    <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
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
import {getCustomResourceDefinition, updateCustomResourceDefinitione} from "@/api/customresourcedefinitions"

export default {
  name: "CustomResourceDefinitionEdit",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
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
      getCustomResourceDefinition(this.cluster,  this.name).then(res => {
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
      updateCustomResourceDefinitione(this.cluster,this.name, data).then(() => {
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
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
