<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-form label-position="top" :model="item">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')">
              <el-input disabled v-model="item.metadata.name"></el-input>
            </el-form-item>
          </el-col>
        </el-form>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick">
            <el-tab-pane label="Labels">
              <ko-labels ref="ko_labels" :key="loading" :labelParentObj="item.metadata"></ko-labels>
            </el-tab-pane>
            <el-tab-pane label="Annotations">
              <ko-annotations ref="ko_annotations"  :key="loading" :annotationsParentObj="item.metadata.annotations"></ko-annotations>
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
import KoLabels from "@/components/ko-workloads/ko-labels"
import KoAnnotations from "@/components/ko-workloads/ko-annotations"
import YamlEditor from "@/components/yaml-editor"
import {getNamespace, updateNamespace} from "@/api/namespaces"

export default {
  name: "NamespaceEdit",
  components: { LayoutContent,KoAnnotations, KoLabels ,YamlEditor},
  props: {
    name: String,
  },
  data () {
    return {
      item: {
        metadata: {
        },
        status: {}
      },
      showYaml: false,
      yaml: {},
      activeName: "",
      loading: false
    }
  },
  methods: {
    onSubmit() {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updateNamespace(this.cluster,this.item.metadata.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Namespaces" })
      }).finally(() => {
        this.loading = false
      })
    },
    onEditYaml() {
      this.yaml = this.transformYaml()
      this.showYaml = true
    },
    backToForm() {
      this.showYaml = false
    },
    getNamespaceByName () {
      this.loading = true
      getNamespace(this.cluster, this.name).then(res => {
        this.item = res
        this.loading = false
        if(this.showYaml) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
      })
    },
    onCancel(){
      this.$router.push({ name: "Namespaces" })
    },
    handleClick(tab) {
      this.activeName = tab.index
    },
    transformYaml () {
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.item))
      // labels
      this.$refs.ko_labels.transformation(formData.metadata)
      // annotations
      this.$refs.ko_annotations.transformation(formData.metadata)
      return formData
    },
  },
  watch: {
    showYaml:function (newValue) {
      this.$router.push({
        path: "/namespaces/edit/"+this.name ,
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow  === "true"
    this.getNamespaceByName()
  }
}
</script>

<style scoped>

</style>