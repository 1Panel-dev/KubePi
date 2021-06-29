<template>
  <div style="margin-top: 20px">
    <div v-for="(item, index) in volumes" :key="index">
      <el-card style="margin-top: 20px">
        <span style="font-size: 18px;">{{item._type}}</span>
        <el-button style="float: right; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">删 除</el-button>
        <div v-if="item._type === 'persistentVolumeClaim'">
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <ko-form-item labelName="Volume Name" itemType="input" v-model="item.name" />
            </el-col>
            <el-col :span="12">
              <ko-form-item labelName="Persistent Volume Claim" itemType="select" v-model="item.weight" :selections="pvc_list" />
            </el-col>
          </el-row>
        </div>
        <div v-if="item._type === 'configMap'">
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <ko-form-item labelName="Volume Name" itemType="input" v-model="item.name" />
            </el-col>
            <el-col :span="12">
              <ko-form-item labelName="Default Mode" itemType="number" v-model.number="item.defaultMode" />
            </el-col>
          </el-row>
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <ko-form-item labelName="ConfigMap" itemType="select" v-model="item.configMap.name" :selections="config_map_name_list" />
            </el-col>
            <el-col :span="12">
              <ko-form-item labelName="Optional" itemType="radio" v-model="item.configMap.optional" :radios="optional_list" />
            </el-col>
          </el-row>
        </div>
        <div v-if="item._type === 'secret'">
          <el-row :gutter="20" style="margin-top: 20px;">
            <el-col :span="12">
              <ko-form-item labelName="Volume Name" itemType="input" v-model="item.name" />
            </el-col>
            <el-col :span="12">
              <ko-form-item labelName="Default Mode" itemType="input" v-model.number="item.defaultMode" />
            </el-col>
          </el-row>
          <el-row :gutter="20" style="margin-top: 20px">
            <el-col :span="12">
              <ko-form-item labelName="Secret" itemType="select" v-model="item.secret.name" :selections="secret_list" />
            </el-col>
            <el-col :span="12">
              <ko-form-item labelName="Optional" itemType="radio" v-model="item.secret.optional" :radios="optional_list" />
            </el-col>
          </el-row>
        </div>
        <div style="margin-top: 30px"><span>Valume Mount</span></div>
        <table style="width: 98%;" class="tab-table">
          <tr>
            <th scope="col" width="43%" align="left"><label>mount point</label></th>
            <th scope="col" width="43%" align="left"><label>sub path in volume</label></th>
            <th scope="col" width="8%" align="left"><label>read only</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in volumeMounts" v-bind:key="index">
            <td>
              <ko-form-item :withoutLabel="true" itemType="input" v-model="row.mountPath" />
            </td>
            <td>
              <ko-form-item :withoutLabel="true" itemType="input" v-model="row.subPath" />
            </td>
            <td>
              <el-checkbox v-model="row.readOnly" />
            </td>
            <td>
              <el-button type="text" style="font-size: 10px" @click="handleMountDelete(index)">
                {{ $t("commons.button.delete") }}
              </el-button>
            </td>
          </tr>
          <tr>
            <td align="left">
              <el-button @click="handleMountAdd()">Add Mount</el-button>
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
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
// import KoCard from "@/components/ko-card/index"
import { listSecrets } from "@/api/secrets"
import { listConfigMaps } from "@/api/configmaps"

export default {
  name: "KoStorage",
  components: { KoFormItem },
  props: {
    storageParentObj: Object,
    containerIndex: Number,
  },
  data() {
    return {
      volumes: [],
      pvc_list: [],
      volumeMounts: [],
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
    loadSecrets() {
      listSecrets(this.$route.params.cluster).then((res) => {
        this.secret_list = []
        for (const secret of res.items) {
          this.secret_list.push({ label: secret.metadata.name, value: secret.metadata.name })
        }
      })
    },
    loadConfigMaps() {
      listConfigMaps(this.$route.params.cluster).then((res) => {
        this.config_map_name_list = []
        for (const cm of res.items) {
          this.config_map_name_list.push({ label: cm.metadata.name, value: cm.metadata.name })
        }
      })
    },
    handleVolumeAdd(type) {
      switch (type) {
        case "configMap":
          this.volumes.push({
            _type: type,
            name: "",
            configMap: {
              name: "",
              defaultMode: "",
              optional: true,
            },
          })
          break
        case "secret":
          this.volumes.push({
            _type: type,
            name: "",
            secret: {
              name: "",
              defaultMode: "",
              optional: true,
            },
          })
          break
        case "persistentVolumeClaim":
          this.volumes.push({
            _type: type,
            name: "",
            persistentVolumeClaim: {
              claimName: "",
              readOnly: false,
            },
          })
          break
      }
    },
    handleVolumeDelete(index) {
      this.volumes.splice(index, 1)
    },

    handleMountAdd() {
      var item = {
        mountPath: "",
        subPath: "",
        readOnly: false,
      }
      this.volumeMounts.push(item)
    },
    handleMountDelete(index) {
      this.volumeMounts.splice(index, 1)
    },

    transformation(parentFrom) {
      parentFrom.volume = this.volumes
      parentFrom.containers[this.containerIndex].volumeMounts = this.volumeMounts
    },
  },
  mounted() {
    this.loadSecrets()
    this.loadConfigMaps()
    this.volumes = []
    this.volumeMounts = []
    if (this.storageParentObj.containers.length - 1 >= this.containerIndex) {
      if (this.storageParentObj.containers[this.containerIndex].volumeMounts) {
        this.volumeMounts = this.storageParentObj.containers[this.containerIndex].volumeMounts
      }
    }
    if (this.storageParentObj.volumes) {
      this.volumes = this.storageParentObj.volumes
    }
  },
}
</script>