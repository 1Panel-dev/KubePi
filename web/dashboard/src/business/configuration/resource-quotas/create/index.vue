<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <div class="grid-content bg-purple-light" v-if="!showYaml">
      <el-form label-position="top" :model="form" :rules="rules" ref="form">
        <el-row :gutter="20" style="margin-left: 25px">
          <el-col :span="4">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="4">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-select style="width:100%" :namespace.sync="form.metadata.namespace"></ko-select>
            </el-form-item>
          </el-col>
        </el-row>
        <ko-resource-quota style="margin-left: 20px" :hasSelector="true" :resourceQuotaObj="form.spec" ref="ko_resource_quota" />
      </el-form>
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
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"
import { createResourceQuota } from "@/api/resourcequota"
import KoSelect from "@/components/ko-select"
import KoResourceQuota from "@/components/ko-configuration/ko-resource-quota.vue"

export default {
  name: "ResourceQuotaCreate",
  components: { YamlEditor, LayoutContent, KoSelect, KoResourceQuota },
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "ResourceQuota",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          hard: {},
          scopeSelector: {
            matchExpressions: [],
          },
        },
      },
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        },
      },

      showYaml: false,
      yaml: undefined,
      cluster: "",
    }
  },
  methods: {
    gatherFormData() {
      this.$refs.ko_resource_quota.transformation(this.form.spec)
    },
    onSubmit() {
      if (this.showYaml) {
        this.onCreate(this.$refs.yaml_editor.getValue())
      } else {
        this.$refs["form"].validate((valid) => {
          if (valid) {
            this.gatherFormData()
            this.onCreate(this.form)
          }
        })
      }
    },
    onCreate(data) {
      this.loading = true
      createResourceQuota(this.cluster, data.metadata.namespace, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "ResourceQuotas" })
        })
        .finally(() => {
          this.loading = false
        })
    },
    onCancel() {
      this.$router.push({ name: "ResourceQuotas" })
    },
    onEditYaml() {
      this.gatherFormData()
      this.yaml = this.form
      this.showYaml = true
    },
    backToForm() {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
  },
}
</script>

<style scoped>
</style>
