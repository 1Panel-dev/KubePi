<template>
  <div style="margin-top: 20px">
    <el-collapse v-model="activeName">
      <el-collapse-item style="margin-top: 10px" :title="item.name" :name="index" v-for="(item, index) in volumeClaimTemplates" :key="index">
        <el-form label-position="top" :disabled="isReadOnly">
          <el-button style="float: right;margin-top: 5px;margin-right: 20px; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">{{$t("commons.button.delete")}}</el-button>
          <el-form-item :label="$t('business.workload.pv_name')">
            <ko-form-item itemType="input" @change="changeName" v-model="item.name" />
          </el-form-item>
          <el-form-item :label="$t('business.workload.type')">
            <ko-form-item itemType="radio" radioLayout="vertical" v-model="item.type" :radios="type_list" />
          </el-form-item>
          <div v-if="item.type === 'new'">
            <el-form-item :label="$t('business.workload.storage_class')">
              <ko-form-item itemType="select2" v-model="item.storageClass" :selections="sc_list" />
            </el-form-item>
            <el-form-item :label="$t('business.workload.size')" v-if="item.type === 'new'">
              <ko-form-item itemType="number" deviderName="Gi" v-model.number="item.storage" />
            </el-form-item>
          </div>
          <el-form-item v-if="item.type === 'existing'" :label="$t('business.workload.pvc')">
            <ko-form-item itemType="select2" v-model="item.volumeName" :selections="pvc_list" />
          </el-form-item>
          <el-form-item :label="$t('business.workload.access_modes')">
            <ko-form-item itemType="checkbox" v-model="item.accessModes" :checks="access_mode_list" />
          </el-form-item>
        </el-form>
      </el-collapse-item>
    </el-collapse>
    <el-row style="margin-top : 10px">
      <el-col :span="12">
        <el-button @click="handleVolumeAdd">{{$t('business.workload.add')}}</el-button>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoVolumeClaim",
  components: { KoFormItem },
  props: {
    volumeClaimParentObj: Object,
    scList: Array,
    pvcList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    scList: {
      handler(newName) {
        this.sc_list = []
        for (const s of newName) {
          this.sc_list.push(s.metadata.name)
        }
      },
      immediate: true,
    },
    pvcList: {
      handler(newName) {
        this.pvc_list = []
        for (const s of newName) {
          this.pvc_list.push(s.metadata.name)
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      activeName: [],
      volumeClaimTemplates: [],
      sc_list: [],
      pvc_list: [],
      type_list: [
        { label: this.$t("business.workload.new_pv"), value: "new" },
        { label: this.$t("business.workload.existing_pv"), value: "existing" },
      ],
      access_mode_list: [
        { label: this.$t("business.workload.single_read_write"), value: "ReadWriteOnce" },
        { label: this.$t("business.workload.many_read_only"), value: "ReadOnlyMany" },
        { label: this.$t("business.workload.many_read_write"), value: "ReadWriteMany" },
      ],
    }
  },
  methods: {
    handleVolumeAdd() {
      this.volumeClaimTemplates.push({
        name: "",
        type: "new",
        storageClass: null,
        accessModes: [],
        storage: 10,
      })
      this.activeName = this.volumeClaimTemplates.length - 1
      this.$emit("loadVolumes", "VolumeClaimTemplates", [], this.volumeClaimTemplates)
    },
    handleVolumeDelete(index) {
      this.volumeClaimTemplates.splice(index, 1)
      this.activeName = this.volumeClaimTemplates.length - 1
      this.$emit("loadVolumes", "VolumeClaimTemplates", [], this.volumeClaimTemplates)
    },
    changeName() {
      this.$emit("loadVolumes", "VolumeClaimTemplates", [], this.volumeClaimTemplates)
    },

    transformation(parentFrom) {
      let volumeClaimTemplates = []
      for (const volume of this.volumeClaimTemplates) {
        let item = {
          metadata: {
            name: volume.name,
          },
          spec: {
            accessModes: volume.accessModes,
          },
        }
        if (volume.type === "new") {
          item.spec.resources = { requests: { storage: volume.storage.toString() + "Gi" } }
          item.spec.storageClassName = volume.storageClass || undefined
        } else {
          item.spec.volumeName = volume.volumeName
        }
        volumeClaimTemplates.push(item)
      }
      parentFrom.volumeClaimTemplates = volumeClaimTemplates.length !== 0 ? volumeClaimTemplates : undefined
    },
  },
  mounted() {
    if (this.volumeClaimParentObj) {
      if (this.volumeClaimParentObj.volumeClaimTemplates) {
        for (const volume of this.volumeClaimParentObj.volumeClaimTemplates) {
          this.volumeClaimTemplates.push({
            name: volume.metadata.name,
            type: volume.spec?.resources?.requests?.storage ? "new" : "existing",
            storageClass: volume.spec.storageClassName || "",
            volumeName: volume.spec.volumeName || "",
            accessModes: volume.spec.accessModes,
            storage: volume.spec?.resources?.requests?.storage ? Number(volume.spec.resources.requests.storage.replace("Gi", "")) : 0,
          })
        }
      }
    }
  },
}
</script>