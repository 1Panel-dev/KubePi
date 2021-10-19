<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{ name: 'NameNodesspaces' }" v-loading="loading">
    <br>
    <el-row :gutter="20">
      <div v-if="!showYaml">
        <el-form label-position="top" :model="item">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')">
              <el-input disabled v-model="item.metadata.name"></el-input>
            </el-form-item>
          </el-col>
        </el-form>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick"
                   v-if="Object.keys(item.metadata).length!==0">
            <el-tab-pane :label="$t('business.workload.labels_annotations')">
              <ko-key-value :title="$t('business.workload.label')"
                            :value.sync="item.metadata.labels"></ko-key-value>
              <ko-key-value :title="$t('business.workload.annotations')"
                            :value.sync="item.metadata.annotations"></ko-key-value>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </div>
      <div v-if="showYaml">
        <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
      </div>
      <div class="bottom-button">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
        <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="loading" @click="onSubmit" type="primary">
          {{ $t("commons.button.submit") }}
        </el-button>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import {getNode, updateNode} from "@/api/nodes"

export default {
  name: "NodeEdit",
  components: { KoKeyValue, LayoutContent, YamlEditor },
  props: {
    name: String,
  },
  data () {
    return {
      item: {
        metadata: {},
        status: {}
      },
      showYaml: false,
      yaml: {},
      activeName: "",
      loading: false
    }
  },
  methods: {
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updateNode(this.cluster, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "Nodes" })
      }).finally(() => {
        this.loading = false
      })
    },
    onEditYaml () {
      this.yaml = this.transformYaml()
      this.showYaml = true
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
    getNodeByName () {
      this.loading = true
      getNode(this.cluster, this.name).then(res => {
        if (res.metadata.resourceVersion) {
          delete res.metadata.resourceVersion
        }
        this.item = res
        this.loading = false
        if (this.showYaml) {
          this.yaml = this.transformYaml()
        }
      })
    },
    onCancel () {
      this.$router.push({ name: "Nodes" })
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    transformYaml () {
      return JSON.parse(JSON.stringify(this.item))
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
    this.getNodeByName()
  }
}
</script>

<style scoped>

</style>
