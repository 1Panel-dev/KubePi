<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="form.metadata.namespace">
                  <el-option v-for="namespace in namespaces"
                             :key="namespace.metadata.name"
                             :label="namespace.metadata.name"
                             :value="namespace.metadata.name">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane label="Data">
                  <ko-data :data-obj.sync="form.data"></ko-data>
                </el-tab-pane>
                <el-tab-pane label="Labels/Annotations">
                  <ko-labels labelTitle="Labels" :label-obj.sync="form.metadata.labels"></ko-labels>
                  <ko-annotations annotations-title="Annotations" :annotations-obj.sync="form.metadata.annotations"></ko-annotations>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
        </div>
        <div>
          <div class="bottom-button">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
            <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
            <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
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
import {listNamespace} from "@/api/namespaces"
import KoData from "@/components/ko-workloads/ko-data"
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import YamlEditor from "@/components/yaml-editor"
import {createConfigMap} from "@/api/configmaps"

export default {
  name: "ConfigMapCreate",
  components: { YamlEditor, KoAnnotations, KoLabels, KoData, LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      showYaml: false,
      form: {
        apiVersion: "v1",
        kind: "ConfigMap",
        metadata: {
          name: "",
          namespace: "default",
          labels: {},
          annotations: {},
        },
        data: {}
      },
      namespaces: [],
      activeName: "",
      yaml: {},
      cluster: ""
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "ConfigMaps" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      this.showYaml = false
    },
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      createConfigMap(this.cluster, this.form.metadata.namespace, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({ name: "ConfigMaps" })
      }).finally(() => {
        this.loading = false
      })
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.form))
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    listNamespace(this.cluster).then(res => {
      this.namespaces = res.items
    })
  }
}
</script>

<style scoped>

</style>
