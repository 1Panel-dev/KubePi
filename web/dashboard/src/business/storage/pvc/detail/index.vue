<template>
  <layout-content :header="$t('commons.form.detail')"  :back-to="{name: 'PersistentVolumeClaim'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!showYaml">
        <el-form label-position="top" :disabled="true" :model="form">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.namespace.namespace')" required>
                <el-select v-model="form.metadata.namespace">
                  <el-option
                      v-for="(item, index) in namespace_list"
                      :key="index"
                      :label="item"
                      :value="item">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="$t('business.storage.pvSource')" required>
                <el-radio v-model="currentVolumeClaimSource" label="sc">{{ this.$t('business.storage.assignSc') }}</el-radio>
                <el-radio v-model="currentVolumeClaimSource" label="pv">{{ this.$t('business.storage.assignPv') }}</el-radio>
              </el-form-item>
            </el-col>
          </el-row>

          <el-tabs v-model="activeName" tab-position="top" type="border-card"
                   @tab-click="handleClick">
            <el-tab-pane :label="$t('commons.form.setting')">
              <div style="margin-top: 20px">
                <ko-card title="Customize">
                  <el-row :gutter="24">
                    <el-col :span="12">
                      <el-form-item :label="$t('business.storage.accessModes')" required>
                        <el-checkbox-group v-model="form.spec.accessModes">
                          <el-checkbox label="ReadWriteOnce">ReadWriteOnce</el-checkbox>
                          <el-checkbox label="ReadOnlyMany">ReadOnlyMany</el-checkbox>
                          <el-checkbox label="ReadWriteMany">ReadWriteMany</el-checkbox>
                        </el-checkbox-group>
                      </el-form-item>
                    </el-col>
                  </el-row>
                </ko-card>
              </div>
              <div style="margin-top: 20px">
                <ko-card title="Volume Claim">
                  <el-row :gutter="24">
                    <el-col :span="6">
                      <el-form-item v-if="currentVolumeClaimSource == 'sc'" :label="$t('business.storage.assignSc')">
                        <el-select v-model="form.spec.storageClassName">
                          <el-option v-for="(sc, index) in storageClasses"
                                     :key="index"
                                     :label="sc"
                                     :value="sc">
                          </el-option>
                        </el-select>
                      </el-form-item>
                      <el-form-item v-if="currentVolumeClaimSource == 'pv'" :label="$t('business.storage.assignPv')">
                        <el-select v-model="form.spec.volumeName">
                          <el-option v-for="(sc, index) in persistentVolumeList"
                                     :key="index"
                                     :disabled="sc.status.phase == 'Bound'"
                                     :label="sc.metadata.name"
                                     :value="sc.metadata.name">
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>
                    <el-col :span="6">
                      <el-form-item v-if="currentVolumeClaimSource == 'sc'" :label="$t('business.storage.capacity')"
                                    required>
                        <el-input-number :min="1" :step="2" @change="setStorageCapacity" clearable
                                         v-model="currentStorageCapacity"></el-input-number>
                      </el-form-item>
                    </el-col>
                  </el-row>
                </ko-card>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-form>
      </div>
      <div v-if="showYaml">
        <yaml-editor :value="yaml" ref="yaml_editor"></yaml-editor>
      </div>
      <div>
        <div style="float: right;margin-top: 10px">
          <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
        </div>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {getPvc} from "@/api/pvc"
import KoCard from "@/components/ko-card/index";

export default {
  name: "PersistentVolumeClaimDetail",
  components: {KoCard, YamlEditor, LayoutContent},
  props: {
    name: String,
    namespace: String
  },
  data() {
    return {
      loading: false,
      showYaml: false,
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
          namespace: "default",
        },
        spec: {
          accessModes: [],
          resources: {
            requests: {
              storage: "1Gi"
            }
          },
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
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    backToForm() {
      this.showYaml = false
    },
    search() {
      getPvc(this.cluster, this.namespace, this.name).then(res => {
        this.form = res
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
  created() {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>