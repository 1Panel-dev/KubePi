<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'CustomResourceDefinitions'}"
                  v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="15">
          <el-card class="el-card">
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" align="left">
                  <h3>{{ $t("business.common.basic") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("commons.table.name") }}</td>
                <td>{{ item.metadata.name }}</td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.created_time") }}</td>
                <td>{{ item.metadata.creationTimestamp | age }}</td>
              </tr>
            </table>
            <div class="bottom-button">
              <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
            </div>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card>
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" align="left">
                  <h3>{{ $t("business.common.resource") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("business.custom_resource.version") }}</td>
                <td>{{ item.status.storedVersions[0] }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.custom_resource.scope") }}</td>
                <td>{{ item.spec.scope }}</td>
              </tr>
              <tr>
                <td>Group</td>
                <td>{{ item.spec.group }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card" v-if="Object.keys(item.spec).length > 0">
            <el-tab-pane :label=" $t('business.custom_resource.names')">
              <table style="width: 100%" class="myTable">
                <tr>
                  <td>{{ $t("business.custom_resource.singular") }}</td>
                  <td>{{ item.spec.names.singular }}</td>
                </tr>
                <tr>
                  <td>{{ $t("business.custom_resource.plural") }}</td>
                  <td>{{ item.spec.names.plural }}</td>
                </tr>
                <tr>
                  <td>Kind</td>
                  <td>{{ item.spec.names.kind }}</td>
                </tr>
                <tr>
                  <td>List Kind</td>
                  <td>{{ item.spec.names.listKind }}</td>
                </tr>
              </table>
            </el-tab-pane>
            <el-tab-pane :label=" $t('business.custom_resource.version')">
              <complex-table :data="item.spec.versions">
                <el-table-column :label="$t('commons.table.name')" prop="name">
                </el-table-column>
                <el-table-column :label="$t('business.custom_resource.served')" prop="served">
                  <template v-slot:default="{row}">
                    {{ row.served }}
                  </template>
                </el-table-column>
                <el-table-column :label="$t('business.custom_resource.storage')" prop="storage">
                  <template v-slot:default="{row}">
                    {{ row.storage }}
                  </template>
                </el-table-column>
              </complex-table>
            </el-tab-pane>
            <el-tab-pane :label=" $t('business.common.conditions')">
              <ko-detail-conditions :conditions="item.status.conditions"></ko-detail-conditions>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getCustomResourceDefinition} from "@/api/customresourcedefinitions"
import ComplexTable from "@/components/complex-table"
import KoDetailConditions from "@/components/detail/detail-conditions"

export default {
  name: "CustomResourceDefinitionDetail",
  components: { KoDetailConditions, ComplexTable, LayoutContent, YamlEditor, },
  props: {
    name: String
  },
  data () {
    return {
      yaml: {},
      yamlShow: false,
      item: {
        metadata: {},
        spec: {},
        status: {
          storedVersions: []
        }
      },
      loading: false
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getCustomResourceDefinition(this.cluster, this.name).then(res => {
        this.item = res
        this.yaml = JSON.parse(JSON.stringify(this.item))
        this.loading = false
      })
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "CustomResourceDefinitionDetail",
        params: { name: this.name },
        query: { yamlShow: newValue }
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
