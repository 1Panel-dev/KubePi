<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-form label-position="top" :model="form" :rules="rules" ref="form">
          <el-col :span="12">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
        </el-form>
        <el-col :span="24">
          <br>
          <el-tabs v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick">
            <el-tab-pane :label="$t('business.workload.labels_annotations')">
              <ko-key-value :title="$t('business.workload.label')" :value.sync="form.metadata.labels"></ko-key-value>
              <ko-key-value :title="$t('business.workload.annotations')" :value.sync="form.metadata.annotations"></ko-key-value>
            </el-tab-pane>

            <el-tab-pane label="Resource Quotas">
              <ko-resource-quota style="margin-left: 20px" :hasSelector="false" :resourceQuotaObj="quotaForm.spec" ref="ko_resource_quota" />
            </el-tab-pane>

            <el-tab-pane label="Limit Range">
              <ko-limit-range style="margin-left: 20px" :limitRangeObj="rangeForm.spec.limits.length === 0 ? null : rangeForm.spec.limits[0]" ref="ko_limit_range" />
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
import { createNamespace } from "@/api/namespaces"
import Rule from "@/utils/rules"
import KoKeyValue from "@/components/ko-configuration/ko-key-value"
import Bus from "@/utils/bus"
import KoLimitRange from "@/components/ko-configuration/ko-limit-range.vue"
import KoResourceQuota from "@/components/ko-configuration/ko-resource-quota.vue"
import { createResourceQuota } from "@/api/resourcequota"
import { createLimitRange } from "@/api/limitranges"

export default {
  name: "NamespaceCreate",
  components: { KoKeyValue, YamlEditor, LayoutContent, KoLimitRange, KoResourceQuota },
  data() {
    return {
      form: {
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
          name: "",
          annotations: {},
          labels: {},
        },
      },
      quotaForm: {
        apiVersion: "v1",
        kind: "ResourceQuota",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {},
      },
      rangeForm: {
        apiVersion: "v1",
        kind: "LimitRange",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: { limits: [] },
      },
      batchCreateForm: { apiVersion: "v1", kind: "List", items: [] },
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
        },
      },
      loading: false,
      showYaml: false,
      yaml: undefined,
      activeName: "",
      annotations: {},
      cluster: "",
    }
  },
  methods: {
    onCancel() {
      this.$router.push({ name: "Namespaces" })
    },
    gatherFormData() {
      this.$refs.ko_limit_range.transformation(this.rangeForm.spec)
      this.$refs.ko_resource_quota.transformation(this.quotaForm.spec)

      this.batchCreateForm.items = []
      let limitLen = this.rangeForm.spec.limits.length
      let hardLen = Object.keys(this.quotaForm.spec.hard).length
      if (limitLen !== 0 || hardLen !== 0) {
        if (limitLen !== 0) {
          this.rangeForm.metadata.name = this.form.metadata.name + "-limit-range"
          this.rangeForm.metadata.namespace = this.form.metadata.name
          this.batchCreateForm.items.push(this.rangeForm)
        }
        if (hardLen !== 0) {
          this.quotaForm.metadata.name = this.form.metadata.name + "-resource-quota"
          this.quotaForm.metadata.namespace = this.form.metadata.name
          this.batchCreateForm.items.push(this.quotaForm)
        }
      }
      return JSON.parse(JSON.stringify(this.form))
    },
    onSubmit() {
      if (this.showYaml) {
        this.onCreate(this.$refs.yaml_editor.getValue())
      } else {
        this.$refs["form"].validate((valid) => {
          if (valid) {
            let data = this.gatherFormData()
            if (this.batchCreateForm.items.length !== 0) {
              this.batchCreateForm.items.push(data)
              this.onCreate(this.batchCreateForm)
            } else {
              this.onCreate(data)
            }
          }
        })
      }
    },
    onCreate(data) {
      this.loading = true
      let namespaceObj
      let limitrange
      let resourcequota
      if (data.kind === "List") {
        this.batchCreateForm = data
        for (const item of this.batchCreateForm.items) {
          switch (item.kind) {
            case "Namespace":
              namespaceObj = item
              break
            case "LimitRange":
              limitrange = item
              break
            case "ResourceQuota":
              resourcequota = item
              break
          }
        }
      } else {
        namespaceObj = data
      }
      createNamespace(this.cluster, namespaceObj)
        .then(() => {
          if (limitrange) {
            createLimitRange(this.cluster, limitrange.metadata.namespace, limitrange)
          }
          if (resourcequota) {
            createResourceQuota(this.cluster, resourcequota.metadata.namespace, resourcequota)
          }
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "Namespaces" })
          Bus.$emit("refresh", "refresh")
        })
        .finally(() => {
          this.loading = false
        })
    },
    onEditYaml() {
      let data = this.gatherFormData()
      if (this.batchCreateForm.items.length !== 0) {
        this.batchCreateForm.items.push(data)
        this.yaml = this.batchCreateForm
      } else {
        this.yaml = data
      }
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
    handleClick(tab) {
      this.activeName = tab.index
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
