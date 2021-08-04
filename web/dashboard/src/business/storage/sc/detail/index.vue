<template>
  <layout-content :header="$t('commons.form.detail')"  :back-to="{name: 'StorageClasses'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!yamlShow">
        <el-form label-position="top" :disabled="true" :model="form">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="$t('business.storage.provisioner')" required>
                <el-input clearable v-model="form.provisioner"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="$t('business.storage.reclaimPolicy')" required>
                <el-input clearable v-model="form.reclaimPolicy"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-tabs v-model="activeName" tab-position="top" type="border-card"
                   @tab-click="handleClick">
            <el-tab-pane label="Customize">
              <div style="margin-top: 20px">
                <ko-card :title="$t('commons.form.parameters')">
                  <table style="width: 100%" class="myTable"  >
                    <tr>
                      <th scope="col"></th>
                    </tr>
                    <tr v-for="(k,v,index) in form.parameters" :key="index">
                      <td>{{v}}</td>
                      <td colspan="3">{{ k }}</td>
                    </tr>
                  </table>
                </ko-card>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-form>
      </div>
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
      </div>
      <div>
        <div style="float: right;margin-top: 10px">
          <el-button v-if="!yamlShow" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="yamlShow" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        </div>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getStorageClass} from "@/api/storageclass";
import KoCard from "@/components/ko-card/index";

export default {
  name: "StorageClassDetail",
  components: {KoCard, YamlEditor, LayoutContent},
  props: {
    name: String,
  },
  data() {
    return {
      loading: false,
      yamlShow: false,
      form: {
      },
      activeName: "",
      yaml: {},
      cluster: "",
    }
  },
  methods: {
    handleClick(tab) {
      this.activeName = tab.index
    },
    onCancel() {
      this.$router.push({name: "StorageClasses"})
    },
    onEditYaml() {
      this.yamlShow = true
      this.yaml = this.transformYaml()
    },
    backToForm() {
      this.yamlShow = false
    },
    search() {
      getStorageClass(this.cluster, this.name).then(res => {
        this.form = res
        if (this.yamlShow) {
          this.onEditYaml()
        }
        this.loading = false
      })
    },
    transformYaml() {
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.form))
      return formData
    },
    setStorageCapacity() {
      this.form.spec.resources.requests.storage = this.currentStorageCapacity.toString() + 'Gi'
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.search()
  }
}
</script>