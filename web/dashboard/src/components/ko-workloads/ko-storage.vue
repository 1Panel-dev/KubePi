<template>
  <div style="margin-top: 20px">
    <ko-card title="Storage">
      <div v-for="(item, index) in volumes" :key="index">
        <el-card style="margin-top: 10px">
          <div slot="header" class="clearfix">
            <span>{{item._type}}</span>
            <el-button style="float: right; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">删 除</el-button>
          </div>
          <el-form label-position="top">
            <div v-if="item._type === 'persistentVolumeClaim'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name" required>
                    <ko-form-item itemType="input" v-model="item.metadata.name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Persistent Volume Claim">
                    <ko-form-item itemType="select2" v-model="item.metadata.weight" :selections="pvc_list" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
            <div v-if="item._type === 'configMap'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name">
                    <ko-form-item itemType="input" v-model="item.metadata.name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Default Mode">
                    <ko-form-item itemType="number" v-model.number="item.defaultMode" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="ConfigMap">
                    <ko-form-item itemType="select2" v-model="item.metadata.configMap.name" :selections="config_map_name_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Optional">
                    <ko-form-item itemType="radio" v-model="item.metadata.configMap.optional" :radios="optional_list" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
            <div v-if="item._type === 'secret'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name">
                    <ko-form-item itemType="input" v-model="item.metadata.name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Default Mode">
                    <ko-form-item itemType="input" v-model="item.defaultMode" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Secret">
                    <ko-form-item itemType="select2" v-model="item.metadata.secret.secretName" :selections="secret_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Optional">
                    <ko-form-item itemType="radio" v-model="item.metadata.secret.optional" :radios="optional_list" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-form>
          <div><span>Valume Mount</span></div>
          <table style="width: 98%;" class="tab-table">
            <tr>
              <th scope="col" width="43%" align="left"><label>mount point</label></th>
              <th scope="col" width="43%" align="left"><label>sub path in volume</label></th>
              <th scope="col" width="8%" align="left"><label>read only</label></th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row, index) in item.volumeMounts" v-bind:key="index">
              <td>
                <ko-form-item itemType="input" v-model="row.mountPath" />
              </td>
              <td>
                <ko-form-item itemType="input" v-model="row.subPath" />
              </td>
              <td>
                <el-checkbox v-model="row.readOnly" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleMountDelete(item, index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleMountAdd(item)">Add Mount</el-button>
              </td>
            </tr>
          </table>
        </el-card>
      </div>
      <el-row>
        <el-col :span="12">
          <el-dropdown split-button @command="handleVolumeAdd">
            Add Volume
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item v-for="(item, index) in volume_type_list" :key="index" :command="item.value">{{item.label}}</el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-col>
      </el-row>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoStorage",
  components: { KoFormItem, KoCard },
  props: {
    storageParentObj: Object,
    currentContainerIndex: Number,
    configMapList: Array,
    secretList: Array,
  },
  watch: {
    currentContainerIndex: {
      handler(newName) {
        this.containerIndex = newName
      },
      immediate: true,
    },
    configMapList: {
      handler(newName) {
        this.config_map_name_list = []
        for (const cm of newName) {
          this.config_map_name_list.push(cm.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
    secretList: {
      handler(newName) {
        this.secret_list = []
        for (const s of newName) {
          this.secret_list.push(s.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      containerIndex: 0,
      volumes: [],
      volumeAdd: [],
      pvc_list: [],
      secret_list: [],
      config_map_name_list: [],
      volume_type_list: [
        { label: "Persistent Volume Claim", value: "persistentVolumeClaim" },
        { label: "ConfigMap", value: "configMap" },
        { label: "Secret", value: "secret" },
      ],
      optional_list: [
        { label: "yes", value: true },
        { label: "no", value: false },
      ],
    }
  },
  methods: {
    handleVolumeAdd(type) {
      let item = {}
      item._type = type
      item.volumeMounts = []
      item.defaultMode = 644
      switch (type) {
        case "configMap":
          item.metadata = {
            name: "",
            configMap: {
              name: "",
              defaultMode: 644,
              optional: true,
            },
          }
          break
        case "secret":
          item.metadata = {
            name: "",
            secret: {
              name: "",
              defaultMode: 644,
              optional: true,
            },
          }
          break
        case "persistentVolumeClaim":
          item.metadata = {
            name: "",
            persistentVolumeClaim: {
              claimName: "",
              readOnly: false,
            },
          }
          break
      }
      this.volumes.push(item)
      this.volumeAdd.push(item)
    },
    handleVolumeDelete(index) {
      for (const vo of this.volumes) {
        for (let i = 0; i < this.volumeAdd; i++) {
          if (vo.metadata.name == this.volumeAdd[i].metadata.name) {
            this.volumeAdd.splice(i, 1)
          }
        }
      }
      this.volumes.splice(index, 1)
    },

    handleMountAdd(item) {
      item.volumeMounts.push({
        mountPath: "",
        subPath: "",
        readOnly: false,
      })
    },
    handleMountDelete(item, index) {
      item.volumeMounts.splice(index, 1)
    },

    checkIsValid() {
      for (const vo of this.volumes) {
        if (vo.metadata.name == "") {
          return false
        }
        for (const mo of vo.volumeMounts) {
          if (mo.mountPath == "") {
            return false
          }
        }
      }
      return true
    },
    transformation(parentFrom) {
      if (!parentFrom.volumes) {
        parentFrom.volumes = []
      }
      parentFrom.containers[this.containerIndex].volumeMounts = []
      for (const volume of this.volumeAdd) {
        switch (volume._type) {
          case "configMap":
            if (volume.defaultMode) {
              volume.configMap.defaultMode = parseInt(volume.defaultMode)
            }
            break
          case "secret":
            if (volume.defaultMode) {
              volume.secret.defaultMode = parseInt(volume.defaultMode)
            }
            break
        }
        parentFrom.volumes.push(volume.metadata)
      }
      for (const volume of this.volumes) {
        if (volume.volumeMounts.length !== 0) {
          for (const mount of volume.volumeMounts) {
            parentFrom.containers[this.containerIndex].volumeMounts.push({ mountPath: mount.mountPath, subPath: mount.subPath, readOnly: mount.readOnly })
          }
        }
      }
    },
  },
  mounted() {
    this.volumes = []
    let volumeMounts = []
    if (this.storageParentObj.containers.length - 1 >= this.containerIndex) {
      if (this.storageParentObj.containers[this.containerIndex].volumeMounts) {
        volumeMounts = this.storageParentObj.containers[this.containerIndex].volumeMounts
      }
    }
    if (this.storageParentObj.volumes) {
      for (const volume of this.storageParentObj.volumes) {
        for (const mount of volumeMounts) {
          if (mount.name === volume.name) {
            let item = {}
            item.volumeMounts = []
            item.metadata = volume
            item.volumeMounts.push(mount)
            if (volume.configMap) {
              if (volume.configMap.defaultMode) {
                item.defaultMode = volume.configMap.defaultMode.toString(8)
              }
              item._type = "configMap"
              this.volumes.push(item)
            }
            if (volume.secret) {
              if (volume.secret.defaultMode) {
                item.defaultMode = volume.secret.defaultMode.toString(8)
              }
              item._type = "secret"
              this.volumes.push(item)
            }
            if (volume.persistentVolumeClaim) {
              item._type = "persistentVolumeClaim"
              this.volumes.push(item)
            }
          }
        }
      }
    }
  },
}
</script>