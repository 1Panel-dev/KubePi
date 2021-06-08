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
                        v-model="form.metadata.annotations['field.cattle.io/description']"></ko-form-item>
        </el-col>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="left" style="background-color: #141418;" @tab-click="handleClick">
            <el-tab-pane label="Container Resource Limit">
              <ko-container-resource-limit :resourceLimit.sync="resourceLimit"></ko-container-resource-limit>
            </el-tab-pane>
            <el-tab-pane label="Labels">
              <ko-labels :labelObj.sync="form.metadata.labels"></ko-labels>
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
import KoContainerResourceLimit from "@/components/ko-namespace/ko-container-resource-limit"
import KoLabels from "@/components/ko-workloads/ko-labels"

export default {
  name: "NamespaceCreate",
  components: { KoLabels, KoContainerResourceLimit, KoFormItem, YamlEditor, LayoutContent },
  data () {
    return {
      form: {
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
          name: "",
          annotations: {
            "field.cattle.io/description": "",
            "field.cattle.io/containerDefaultResourceLimit": ""
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
      labels: {},
      activeName: ""
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {

    },
    onEditYaml () {
      let formData = {}
      formData = JSON.parse(JSON.stringify((this.form)))
      if (formData.metadata.annotations["field.cattle.io/description"] === "") {
        delete formData.metadata.annotations["field.cattle.io/description"]
      }
      let resourceData = {}
      resourceData = JSON.parse(JSON.stringify(this.resourceLimit))
      for (const key in resourceData) {
        if (Object.prototype.hasOwnProperty.call(resourceData, key)) {
          if (resourceData[key] === "") {
            delete resourceData[key]
          } else {
            if (key === "limitsCpu") {
              resourceData.limitsCpu = resourceData.limitsCpu + "m"
            } else {
              resourceData[key] = resourceData[key] + "Mi"
            }
          }
        }
      }
      formData.metadata.annotations["field.cattle.io/containerDefaultResourceLimit"] = JSON.stringify(resourceData).replace(/[\\]/g, "")
      this.yaml = formData
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