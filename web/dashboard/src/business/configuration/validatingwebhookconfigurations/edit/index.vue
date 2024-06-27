<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'Validatingwebhookconfigurations'}"
                  v-loading="loading">
    <yaml-editor ref="yaml_editor" :is-edit="true"  :value="form"></yaml-editor>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
      <el-button v-loading="loading" @click="onSubmit" type="primary">
        {{ $t("commons.button.submit") }}
      </el-button>
    </div>
  </layout-content>
</template>

<script>
import YamlEditor from "@/components/yaml-editor";
import LayoutContent from "@/components/layout/LayoutContent";
import {updateValidatingwebhookconfiguration, getValidatingwebhookconfiguration} from "@/api/validatingwebhookconfiguration";

export default {
  name: "ValidatingwebhookconfigurationEdit",
  components: {YamlEditor, LayoutContent},
  props: {
    name: String,
  },
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "admissionregistration.k8s.io/v1",
        kind: "ValidatingWebhookConfiguration",
        metadata: {
          name: "",
        },
        provisioner: "",
        parameters: {}
      },
      cluster: ""
    }
  },
  methods: {
    search() {
      getValidatingwebhookconfiguration(this.cluster, this.name).then(res => {
        this.form = res
        if (this.yamlShow) {
          this.onEditYaml()
        }
        this.loading = false
      })
    },
    onSubmit() {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      updateValidatingwebhookconfiguration(this.cluster, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({name: "Validatingwebhookconfigurations"})
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel() {
      this.$router.push({name: "Validatingwebhookconfigurations"})
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    transformYaml () {
      if (this.form?.subjects && this.form?.subjects?.length === 0) {
        delete this.form.subjects
      }
      return JSON.parse(JSON.stringify(this.form))
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.search()
  }
}

</script>

<style scoped>

</style>

