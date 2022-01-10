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
import {getCluster} from "@/api/clusters"

export default {
  name: "YamlCreate",
  components: { YamlEditor, LayoutContent },
  props: {},
  data () {
    return {
      type: "",
      yamlValue: {},
      namespace: "",
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
      postYaml(this.clusterName,this.type, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.go(-1)
      })
    },
    checkValid (data) {
      let result = true
      let message = ""
      for (const key of this.requires) {
        if (data[key] === undefined || data[key] === null || data[key] === "") {
          message = key + " is required \n"
          result = false
          break
        }
      }
      if (data.metadata.name === undefined || data.metadata.name === "" || data.metadata.name === null) {
        message = message + "name is required \n"
        result = false
      }
      if (data.metadata.namespace !== undefined && (data.metadata.namespace === "" || data.metadata.namespace === null)) {
        message = message + "namespace is required \n"
        result = false
      }
      if (!result) {
        this.$message({
          type: "warning",
          message: message,
        })
      }
      return result
    }
  },
  created () {
    this.type = this.$route.query.type
    this.clusterName = this.$route.query.cluster
    this.namespace = sessionStorage.getItem("namespace")?sessionStorage.getItem("namespace"):"default"
    getCluster(this.clusterName).then(res => {
      this.yamlValue = getK8sObject(this.type, this.namespace, res.data.status.version)
    })
  }
}
</script>

<style scoped>

</style>
