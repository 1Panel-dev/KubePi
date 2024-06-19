<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'Mutatingwebhookconfigurations'}"
                  v-loading="loading">
    <yaml-editor ref="yaml_editor" :value="form"></yaml-editor>
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
import {createMutatingwebhookconfiguration} from "@/api/mutatingwebhookconfiguration";

export default {
  name: "MutatingWebhookConfigurationCreate",
  components: {YamlEditor, LayoutContent},
  props: {},
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "admissionregistration.k8s.io/v1",
        kind: "MutatingWebhookConfiguration",
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
    onSubmit() {
      const data = this.$refs.yaml_editor.getValue()
      this.loading = true
      createMutatingwebhookconfiguration(this.cluster, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
        })
        this.$router.push({name: "Mutatingwebhookconfigurations"})
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel() {
      this.$router.push({name: "Mutatingwebhookconfigurations"})
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
  }
}

</script>

<style scoped>

</style>
