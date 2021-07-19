<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-form label-position="top" :model="form" :rules="rules" ref="form">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')" required prop="metadata.name">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
        </el-form>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick">
            <el-tab-pane label="Labels/Annotations">
              <ko-labels labelTitle="Labels" :label-obj.sync="form.metadata.labels"></ko-labels>
              <ko-annotations annotations-title="Annotations"
                              :annotations-obj.sync="form.metadata.annotations"></ko-annotations>
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
          <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.submit") }}
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
import {createNamespace} from "@/api/namespaces"
import Rule from "@/utils/rules"

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
          annotations: {},
          labels: {}
        },
      },
      rules: {
        metadata: {
          name: [Rule.RequiredRule]
        }
      },
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
      cluster: ""
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
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
      createNamespace(this.cluster, data).then(() => {
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
      return JSON.parse(JSON.stringify(this.form))
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
  },
  created () {
    this.cluster = this.$route.query.cluster
  }
}
</script>

<style scoped>

</style>
