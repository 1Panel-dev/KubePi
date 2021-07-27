<template>
  <div style="margin-top: 20px">
    <ko-card title="Storage">
      <el-form label-position="top" :disabled="isReadOnly">
        <div v-for="(item, index) in volumes" :key="index">
          <el-card style="margin-top: 10px">
            <div slot="header" class="clearfix">
              <span>{{item.type}}</span>
              <el-button style="float: right; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">删 除</el-button>
            </div>
            <div v-if="item.type === 'persistentVolumeClaim'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name" required>
                    <ko-form-item itemType="input" v-model="item.name" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Persistent Volume Claim">
                    <ko-form-item itemType="select2" v-model="item.resource" :selections="pvc_list" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-checkbox style="margin-bottom: 20px" v-model="item.readOnly" >Read Only</el-checkbox>
            </div>
            <div v-if="item.type === 'configMap'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name">
                    <ko-form-item itemType="input" v-model="item.name" />
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
                    <ko-form-item itemType="select2" v-model="item.resource" :selections="config_map_name_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Optional">
                    <ko-form-item itemType="radio" v-model="item.optional" :radios="optional_list" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
            <div v-if="item.type === 'secret'">
              <el-row :gutter="20">
                <el-col :span="12">
                  <el-form-item label="Volume Name">
                    <ko-form-item itemType="input" v-model="item.name" />
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
                    <ko-form-item itemType="select2" v-model="item.resource" :selections="secret_list" />
                  </el-form-item>
                </el-col>
                <el-col :span="12">
                  <el-form-item label="Optional">
                    <ko-form-item itemType="radio" v-model="item.optional" :radios="optional_list" />
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
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
      </el-form>
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
    isReadOnly: Boolean,
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
      this.volumes.push({
        type: type,
        name: "",
        resource: "",
        defaultMode: 644,
        optional: null,
        volumeMounts: [],
      })
    },
    handleVolumeDelete(index) {
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
        if (vo.name == "") {
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
      for (let i = 0; i < parentFrom.volumes.length; i++) {
        if (parentFrom.volumes[i].configMap || parentFrom.volumes[i].secret || parentFrom.volumes[i].persistentVolumeClaim) {
          parentFrom.volumes.splice(i, 1)
        }
      }
      parentFrom.containers[this.containerIndex].volumeMounts = []
      for (const volume of this.volumes) {
        let item = {}
        if (volume.name) {
          item.name = volume.name
        }
        switch (volume.type) {
          case "configMap":
            item.configMap = {}
            if (volume.defaultMode) {
              item.configMap.defaultMode = parseInt("0" + volume.defaultMode.toString(), 8)
            }
            if (volume.resource) {
              item.configMap.name = volume.resource
            }
            if (volume.optional !== null) {
              item.configMap.optional = volume.optional
            }
            break
          case "secret":
            item.secret = {}
            if (volume.defaultMode) {
              item.secret.defaultMode = parseInt("0" + volume.defaultMode.toString(), 8)
            }
            if (volume.resource) {
              item.secret.secretName = volume.resource
            }
            if (volume.optional !== null) {
              item.secret.optional = volume.optional
            }
            break
          case "persistentVolumeClaim":
            item.persistentVolumeClaim = {}
            if (volume.resource) {
              item.persistentVolumeClaim.name = volume.resource
            }
            if (volume.readOnly !== null) {
              item.persistentVolumeClaim.readOnly = volume.readOnly
            }
            break
        }
        parentFrom.volumes.push(item)
        if (volume.volumeMounts.length !== 0) {
          for (const mount of volume.volumeMounts) {
            parentFrom.containers[this.containerIndex].volumeMounts.push({ name: item.name, mountPath: mount.mountPath, subPath: mount.subPath, readOnly: mount.readOnly })
          }
        } else {
          parentFrom.containers[this.containerIndex].volumeMounts.push({ name: item.name })
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
        if (volume.configMap || volume.secret || volume.persistentVolumeClaim) {
          let item = {}
          item.volumeMounts = []
          for (const mount of volumeMounts) {
            if (mount.name === volume.name) {
              item.volumeMounts.push(mount)
            }
            if (volume.name) {
              item.name = volume.name
            }
            if (volume.configMap) {
              item.type = "configMap"
              if (volume.configMap.defaultMode) {
                item.defaultMode = volume.configMap.defaultMode.toString(8)
              }
              if (volume.configMap.optional !== undefined) {
                item.optional = volume.configMap.optional
              }
              if (volume.configMap.name) {
                item.resource = volume.configMap.name
              }
            }
            if (volume.secret) {
              item.type = "secret"
              if (volume.secret.defaultMode) {
                item.defaultMode = volume.secret.defaultMode.toString(8)
              }
              if (volume.secret.optional !== undefined) {
                item.optional = volume.secret.optional
              }
              if (volume.secret.secretName) {
                item.resource = volume.secret.secretName
              }
            }
            if (volume.persistentVolumeClaim) {
              item.type = "persistentVolumeClaim"
              if (volume.persistentVolumeClaim.name) {
                item.resource = volume.persistentVolumeClaim.name
              }
              if (volume.persistentVolumeClaim.readOnly !== undefined) {
                item.readOnly = volume.persistentVolumeClaim.readOnly
              }
            }
          }
          this.volumes.push(item)
        }
      }
    }
  },
}
</script>