<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }">
    <el-row :gutter="20" v-if="!showYaml">
      <div class="grid-content bg-purple-light">
        <el-form :model="form" label-position="top" :rules="rules" label-width="120px">
          <el-col :span="6" :offset="4">
            <el-form-item :label="$t('commons.table.name')" prop="name">
              <el-input v-model="form.metadata.name" clearable></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item :label="$t('business.namespace.description')" prop="description">
              <el-input v-model="form.metadata.annotations['field.cattle.io/description']" clearable></el-input>
            </el-form-item>
            <el-form-item style="float: right">
              <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
              <el-button @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
              <el-button v-loading="loading" @click="onSubmit" type="primary">
                {{ $t("commons.button.create") }}
              </el-button>
            </el-form-item>
          </el-col>
        </el-form>
      </div>
    </el-row>
    <el-row v-if="showYaml">
      <yaml-editor :value="form"></yaml-editor>
      <br>
      <div style="float: right">
        <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
        <el-button @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        <el-button v-loading="loading" @click="onSubmit" type="primary">
          {{ $t("commons.button.create") }}
        </el-button>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "NamespaceCreate",
  components: { YamlEditor, LayoutContent },
  data () {
    return {
      form: {
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
          name: "",
          annotations: {
              "field.cattle.io/description":""
          },
          labels:{

          }
        },
      },
      rules: {},
      loading: false,
      showYaml: false,
      yaml:{}
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {

    },
    onEditYaml(){
      let obj={};
      obj=JSON.parse(JSON.stringify(this.form));
      this.yaml = obj
      this.showYaml = true
    },
    backToForm() {
      this.showYaml = false
    }
  }
}
</script>

<style scoped>

</style>