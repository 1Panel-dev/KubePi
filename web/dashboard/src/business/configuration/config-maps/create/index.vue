<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'ConfigMaps'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <el-row :gutter="20">
        <div v-if="!showYaml">
          <el-form label-position="top" :model="form" :rules="rules" ref="form">
            <el-col :span="6">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="form.metadata.namespace">
                  <el-option v-for="namespace in namespaces"
                             :key="namespace"
                             :label="namespace"
                             :value="namespace">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-tabs v-model="activeName" tab-position="top" type="border-card"
                       @tab-click="handleClick">
                <el-tab-pane :label="$t('business.configuration.data')">
                  <ko-data :data-obj.sync="form.data"></ko-data>
                </el-tab-pane>
                <el-tab-pane :label="$t('business.workload.labels_annotations')">
                  <ko-key-value :title="$t('business.workload.label')"
                                :value.sync="form.metadata.labels"></ko-key-value>
                  <ko-key-value :title="$t('business.workload.labels_annotations')"
                                :value.sync="form.metadata.annotations"></ko-key-value>
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
import KoData from "@/components/ko-workloads/ko-data"
import YamlEditor from "@/components/yaml-editor"
import {createConfigMap} from "@/api/configmaps"
import Rule from "@/utils/rules"
import {getNamespaces} from "@/api/auth"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"

export default {
  name: "ConfigMapCreate",
  components: { KoKeyValue, YamlEditor, KoData, LayoutContent },
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
          namespace: "",
          labels: {},
          annotations: {},
        },
        data: {}
      },
      namespaces: [],
      activeName: "",
      yaml: {},
      cluster: "",
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        }
      }
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
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
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
    getNamespaces(this.cluster).then(res => {
      this.namespaces = res.data
      this.form.metadata.namespace = this.namespaces[0]
    })
  }
}
</script>

<style scoped>

</style>
