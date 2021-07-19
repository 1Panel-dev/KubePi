<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'HPA'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form" ref="form" :rules="rules">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required prop="metadata.namespace">
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
                <el-tab-pane label="Target">
                  <ko-hpa-target :namespace="form.metadata.namespace" :cluster="cluster"
                                 :spec-obj.sync="form.spec"></ko-hpa-target>
                </el-tab-pane>
                <el-tab-pane label="Metrics">
                  <ko-hpa-metrics :metrics-obj.sync="form.spec.metrics"></ko-hpa-metrics>
                </el-tab-pane>
                <el-tab-pane label="Labels/Annotations">
                  <ko-labels ref="ko_labels" :labelParentObj="form.metadata"></ko-labels>
                  <ko-annotations ref="ko_annotations" :labelParentObj="form.metadata"></ko-annotations>
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
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import YamlEditor from "@/components/yaml-editor"
import KoHpaTarget from "@/components/ko-configuration/ko-hpa-target"
import {listNamespace} from "@/api/namespaces"
import {createHpa} from "@/api/hpa"
import Rule from "@/utils/rules"
import KoHpaMetrics from "@/components/ko-configuration/ko-hpa-metrics"

export default {
  name: "HPACreate",
  components: { KoHpaMetrics, KoHpaTarget, LayoutContent, YamlEditor, KoAnnotations, KoLabels },
  props: {},
  data () {
    return {
      form: {
        apiVersion: "autoscaling/v2beta2",
        kind: "HorizontalPodAutoscaler",
        metadata: {
          namespace: "default"
        },
        spec: {}
      },
      loading: false,
      showYaml: false,
      cluster: "",
      namespaces: [],
      activeName: "",
      yaml: {},
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        }
      }
    }
  },
  methods: {
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
        this.$refs["form"].validate((valid) => {
          if (valid) {
            this.onCreate(this.transformYaml())
          }
        })
      }
    },
    onCreate (data) {
      this.loading = true
      createHpa(this.cluster, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
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
      return JSON.parse(JSON.stringify(this.form))
    }
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
