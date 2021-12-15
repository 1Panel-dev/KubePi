<template>
  <layout-content :header="$t('commons.button.create')" :back="true">
    <div v-if="yamlValue!=={}">
      <yaml-editor :value="yamlValue" ref="yaml_editor"></yaml-editor>
    </div>
    <div style="float: right;margin-top: 10px">
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
import {getK8sObject} from "@/utils/k8s"
import {postYaml} from "@/api/k8s"

export default {
  name: "YamlCreate",
  components: { YamlEditor, LayoutContent },
  props: {},
  data () {
    return {
      type: "",
      namespace: "",
      yamlValue: {},
      loading: false,
      clusterName: "",
      requires: ["apiVersion", "kind", "metadata"]
    }
  },
  methods: {
    onCancel () {
      this.$router.go(-1)
    },
    onSubmit () {
      const data = this.$refs.yaml_editor.getValue()
      if (!this.checkValid(data)) {
        return
      }
      postYaml(this.clusterName, this.getGroup(), this.yamlValue.apiVersion, this.type, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.go(-1)
      })
    },
    getGroup () {
      return ""
    },
    checkValid (data) {
      for (const key of this.requires) {
        if (data[key] === undefined || data[key] === "") {
          this.$message({
            type: "warning",
            message: key + " is required",
          })
          return false
        }
      }
      return true
    }
  },
  created () {
    this.type = this.$route.query.type
    this.namespace = this.$route.query.namespace
    this.clusterName = this.$route.query.cluster
    this.yamlValue = getK8sObject(this.type, this.namespace)
  }
}
</script>

<style scoped>

</style>
