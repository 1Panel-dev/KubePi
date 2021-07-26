<template>
  <div style="margin-top: 20px">
    <ko-card title="Volume Claim Templates">
      <el-form label-position="top" :disabled="isReadOnly">
        <div v-for="(item, index) in volumeClaimTemplates" :key="index">
          <el-card style="margin-top: 10px">
            <div slot="header" class="clearfix">
              <el-button style="float: right; padding: 3px 0" type="text" @click="handleVolumeDelete(index)">删 除</el-button>
            </div>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Persistent Volume Name">
                  <ko-form-item itemType="input" v-model="item.name" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Type">
                  <ko-form-item itemType="radio" v-model="item.type" :radios="type_list" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Storage Class">
                  <ko-form-item itemType="select2" v-model="item.storageClass" :selections="sc_list" />
                </el-form-item>
              </el-col>
              <el-col :span="12" v-if="item.type === 'new'">
                <el-form-item label="Storage Class">
                  <ko-form-item itemType="number" deviderName="GiB" v-model.number="item.storage" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="Access Modes">
                  <ko-form-item itemType="checkbox" v-model="item.accessModes" :checks="access_mode_list" />
                </el-form-item>
              </el-col>
            </el-row>
          </el-card>
        </div>
        <el-row>
          <el-col :span="12">
            <el-button @click="handleVolumeAdd">Add Claim Templates</el-button>
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
  name: "KoVolumeClaim",
  components: { KoFormItem, KoCard },
  props: {
    volumeClaimParentObj: Object,
    currentNamespace: String,
    scList: Array,
    isReadOnly: Boolean,
  },
  watch: {
    currentNamespace: {
      handler(newName) {
        this.namespace = newName
      },
      immediate: true,
    },
    scList: {
      handler(newName) {
        this.sc_list = []
        for (const s of newName) {
          this.sc_list.push(s.provisioner)
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      namespace: "",
      volumeClaimTemplates: [],
      sc_list: [],
      type_list: [
        { label: "Use a Storage Class to provision a new Persistent Volume", value: "new" },
        { label: "Use an existing Persistent Volume", value: "existing" },
      ],
      access_mode_list: [
        { label: "Single-Node(Read/Wrete)", value: "ReadWriteOnce" },
        { label: "Many-Node(Read-Only)", value: "ReadOnlyMany" },
        { label: "Many-Node(Read/Wrete)", value: "ReadWriteMany" },
      ],
    }
  },
  methods: {
    handleVolumeAdd() {
      this.volumeClaimTemplates.push({
        name: "",
        type: "new",
        storageClass: "",
        accessModes: [],
        storage: 10,
      })
    },
    handleVolumeDelete(index) {
      this.volumeClaimTemplates.splice(index, 1)
    },

    transformation(parentFrom) {
      parentFrom.volumeClaimTemplates = []
      for (const volume of this.volumeClaimTemplates) {
        let item = {
          type: "persistentvolumeclaim",
          metadata: {
            namespace: this.currentNamespace,
            name: volume.name,
          },
          spec: {
            accessModes: volume.accessModes,
            storageClassName: volume.storageClass,
            resources: {
              requests: {},
            },
          },
        }
        if (volume.type === "new") {
          item.spec.resources.requests.storage = volume.storage.toString() + "Gi"
        }
        parentFrom.volumeClaimTemplates.push(item)
      }
    },
  },
  mounted() {
    if (this.volumeClaimParentObj) {
      if (this.volumeClaimParentObj.volumeClaimTemplates) {
        for (const volume of this.volumeClaimParentObj.volumeClaimTemplates) {
          this.volumeClaimTemplates.push({
            name: volume.metadata.name,
            type: volume.spec?.resources?.requests?.storage ? "new" : "existing",
            storageClass: volume.spec.storageClassName,
            accessModes: volume.spec.accessModes,
            storage: volume.spec?.resources?.requests?.storage ? Number(volume.spec.resources.requests.storage.replace("Gi", "")) : 0,
          })
        }
      }
    }
  },
}
</script>