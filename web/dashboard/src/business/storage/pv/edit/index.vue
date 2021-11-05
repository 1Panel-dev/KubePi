<template>
  <layout-content :header="$t('commons.button.edit')" :back-to="{name: 'PersistentVolumes'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!yamlShow">
        <el-form label-position="top" :model="form">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input disabled clearable v-model="form.metadata.name"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="6">
              <el-form-item :label="$t('business.storage.capacity')" required>
                <el-input-number :min="1" @change="setStorageCapacity" clearable
                                 v-model="currentStorageCapacity"></el-input-number>
              </el-form-item>
            </el-col>
            <el-col :span="8">
              <el-form-item :label="$t('commons.table.type')" required>
                <el-select disabled v-model="currentStorageType">
                  <el-option v-for="storage in storageTypes"
                             :key="storage.value"
                             :label="storage.label"
                             :value="storage.value">
                  </el-option>
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-tabs v-model="activeName" tab-position="top" type="border-card"
                   @tab-click="handleClick"  v-if="Object.keys(form.metadata).length !== 0">
            <el-tab-pane label="Customize">
              <div style="margin-top: 20px">
                <ko-card title="Customize">
                  <el-row :gutter="24">
                    <el-col :span="12">
                      <el-form-item :label="$t('business.storage.assignSc')">
                        <el-select clearable v-model="form.spec.storageClassName">
                          <el-option v-for="(sc, index) in storageClasses"
                                     :key="index"
                                     :label="sc"
                                     :value="sc">
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>
                    <el-col :span="12">
                      <el-form-item :label="$t('business.storage.accessModes')" required>
                        <el-radio-group disabled v-model="currentAccessModes">
                          <el-radio @change="setFormAccessModel" label="ReadWriteOnce">ReadWriteOnce</el-radio>
                          <el-radio @change="setFormAccessModel" label="ReadOnlyMany">ReadOnlyMany</el-radio>
                          <el-radio @change="setFormAccessModel" label="ReadWriteMany">ReadWriteMany</el-radio>
                        </el-radio-group>
                      </el-form-item>
                    </el-col>
                  </el-row>
                </ko-card>
                <ko-node-scheduling ref="ko_node_scheduling" :nodeSchedulingType="'matching_rules'"
                                    :nodeList="[]" :nodeSchedulingParentObj="form.spec"/>
              </div>
            </el-tab-pane>
            <el-tab-pane label="Plugin Configuration">
              <div style="margin-top: 20px">
                <ko-card title="Plugin Configuration">
                  <el-form v-if="currentStorageType === 'NFS'">
                    <el-row :gutter="24">
                      <el-col :span="8">
                        <el-form-item label="PATH" required>
                          <el-input disabled clearable placeholder="eg: /nfs-share"
                                    v-model="form.spec.nfs.path"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="8">
                        <el-form-item label="Server" required>
                          <el-input disabled clearable placeholder="eg: 172.16.10.100"
                                    v-model="form.spec.nfs.server"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="6">
                        <el-form-item label="ReadOnly">
                          <el-radio disabled v-model="form.spec.nfs.readOnly" :label="true">Yes</el-radio>
                          <el-radio disabled v-model="form.spec.nfs.readOnly" :label="false">No</el-radio>
                        </el-form-item>
                      </el-col>
                    </el-row>

                  </el-form>
                  <el-form v-if="currentStorageType === 'Local'">
                    <el-col :span="8">
                      <el-form-item label="Path on the Node" required>
                        <el-input disabled clearable placeholder="eg: /data" v-model="form.spec.local.path"></el-input>
                      </el-form-item>
                    </el-col>
                  </el-form>
                  <el-form v-if="currentStorageType === 'Host'">
                    <el-row :gutter="24">
                      <el-col :span="8">
                        <el-form-item label="PATH" required>
                          <el-input disabled clearable placeholder="eg: /data"
                                    v-model="form.spec.hostPath.path"></el-input>
                        </el-form-item>
                      </el-col>
                      <el-col :span="12">
                        <el-form-item label="Path on the Node" required>
                          <el-select disabled v-model="form.spec.hostPath.type">
                            <el-option v-for="(h, index) in hostPathTypes"
                                       :key="index"
                                       :label="h.label"
                                       :value="h.value">
                            </el-option>
                          </el-select>
                        </el-form-item>
                      </el-col>
                    </el-row>
                  </el-form>
                </ko-card>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-form>
      </div>
      <div v-if="yamlShow">
        <yaml-editor :value="yaml" :is-edit="true" ref="yaml_editor"></yaml-editor>
      </div>
      <div>
        <div style="float: right;margin-top: 10px">
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button v-if="!yamlShow" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="yamlShow" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
          <el-button v-loading="loading" @click="onSubmit" type="primary">
            {{ $t("commons.button.submit") }}
          </el-button>
        </div>
      </div>
    </div>
  </layout-content>

