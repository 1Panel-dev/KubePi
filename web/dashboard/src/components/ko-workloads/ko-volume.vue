<template>
  <div style="margin-top: 20px">
    <el-collapse v-model="activeName">
      <el-collapse-item style="margin-top: 10px" :title="item.name + '(' + item.type + ')'" :name="index" v-for="(item, index) in volumes" :key="index">
        <el-form label-position="top" :disabled="isReadOnly">
          <el-button style="float: right;margin-top: 5px;margin-right: 20px; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">{{$t("commons.button.delete")}}</el-button>
          <el-form-item :label="$t('business.workload.volume_name')" required>
            <ko-form-item itemType="input" v-model="item.name" @change="changeName" />
          </el-form-item>
          <el-form-item v-if="hasDefaultMode(item.type)" :label="$t('business.workload.default_mode')">
            <ko-form-item itemType="number" placeholder="644" v-model.number="item.defaultMode" />
          </el-form-item>
          <el-form-item :label="item.type" v-if="hasResource(item.type)">
            <ko-form-item v-if="item.type === 'ConfigMap'" itemType="select2" v-model="item.resource" :selections="config_map_name_list" @change="changeConfigMap(item)" />
            <ko-form-item v-if="item.type === 'Secret'" itemType="select2" v-model="item.resource" :selections="secret_name_list" @change="changeSecret(item)" />
            <ko-form-item v-if="item.type === 'PVC'" itemType="select2" v-model="item.resource" :selections="pvc_list" />
          </el-form-item>
          <table v-if="hasKvData(item.type)" style="width: 98%" class="tab-table">
            <tr>
              <th scope="col" width="43%" align="left">
                <label>{{$t('business.workload.key')}}</label>
              </th>
              <th scope="col" width="43%" align="left">
                <label>{{$t('business.workload.path')}}</label>
              </th>
              <th align="left"></th>
            </tr>
            <tr v-for="(row, index) in item.items" v-bind:key="index">
              <td>
                <ko-form-item itemType="select2" :selections="item.options" v-model="row.key" />
              </td>
              <td>
                <ko-form-item itemType="input" v-model="row.path" />
              </td>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleDelete(item, index)">
                  {{ $t("commons.button.delete") }}
                </el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button :disabled="item.resource === ''" @click="handleAdd(item)">{{ $t("commons.button.add") }}</el-button>
              </td>
            </tr>
          </table>
          <el-form-item label="Optional" v-if="hasOptional(item.type)">
            <ko-form-item itemType="radio" v-model="item.optional" :radios="optional_list" />
          </el-form-item>
          <div v-if="item.type === 'NFS'">
            <el-form-item label="Path">
              <ko-form-item itemType="input" v-model="item.path" />
            </el-form-item>
            <el-form-item label="Server">
              <ko-form-item itemType="input" v-model="item.server" />
            </el-form-item>
          </div>
          <div v-if="item.type === 'HostPath'">
            <el-form-item :label="$t('business.storage.path_or_node')" required>
              <ko-form-item itemType="input" v-model="item.path" />
            </el-form-item>
            <el-form-item :label="$t('business.workload.type')">
              <ko-form-item v-model="item.hostType" itemType="select2" :selections="host_path_list" />
            </el-form-item>
          </div>
        </el-form>
      </el-collapse-item>
    </el-collapse>
    <el-row style="margin-top: 10px">
      <el-col :span="12">
        <el-dropdown placement="bottom" trigger="click" @command="handleVolumeAdd">
          <el-button class="search-btn">
            {{$t('business.workload.add')}}
            <i class="el-icon-arrow-down el-icon--right"></i>
          </el-button>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item v-for="(item, index) in volume_type_list" :key="index" :command="item.value">{{item.label}}</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoStorage",
  components: { KoFormItem },
  props: {
    volumeParentObj: Object,
    configMapList: Array,
    secretList: Array,
    pvcList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    configMapList: {
      handler(newName) {
        this.config_map_name_list = []
        this.config_map_list = []
        for (const cm of newName) {
          this.config_map_name_list.push(cm.metadata.name)
          this.config_map_list.push(cm)
        }
      },
      immediate: true,
      deep: true,
    },
    secretList: {
      handler(newName) {
        this.secret_list = newName
        this.secret_name_list = []
        for (const s of newName) {
          this.secret_name_list.push(s.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    },
    pvcList: {
      handler(newName) {
        this.pvc_list = []
        for (const pvc of newName) {
          this.pvc_list.push(pvc.metadata.name)
        }
      },
      immediate: true,
      deep: true,
    }, 
  },
  data() {
    return {
      activeName: [],
      containerIndex: 0,
      volumes: [],
      pvc_list: [],
      config_map_list: [],
      config_map_name_list: [],
      secret_list: [],
      secret_name_list: [],
      volume_type_list: [
        { label: "NFS", value: "NFS" },
        { label: "PVC", value: "PVC" },
        { label: "EmptyDir", value: "EmptyDir" },
        { label: "HostPath", value: "HostPath" },
        { label: "ConfigMap", value: "ConfigMap" },
        { label: "Secret", value: "Secret" },
      ],
      optional_list: [
        { label: this.$t("business.workload.yes"), value: true },
        { label: this.$t("business.workload.no"), value: false },
      ],
      host_path_list: ["DirectoryOrCreate", "Directory", "FileOrCreate", "File", "Socket", "CharDevice", "BlockDevice"],
    }
  },
  methods: {
    handleAdd(row) {
      const item = {
        key: "",
        path: "",
      }
      row.items.push(item)
    },
    handleDelete(row, index) {
      row.items.splice(index, 1)
    },
    changeConfigMap(row) {
      row.options = []
      for (const cm of this.config_map_list) {
        if (row.resource === cm.metadata.name) {
          for (const item of Object.keys(cm.data)) {
            row.options.push(item)
          }
        }
      }
    },
    changeSecret(row) {
      row.options = []
      for (const se of this.secret_list) {
        if (row.resource === se.metadata.name) {
          for (const item of Object.keys(se.data)) {
            row.options.push(item)
          }
        }
      }
    },

    handleVolumeAdd(type) {
      this.volumes.push({
        type: type,
        name: "",
        resource: "",
        defaultMode: null,
        optional: null,
        path: null,
        server: null,
        hostType: null,
        options: [],
        items: [],
      })
      this.activeName = this.volumes.length - 1
      this.$emit("loadVolumes", "volume", this.volumes, [])
    },
    handleVolumeDelete(index) {
      this.volumes.splice(index, 1)
      this.activeName = this.volumes.length - 1
      this.$emit("loadVolumes", "volume", this.volumes, [])
    },
    changeName() {
      this.$emit("loadVolumes", "volume", this.volumes, [])
    },
    hasResource(type) {
      return type === "ConfigMap" || type === "PVC" || type === "Secret"
    },
    hasOptional(type) {
      return type === "ConfigMap" || type === "Secret"
    },
    hasKvData(type) {
      return type === "ConfigMap" || type === "Secret"
    },
    hasDefaultMode(type) {
      return type === "ConfigMap" || type === "Secret"
    },
    isSupport(volumes) {
      return volumes.configMap || volumes.secret || volumes.persistentVolumeClaim || volumes.emptyDir || volumes.nfs || volumes.hostPath
    },
    // 先直接去掉支持的 6 种类型，然后再往里push，保留不支持的内容
    transformation(parentFrom) {
      if (!parentFrom.volumes) {
        parentFrom.volumes = []
      }
      parentFrom.volumes = parentFrom.volumes.filter((item) => {
        return !this.isSupport(item)
      })
      for (const volume of this.volumes) {
        let item = {}
        let itemClild = {}
        item.name = volume.name || undefined
        itemClild.defaultMode = volume.defaultMode ? parseInt("0" + volume.defaultMode.toString(), 8) : undefined
        itemClild.optional = volume.optional !== null ? volume.optional : undefined
        itemClild.path = volume.path || undefined
        itemClild.server = volume.server || undefined
        itemClild.type = volume.hostType || undefined

        switch (volume.type) {
          case "ConfigMap":
            if (volume.items.length !== 0) {
              item.items = volume.items
            }
            item.configMap = itemClild
            item.configMap.name = volume.resource || undefined
            break
          case "Secret":
            if (volume.items.length !== 0) {
              item.items = volume.items
            }
            item.secret = itemClild
            item.secret.secretName = volume.resource || undefined
            break
          case "PVC":
            item.persistentVolumeClaim = itemClild
            item.persistentVolumeClaim.claimName = volume.resource || undefined
            break
          case "EmptyDir":
            item.emptyDir = {}
            break
          case "NFS":
            item.nfs = itemClild
            break
          case "HostPath":
            item.hostPath = itemClild
            break
        }
        parentFrom.volumes.push(item)
      }
    },
  },
  mounted() {
    this.volumes = []
    if (this.volumeParentObj.volumes) {
      for (const volume of this.volumeParentObj.volumes) {
        if (this.isSupport(volume)) {
          let item = {}
          if (volume.name) {
            item.name = volume.name
          }
          let itemClild = {}
          if (volume.configMap) {
            item.type = "ConfigMap"
            item.items = volume.items || []
            itemClild = volume.configMap
          }
          if (volume.secret) {
            item.type = "Secret"
            item.items = volume.items || []
            itemClild = volume.secret
          }
          if (volume.persistentVolumeClaim) {
            item.type = "PVC"
            itemClild = volume.persistentVolumeClaim
          }
          if (volume.emptyDir) {
            item.type = "EmptyDir"
            itemClild = {}
          }
          if (volume.nfs) {
            item.type = "NFS"
            itemClild = volume.nfs
          }
          if (volume.hostPath) {
            item.type = "HostPath"
            itemClild = volume.hostPath
          }

          if (itemClild.defaultMode) {
            item.defaultMode = itemClild.defaultMode.toString(8)
          }
          if (itemClild.optional !== undefined) {
            item.optional = itemClild.optional
          }
          if (itemClild.name) {
            item.resource = itemClild.name
          }
          if (itemClild.path) {
            item.path = itemClild.path
          }
          if (itemClild.server) {
            item.server = itemClild.server
          }
          if (itemClild.type) {
            item.hostType = itemClild.type
          }

          this.volumes.push(item)
        }
      }
    }
  },
}
</script>