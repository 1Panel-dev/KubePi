<template>
  <layout-content :header="'Deployment - ' + operation + ' - With - Yaml'">
    <div class="grid-content bg-purple-light">
      <yaml-editor :value="yaml" :key="isRefresh"></yaml-editor>
    </div>
    <div class="grid-content bg-purple-light">
      <div style="float: right;margin-top: 10px">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="operation === 'Create'" v-loading="loading" @click="onSubmit()">{{ $t("commons.button.create") }}</el-button>
        <el-button v-if="operation === 'Edit'" v-loading="loading" @click="onUpdate()">{{ $t("commons.button.confirm") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { getDeploymentByName } from "@/api/workloads"

export default {
  name: "DeploymentYaml",
  components: { LayoutContent, YamlEditor },
  data() {
    return {
      loading: false,
      yaml: "",
      isRefresh: false,
      operation: "",
    }
  },
  methods: {
    search() {
      getDeploymentByName(this.$route.params.cluster, this.$route.params.namespace, this.$route.params.name)
        .then((res) => {
          this.yaml = res
          this.isRefresh = !this.isRefresh
        })
    },
    onCancel() {
      this.$router.push({ name: "Deployments" })
    },
  },
  mounted() {
    if (this.$route.params.name) {
      this.search()
      this.operation = "Edit"
    } else {
      this.operation = "Create"
    }
  },
}
</script>
