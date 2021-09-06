<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'PersistentVolumeClaim'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!yamlShow">
        <el-form label-position="top" :model="form">
          <el-row :gutter="20" class="row-box">
            <el-col :span="24">
              <el-card class="el-card">
                <ko-detail-basic :item="form" :yaml-show.sync="yamlShow"></ko-detail-basic>
              </el-card>
            </el-col>
          </el-row>
          <el-tabs style="margin-top: 20px" v-model="activeName" tab-position="top" type="border-card"
                   @tab-click="handleClick">
            <el-tab-pane lazy :label="$t('commons.table.resourceInformation')">
              <complex-table :data="[form.status]">
                <el-table-column :label="$t('commons.table.status')" min-width="30">
                  <template v-slot:default="{row}">
                    <el-button v-if="row.phase && row.phase === 'Bound'" type="success" size="mini" plain round>
                      {{ row.phase }}
                    </el-button>
                    <el-button v-if="row.phase && row.phase === 'Available'" type="success" size="mini" plain round>
                      {{ row.phase }}
                    </el-button>
                    <el-button v-if="row.phase && row.phase === 'Pending'" type="warning" size="mini" plain round>
                      {{ row.phase }}
                    </el-button>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('business.storage.capacity')" min-width="30"
                                 prop="capacity.storage"/>
                <el-table-column :label="$t('business.storage.accessModes')" min-width="30">
                  <template v-slot:default="{row}">
                    <div v-for="(name,index) in row.accessModes " :key="index" style="display:inline-block">
                      <el-tag>{{ name }}</el-tag>
                    </div>
                  </template>
                </el-table-column>
              </complex-table>
            </el-tab-pane>
          </el-tabs>
        </el-form>
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
import {getPvc} from "@/api/pvc"
import KoDetailBasic from "@/components/detail/detail-basic";
import ComplexTable from "@/components/complex-table"

export default {
  name: "PersistentVolumeClaimDetail",
  components: {KoDetailBasic, ComplexTable, YamlEditor, LayoutContent},
  props: {
    name: String,
    namespace: String
  },
  data() {
    return {
      loading: false,
      yamlShow: false,
      page: {
        pageSize: 10,
        nextToken: "",
      },
      conditions: "",
      form: {
        apiVersion: "v1",
        kind: "PersistentVolumeClaim",
        metadata: {
          name: "",
          namespace: "",
        },
        spec: {
          accessModes: [],
          resources: {
            requests: {
              storage: ""
            }
          },
        },
        status: {
          phase: "",
          accessModes: [],
          capacity: []
        }
      },
      namespace_list: [],
      activeName: "",
      yaml: {},
      cluster: "",
      persistentVolumeList: [],
      currentStorageCapacity: 1,
      storageClasses: [],
      currentVolumeClaimSource: "sc"
    }
  },
  methods: {
    handleClick(tab) {
      this.activeName = tab.index
    },
    onCancel() {
      this.$router.push({name: "PersistentVolumes"})
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
      getPvc(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
        this.yaml = JSON.parse(JSON.stringify(this.form))
        this.currentStorageCapacity = parseInt(this.form.spec.resources.requests.storage)
      }).finally(() => {
        this.loading = false
      })
    },
    transformYaml() {
      let formData = {}
      if (this.form.spec.storageClassName == 'Default StorageClass') {
        this.form.spec.storageClassName = ''
      }
      formData = JSON.parse(JSON.stringify(this.form))
      return formData
    },
    setStorageCapacity() {
      this.form.spec.resources.requests.storage = this.currentStorageCapacity.toString() + 'Gi'
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "PersistentVolumeClaimDetail",
        params: {name: this.name, namespace: this.namespace},
        query: {yamlShow: newValue}
      })
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.search()
  }
}
</script>
