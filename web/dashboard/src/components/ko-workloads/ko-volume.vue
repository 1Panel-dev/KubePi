<template>
  <div style="margin-top: 20px">
    <el-form label-position="top" :disabled="isReadOnly">
      <div v-for="(item, index) in volumes" :key="index">
        <div style="margin-top: 20px">
          <span>{{item.type}}</span>
          <el-button style="float: right; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">{{$t("commons.button.delete")}}</el-button>
          <el-divider />
          <el-form-item :label="$t('business.workload.volume_name')" required>
            <ko-form-item itemType="input" v-model="item.name" @change="changeName" />
          </el-form-item>
          <el-form-item v-if="hasDefaultMode(item.type)" :label="$t('business.workload.default_mode')">
            <ko-form-item itemType="number" placeholder="644" v-model.number="item.defaultMode" />
          </el-form-item>
          <el-form-item :label="item.type" v-if="hasResource(item.type)">
            <ko-form-item v-if="item.type === 'ConfigMap'" itemType="select2" v-model="item.resource" :selections="config_map_name_list" />
            <ko-form-item v-if="item.type === 'Secret'" itemType="select2" v-model="item.resource" :selections="secret_list" />
            <ko-form-item v-if="item.type === 'PVC'" itemType="select2" v-model="item.resource" :selections="pvc_list" />
          </el-form-item>
          <el-form-item :label="$t('business.workload.optional')" v-if="hasOptional(item.type)">
            <ko-form-item itemType="radio" v-model="item.optional" :radios="optional_list" />
          </el-form-item>
          <el-form-item :label="$t('business.workload.read_only')" v-if="hasReadOnly(item.type)">
            <ko-form-item itemType="radio" v-model="item.readOnly" :radios="optional_list" />
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
        </div>
      </div>
      <el-row>
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
    </el-form>
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
    isReadOnly: Boolean,
  },
  watch: {
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
    handleVolumeAdd(type) {
      this.volumes.push({
        type: type,
        name: "",
        resource: "",
        defaultMode: null,
        optional: null,
        readOnly: null,
        path: null,
        server: null,
        hostType: null,
      })
      this.$emit("loadVolumes", "volume", this.volumes, [])
    },
    handleVolumeDelete(index) {
      this.volumes.splice(index, 1)
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
    hasReadOnly(type) {
      return type === "NFS" || type === "PVC"
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
        console.log(volume)
        item.name = volume.name || undefined
        itemClild.defaultMode = volume.defaultMode ? parseInt("0" + volume.defaultMode.toString(), 8) : undefined
        itemClild.optional = volume.optional !== null ? volume.optional : undefined
        itemClild.readOnly = volume.readOnly !== null ? volume.readOnly : undefined
        itemClild.path = volume.path || undefined
        itemClild.server = volume.server || undefined
        itemClild.type = volume.hostType || undefined

        switch (volume.type) {
          case "ConfigMap":
            item.configMap = itemClild
            item.configMap.name = volume.resource || undefined
            break
          case "Secret":
            item.secret = itemClild
            item.secret.secretName = volume.resource || undefined
            break
          case "PVC":
            item.persistentVolumeClaim = itemClild
            item.persistentVolumeClaim.name = volume.resource || undefined
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
            itemClild = volume.configMap
          }
          if (volume.secret) {
            item.type = "Secret"
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
          if (itemClild.readOnly !== undefined) {
            item.readOnly = itemClild.readOnly
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