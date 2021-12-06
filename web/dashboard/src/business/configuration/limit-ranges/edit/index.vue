<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }" v-loading="loading">
    <br>
    <div class="grid-content bg-purple-light" v-if="!showYaml">
      <el-form label-position="top" :model="form" :rules="rules" ref="form">
        <el-row :gutter="20" style="margin-left: 25px;">
          <el-col :span="6">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <el-input disabled clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-form-item disabled itemType="select2" :selections="[]" v-model="form.metadata.namespace" />
            </el-form-item>
          </el-col>
        </el-row>
        <ko-limit-range style="margin-left: 20px;" :limitRangeObj="form.spec.limits.length === 0 ? null : form.spec.limits[0]" ref="ko_limit_range" />
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
import { getLimitRange, updateLimitRange } from "@/api/limitranges"
import KoLimitRange from "@/components/ko-configuration/ko-limit-range.vue"
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "LimitRangeCreate",
  components: { YamlEditor, LayoutContent, KoLimitRange, KoFormItem },
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "v1",
        kind: "LimitRange",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          limits: [],
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
    getDetail() {
      this.loading = true
      getLimitRange(this.cluster, this.namespace, this.name).then((res) => {
        this.form = res
        this.loading = false
        if (this.showYaml) {
          this.yaml = JSON.parse(JSON.stringify(this.form))
        }
      })
    },
    gatherFormData() {
      this.$refs.ko_limit_range.transformation(this.form.spec)
    },
    onSubmit() {
      this.loading = true
      let data
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        this.gatherFormData()
        data = this.form
      }
      updateLimitRange(this.cluster, data.metadata.namespace, data.metadata.name, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.$router.push({ name: "LimitRanges" })
        })
        .finally(() => {
          this.loading = false
        })
    },
    onCancel() {
      this.$router.push({ name: "LimitRanges" })
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
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
