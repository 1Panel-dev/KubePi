<template>
  <layout-content :header="$t('commons.button.create')" :back-to="{name: 'PersistentVolumes'}" v-loading="loading">
    <div class="grid-content bg-purple-light">
      <div v-if="!showYaml">
        <el-form label-position="top" :model="form">
          <el-row :gutter="24">
            <el-col :span="8">
              <el-form-item :label="$t('commons.table.name')" required>
                <el-input clearable v-model="form.metadata.name"></el-input>
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
                <el-select v-model="currentStorageType">
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
                   @tab-click="handleClick">
            <el-tab-pane label="Customize">
              <div style="margin-top: 20px">
                <ko-card title="Customize">
                  <el-row :gutter="24">
                    <el-col :span="12">
                      <el-form-item :label="$t('business.storage.assignSc')">
                        <el-select v-model="form.spec.storageClassName">
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
                        <el-checkbox-group v-model="form.spec.accessModes">
                          <el-checkbox label="ReadWriteOnce">ReadWriteOnce</el-checkbox>
                          <el-checkbox label="ReadOnlyMany">ReadOnlyMany</el-checkbox>
                          <el-checkbox label="ReadWriteMany">ReadWriteMany</el-checkbox>
                        </el-checkbox-group>
                      </el-form-item>
                    </el-col>
                  </el-row>
                </ko-card>
                <ko-node-scheduling :isReadOnly="readOnly" ref="ko_node_scheduling"
                                    :nodeSchedulingType="'matching_rules'" :nodeList="[]"
                                    :nodeSchedulingParentObj="form.spec"/>
              </div>
            </el-tab-pane>
            <el-tab-pane label="Plugin Configuration">
              <div style="margin-top: 20px">
                <ko-card title="Plugin Configuration">
                  <pv-nfs v-if="currentStorageType === 'NFS'" :spec.sync="form.spec"></pv-nfs>
                  <local v-if="currentStorageType === 'Local'" :spec.sync="form.spec"></local>
                  <host-path v-if="currentStorageType === 'Host'" :spec.sync="form.spec"></host-path>
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
          <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
          <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
          <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
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
import {createPv} from "@/api/pv"
import KoCard from "@/components/ko-card/index"
import {listStorageClasses} from "@/api/storageclass"
import KoNodeScheduling from "@/components/ko-workloads/ko-node-scheduling.vue"
import {checkPermissions} from "@/utils/permission"
import PvNfs from "@/components/ko-storage/pv-nfs"
import Local from "@/components/ko-storage/local"
import HostPath from "@/components/ko-storage/host-path"

export default {
  name: "PersistentVolumeCreate",
  components: { HostPath, Local, PvNfs, KoCard, YamlEditor, LayoutContent, KoNodeScheduling },
  data () {
    return {
      loading: false,
      showYaml: false,
      page: {
        pageSize: 10,
        nextToken: "",
      },
      conditions: "",
      readOnly: false,
      form: {
        apiVersion: "v1",
        kind: "PersistentVolume",
        metadata: {
          name: "",
        },
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
            storage: "1Gi"
          },
          nodeAffinity: {
            required: {
              nodeSelectorTerms: [],
            },
          },
          accessModes: []
        }
      },
      namespaces: [],
      activeName: "",
      yaml: undefined,
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
      currentStorageType: "Local",
      currentStorageCapacity: 1,
      storageClasses: []
    }
  },
  methods: {
    handleClick (tab) {
      this.activeName = tab.index
    },
    onCancel () {
      this.$router.push({ name: "PersistentVolumes" })
    },
    onEditYaml () {
      this.showYaml = true
      this.yaml = this.transformYaml()
    },
    backToForm () {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.continue_edit"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    onSubmit () {
      let data = {}
      if (this.showYaml) {
        data = this.$refs.yaml_editor.getValue()
      } else {
        data = this.transformYaml()
      }
      if (!data.spec.nodeAffinity?.required?.nodeSelectorTerms?.length > 0) {
        delete data.spec.nodeAffinity
      }
      this.loading = true
      createPv(this.cluster, data).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success"),
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
      let formData = {}
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
          // if (this.form.spec.nfs.readOnly) {
          //   delete this.form.spec.nfs["readOnly"]
          // }
          break
      }

      this.$refs.ko_node_scheduling.transformation(this.form.spec)
      // if (this.form.spec.nodeAffinity?.required.nodeSelectorTerms.length === 0) {
      //   delete this.form.spec.nodeAffinity.required.nodeSelectorTerms
      // }
      formData = JSON.parse(JSON.stringify(this.form))
      return formData
    },
    setStorageCapacity () {
      this.form.spec.capacity.storage = this.currentStorageCapacity.toString() + "Gi"
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.showYaml = this.$route.query.yamlShow === "true"
    if (checkPermissions({ scope: "cluster", apiGroup: "storage.k8s.io", resource: "storageclasses", verb: "list" })) {
      this.loadStorageClasses()
    }
  }
}
</script>

<style scoped>

</style>
