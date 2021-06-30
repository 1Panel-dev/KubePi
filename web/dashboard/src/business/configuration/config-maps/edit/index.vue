<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!yamlShow">
          <el-form label-position="top" :model="item">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="item.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="item.metadata.namespace" disabled>
                  <el-option :key="item.metadata.namespace"
                             :label="item.metadata.namespace"
                             :value="item.metadata.namespace"
                  ></el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane label="Data">
                  <ko-data ref="ko_data" :key="keyNum+1" :labelParentObj="item.data"></ko-data>
                </el-tab-pane>
                <el-tab-pane label="Labels/Annotations">
                  <ko-labels ref="ko_labels"  :key="keyNum+2" :labelParentObj="item.metadata"></ko-labels>
                  <ko-annotations ref="ko_annotations"   :key="keyNum+3"  :labelParentObj="item.metadata"></ko-annotations>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="yamlShow">
          <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
        </div>
        <div>
          <div style="float: right;margin-top: 10px">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
            <el-button v-if="!yamlShow" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
            <el-button v-if="yamlShow" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
            <el-button v-loading="loading" @click="onSubmit" type="primary">
              {{ $t("commons.button.submit") }}
            </el-button>
          </div>
        </div>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {updateConfigMap, getConfigMap} from "@/api/configmaps"
import YamlEditor from "@/components/yaml-editor"
import KoData from "@/components/ko-workloads/ko-data"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"

export default {
  name: "ConfigMapEdit",
  components: { LayoutContent, YamlEditor, KoAnnotations, KoLabels, KoData },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      item: {
        metadata: {},
        data: {}
      },
      loading: false,
      yaml: {},
      yamlShow: false,
      activeName: "",
      keyNum: 0,
      cluster: ""
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getConfigMap(this.cluster, this.namespace, this.name).then(res => {
        this.keyNum++
        this.loading = false
        this.item = res
      })
    },
    onSubmit () {
      let data = {}
      if (this.yamlShow) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updateConfigMap(this.cluster, this.namespace, this.item.metadata.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "ConfigMaps" })
      }).finally(() => {
        this.loading = false
      })
    },
    onEditYaml () {
      this.yamlShow = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      this.yamlShow = false
    },
    onCancel () {
      this.$router.push({ name: "ConfigMaps" })
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    transformYaml () {
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.item))
      // labels
      this.$refs.ko_labels.transformation(formData.metadata)
      // annotations
      this.$refs.ko_annotations.transformation(formData.metadata)
      this.$refs.ko_data.transformation(formData)
      return formData
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>