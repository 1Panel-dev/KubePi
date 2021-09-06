<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'StorageClasses'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!yamlShow">
          <el-row :gutter="24" class="row-box">
            <el-col :span="14">
              <el-card class="el-card">
                <ko-detail-basic :item="form" :yaml-show.sync="yamlShow"></ko-detail-basic>
              </el-card>
            </el-col>
            <el-col :span="10">
              <el-card class="el-card">
                <table style="width: 100%" class="myTable">
                  <th scope="col" width="60%" align="left">
                    <h3>{{ $t('commons.form.parameters') }}</h3>
                  </th>
                  <tr v-for="(k,v,index) in form.parameters" :key="index">
                    <td>{{ v }}</td>
                    <td colspan="1">{{ k }}</td>
                  </tr>
                </table>
              </el-card>
            </el-col>
          </el-row>
          <el-tabs style="margin-top: 20px" v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick">
            <el-tab-pane lazy :label="$t('commons.table.resourceInformation')">
              <complex-table :data="[form]">
                <el-table-column :label="$t('business.storage.provisioner')" min-width="30" prop="provisioner" />
                <el-table-column :label="$t('business.storage.reclaimPolicy')" min-width="30" prop="reclaimPolicy" />
                <el-table-column :label="$t('business.storage.volumeBindingMode')" min-width="30" prop="volumeBindingMode" />
              </complex-table>
            </el-tab-pane>
          </el-tabs>
      </div>
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" ref="yaml_editor" :read-only="true"></yaml-editor>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
        </div>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { getStorageClass } from "@/api/storageclass"
import KoDetailBasic from "@/components/detail/detail-basic"
import ComplexTable from "@/components/complex-table"

export default {
  name: "StorageClassDetail",
  components: { YamlEditor, LayoutContent, KoDetailBasic, ComplexTable },
  props: {
    name: String,
  },
  data() {
    return {
      loading: false,
      yamlShow: false,
      form: {
        metadata: {},
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
      this.$router.push({ name: "StorageClasses" })
    },
    onEditYaml() {
      this.yamlShow = true
      this.yaml = this.transformYaml()
    },
    backToForm() {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.yamlShow = false
      })
    },
    search() {
      getStorageClass(this.cluster, this.name)
        .then((res) => {
          this.form = res
          this.yaml = JSON.parse(JSON.stringify(this.form))
        })
        .finally(() => {
          this.loading = false
        })
    },
    transformYaml() {
      let formData = {}
      formData = JSON.parse(JSON.stringify(this.form))
      return formData
    },
    setStorageCapacity() {
      this.form.spec.resources.requests.storage = this.currentStorageCapacity.toString() + "Gi"
    },
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "StorageClassDetail",
        params: { name: this.name },
        query: { yamlShow: newValue },
      })
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.search()
  },
}
</script>