</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import {updatePv, getPv} from "@/api/pv"
import KoCard from "@/components/ko-card/index"
import {listStorageClasses} from "@/api/storageclass"
import KoNodeScheduling from "@/components/ko-workloads/ko-node-scheduling.vue"

export default {
  name: "PersistentVolumeEdit",
  props: {
    name: String
  },
  components: { KoCard, YamlEditor, LayoutContent, KoNodeScheduling },
  data () {
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
        kind: "PersistentVolume",
        metadata: {},
        spec: {
          nfs: {
            server: "",
            path: "",
            readOnly: false
          },
          hostPath: {
            path: "",
            type: ""
          },
          local: {
            path: "",
          },
          capacity: {
            storage: ""
          },
          accessModes: [],
          storageClassName: "",
          nodeAffinity: {
            required: {
              nodeSelectorTerms: [
                {
                  matchExpressions: [
                    {
                      key: "",
                      operator: "",
                      values: [],
                    },
                  ],
                },
              ],
            },
          },
        }
      },
      activeName: "",
      yaml: {},
      cluster: "",
      storageTypes: [{
        value: "NFS",
        label: "NFS Share"
      }, {
        value: "Local",
        label: "Local Volume"
      }, {
        value: "Host",
        label: "Host Path"
      }],
      hostPathTypes: [{
        value: "DirectoryOrCreate",
        label: this.$t("business.storage.DirectoryOrCreateLabel")
      }, {
        value: "Directory",
        label: this.$t("business.storage.DirectoryLabel")
      }, {
        value: "FileOrCreate",
        label: this.$t("business.storage.FileOrCreateLabel")
      }, {
        value: "File",
        label: this.$t("business.storage.FileLabel")
      }, {
        value: "Socket",
        label: this.$t("business.storage.SocketLabel")
      }, {
        value: "CharDevice",
        label: this.$t("business.storage.CharDeviceLabel")
      }, {
        value: "BlockDevice",
        label: this.$t("business.storage.BlockDeviceLabel")
      }],
      currentStorageType: "",
      currentAccessModes: "ReadWriteOnce",
      currentStorageCapacity: 1,
      storageClasses: []
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getPv(this.cluster, this.name).then(res => {
        this.form = res
        if (this.form.spec.hostPath) {
          this.currentStorageType = "Host"
        }
        if (this.form.spec.local) {
          this.currentStorageType = "Local"
        }
        if (this.form.spec.nfs) {
          this.currentStorageType = "NFS"
        }
        this.loading = false
        if (this.yamlShow) {
          this.yaml = this.transformYaml()
        }
        this.currentStorageCapacity = this.form.spec.capacity.storage.replace(/[^0-9]/ig, "")
      })
    },
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "PersistentVolumes" })
    },
    onEditYaml () {
      this.yamlShow = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.yamlShow = false
      })
    },
    onSubmit () {
      let data = {}
      if (this.yamlShow) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      this.loading = true
      updatePv(this.cluster, data, this.name).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.update_success"),
        })
        this.$router.push({ name: "PersistentVolumes" })
      }).finally(() => {
        this.loading = false
      })
    },
    loadStorageClasses () {
      this.storageClasses = []
      listStorageClasses(this.cluster).then((res) => {
        for (const sc of res.items) {
          this.storageClasses.push(sc.metadata.name)
        }
      })
    },
    transformYaml () {
      this.form.spec["storageClassName"] = this.form.spec["storageClassName"] || undefined
      switch (this.currentStorageType) {
        case "Local":
          delete this.form.spec["hostPath"]
          delete this.form.spec["nfs"]
          break
        case "Host":
          delete this.form.spec["local"]
          delete this.form.spec["nfs"]
          delete this.form.spec["storageClassName"]
          break
        case "NFS":
          delete this.form.spec["hostPath"]
          delete this.form.spec["local"]
          if (this.form.spec.nfs.readOnly) {
            delete this.form.spec.nfs["readOnly"]
          }
          break
      }
      if (this.form.spec.nodeAffinity && this.$refs.ko_node_scheduling) {
        this.$refs.ko_node_scheduling.transformation(this.form.spec)
      }
      return JSON.parse(JSON.stringify(this.form))
    },
    setFormAccessModel () {
      this.form.spec.accessModes = []
      this.form.spec.accessModes.push(this.currentAccessModes)
    },
    setStorageCapacity () {
      this.form.spec.capacity.storage = this.currentStorageCapacity.toString() + "Gi"
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.loadStorageClasses()
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
