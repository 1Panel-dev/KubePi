<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'StorageClasses'}"
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
import {updateStorageClass, getStorageClass} from "@/api/storageclass";

export default {
  name: "StorageClassEdit",
  components: {YamlEditor, LayoutContent},
  props: {
    name: String,
  },
  data() {
    return {
      loading: false,
      form: {
        apiVersion: "storage.k8s.io/v1",
        kind: "StorageClass",
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
      getStorageClass(this.cluster, this.name).then(res => {
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
      updateStorageClass(this.cluster, this.name, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({name: "StorageClasses"})
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel() {
      this.$router.push({name: "StorageClasses"})
    }
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
