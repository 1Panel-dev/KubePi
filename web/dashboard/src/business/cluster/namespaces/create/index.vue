<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-form label-position="top" :model="form">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.namespace.description')">
              <el-input clearable v-model="form.metadata.annotations['description']"></el-input>
            </el-form-item>
          </el-col>
        </el-form>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="left" style="background-color: #141418;" @tab-click="handleClick">
            <el-tab-pane label="Labels">
              <ko-labels ref="ko_labels" :labelParentObj="form.metadata"></ko-labels>
            </el-tab-pane>
            <el-tab-pane label="Annotations">
              <ko-annotations ref="ko_annotations" :annotationsParentObj="annotations"></ko-annotations>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </div>
      <div class="grid-content bg-purple-light" v-if="showYaml">
        <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
      </div>
      <div class="grid-content bg-purple-light">
        <div style="float: right;margin-top: 10px">
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
          <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.create") }}
          </el-button>
        </div>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import {createNamespace} from "@/api/namespace"

export default {
  name: "NamespaceCreate",
  components: { KoAnnotations, KoLabels, YamlEditor, LayoutContent },
  data () {
    return {
      form: {
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
          name: "",
          annotations: {
            "description": ""
          },
          labels: {}
        },
      },
      rules: {},
      loading: false,
      showYaml: false,
      yaml: {},
      resourceLimit: {
        limitsCpu: "",
        limitsMemory: "",
        requestsCpu: "",
        requestsMemory: ""
      },
      activeName: "",
      annotations: {},
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      createNamespace(this.clusterName, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "Namespaces" })
      }).finally(() => {
        this.loading = false
      })
    },
    transformYaml () {
      let formData = {}
      formData = JSON.parse(JSON.stringify((this.form)))
      if (formData.metadata.annotations["description"] === "") {
        delete formData.metadata.annotations["description"]
      }
      // labels
      this.$refs.ko_labels.transformation(formData.metadata)
      // annotations
      this.$refs.ko_annotations.transformation(formData.metadata)
      return formData
    },
    onEditYaml () {
      this.yaml = this.transformYaml()
      this.showYaml = true
    },
    backToForm () {
      this.showYaml = false
    },
    handleClick (tab) {
      this.activeName = tab.index
    }
  }
}
</script>

<style scoped>

</style>