<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-col :span="12">
          <ko-form-item :labelName="$t('commons.table.name')" clearable itemType="input"
                        v-model="form.metadata.name"></ko-form-item>
        </el-col>
        <el-col :span="12">
          <ko-form-item :labelName="$t('business.namespace.description')" clearable itemType="input"
                        v-model="form.metadata.annotations['description']"></ko-form-item>
        </el-col>
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
        <yaml-editor :value="yaml"></yaml-editor>
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
import KoFormItem from "@/components/ko-form-item"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import {createNamespace} from "@/api/namespace"

export default {
  name: "NamespaceCreate",
  components: { KoAnnotations, KoLabels, KoFormItem, YamlEditor, LayoutContent },
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
      annotations: {}
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {
      const data = this.transformYaml()
      createNamespace("test1", data).then(res => {
        console.log(res)
      })
    },
    transformYaml() {
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