<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'HPA'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="item" ref="item">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="item.metadata.name" disabled></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required prop="metadata.namespace">
                <el-select v-model="item.metadata.namespace" disabled>
                  <el-option
                          :key="item.metadata.namespace"
                          :label="item.metadata.name"
                          :value="item.metadata.name">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick" v-if="Object.keys(item.metadata).length!==0">
                <el-tab-pane :label="$t('business.configuration.target')">
                  <ko-hpa-target :namespace="item.metadata.namespace" :cluster="cluster"
                                 :spec-obj.sync="item.spec"></ko-hpa-target>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.configuration.metrics')">
                  <ko-hpa-metrics :metrics-obj.sync="item.spec.metrics"></ko-hpa-metrics>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.workload.labels_annotations')">
                  <ko-key-value :title="$t('business.workload.label')"
                                :value.sync="item.metadata.labels"></ko-key-value>
                  <ko-key-value :title="$t('business.workload.annotations')"
                                :value.sync="item.metadata.annotations"></ko-key-value>
                </el-tab-pane>
              </el-tabs>
            </el-col>
          </el-form>
        </div>
        <div v-if="showYaml">
          <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
        </div>
        <div>
          <div class="bottom-button">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
            <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
            <el-button v-if="showYaml" @click="showYaml=false">{{ $t("commons.button.back_form") }}</el-button>
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
import YamlEditor from "@/components/yaml-editor"
import KoHpaTarget from "@/components/ko-configuration/ko-hpa-target"
import {getHpa, updateHpa} from "@/api/hpa"
import KoHpaMetrics from "@/components/ko-configuration/ko-hpa-metrics"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "HPAEdit",
  components: { KoKeyValue, KoHpaMetrics, KoHpaTarget, LayoutContent, YamlEditor },
  props: {
    namespace: String,
    name: String
  },
  data () {
    return {
      item: {
        metadata: {},
        spec: {}
      },
      loading: false,
      showYaml: false,
      cluster: "",
      activeName: "",
      yaml: {}
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getHpa(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
        if (this.showYaml) {
          this.yaml = this.transformYaml()
        }
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "HPA" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    onSubmit () {
      if (this.showYaml) {
        this.onCreate(this.$refs.yaml_editor.getValue())
      } else {
        this.onCreate(this.transformYaml())
      }
    },
    onCreate (data) {
      this.loading = true
      updateHpa(this.cluster, this.namespace, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "HPA" })
      }).finally(() => {
        this.loading = false
      })
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.item))
    },
  },
  watch: {
    showYaml: function (newValue) {
      this.$router.push({
        name: "HPAEdit",
        params: {
          name: this.name,
          namespace: this.namespace,
        },
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
