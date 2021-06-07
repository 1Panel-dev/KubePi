<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{ name: 'Namespaces' }">
    <br>
    <el-row :gutter="20">
      <div class="grid-content bg-purple-light" v-if="!showYaml">
        <el-col :span="12">
          <ko-form-item :labelName="$t('commons.table.name')"  clearable itemType="input" v-model="form.metadata.name"></ko-form-item>
        </el-col>
        <el-col :span="12">
          <ko-form-item :labelName="$t('business.namespace.description')"  clearable itemType="input" v-model="form.metadata.annotations['field.cattle.io/description']"></ko-form-item>
        </el-col>
        <el-col :span="24">
          <br>
          <el-tabs tab-position="left" style="background-color: #141418;">
            <el-tab-pane label="Container Resource Limit">
              <ko-container-resource-limit :resourceLimit.sync = "resourceLimit"></ko-container-resource-limit>
            </el-tab-pane>
            <el-tab-pane label="Labels and Annotations">

            </el-tab-pane>
          </el-tabs>
        </el-col>
        <el-col :span="24" >
          <br>
          <div style="float: right">
            <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
            <el-button @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>
            <el-button v-loading="loading" @click="onSubmit" type="primary">
              {{ $t("commons.button.create") }}
            </el-button>
          </div>
        </el-col>
<!--        <el-form :model="form" label-position="top" :rules="rules" label-width="120px">-->
<!--          <el-col :span="6" :offset="4">-->
<!--            <el-form-item :label="$t('commons.table.name')" prop="name">-->
<!--              <el-input v-model="form.metadata.name" clearable></el-input>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--          <el-col :span="6">-->
<!--            <el-form-item :label="$t('business.namespace.description')" prop="description">-->
<!--              <el-input v-model="form.metadata.annotations['field.cattle.io/description']" clearable></el-input>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--          <el-col :span="24">-->
<!--            <el-tabs tab-position="top" >-->
<!--              <el-tab-pane label="Container Resource Limit">-->
<!--                  <ko-resources></ko-resources>-->
<!--              </el-tab-pane>-->
<!--              <el-tab-pane label="Labels and Annotations">-->

<!--              </el-tab-pane>-->
<!--            </el-tabs>-->
<!--          </el-col>-->
<!--          <el-col :span="6" :offset="10">-->
<!--            <el-form-item style="float: right">-->
<!--              <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>-->
<!--              <el-button @click="onEditYaml()">{{ $t("commons.button.edit_yaml") }}</el-button>-->
<!--              <el-button v-loading="loading" @click="onSubmit" type="primary">-->
<!--                {{ $t("commons.button.create") }}-->
<!--              </el-button>-->
<!--            </el-form-item>-->
<!--          </el-col>-->
<!--        </el-form>-->
      </div>
      <div class="grid-content bg-purple-light" v-if="showYaml">
        <yaml-editor :value="yaml"></yaml-editor>
        <br>
        <div style="float: right">
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.create") }}
          </el-button>
        </div>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import KoFormItem from "@/components/ko-form-item"
import KoContainerResourceLimit from "@/components/ko-namespace/ko-container-resource-limit"

export default {
  name: "NamespaceCreate",
  components: { KoContainerResourceLimit, KoFormItem, YamlEditor, LayoutContent },
  data () {
    return {
      form: {
        apiVersion: "v1",
        kind: "Namespace",
        metadata: {
          name: "",
          annotations: {
            "field.cattle.io/description":"",
            "field.cattle.io/containerDefaultResourceLimit":""
          },
          labels:{
          }
        },
      },
      rules: {},
      loading: false,
      showYaml: false,
      yaml:{},
      resourceLimit:{
        limitsCpu: "",
        limitsMemory: "",
        requestsCpu: "",
        requestsMemory: ""
      }
    }
  },
  methods: {
    onCancel () {
      this.$router.push({ name: "Namespaces" })
    },
    onSubmit () {

    },
    onEditYaml(){
      let formData={};
      formData =JSON.parse(JSON.stringify((this.form)))
      if (formData.metadata.annotations['field.cattle.io/description'] ==="" ){
        delete  formData.metadata.annotations['field.cattle.io/description']
      }
      let resourceData = {}
      resourceData =JSON.parse(JSON.stringify((this.resourceLimit)))
      resourceData.limitsCpu = resourceData.limitsCpu +"m"
      resourceData.limitsMemory = resourceData.limitsMemory +"Mi"
      resourceData.requestsCpu = resourceData.requestsCpu +"Mi"
      resourceData.requestsMemory = resourceData.requestsMemory +"Mi"
      formData.metadata.annotations["field.cattle.io/containerDefaultResourceLimit"] = JSON.stringify(resourceData).replace(/[\\]/g,'')
      this.yaml = formData
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