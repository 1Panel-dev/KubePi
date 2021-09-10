<template>
  <div>
    <table style="width: 100%;" class="tab-table">
      <tr>
        <th scope="col" width="15%" align="left"><label>{{$t('business.workload.volume')}}</label></th>
        <th scope="col" width="15%" align="left"><label>{{$t('business.workload.type')}}</label></th>
        <th scope="col" width="28%" align="left"><label>{{$t('business.workload.mount_point')}}</label></th>
        <th scope="col" width="26%" align="left"><label>{{$t('business.workload.sub_path_in_volume')}}</label></th>
        <th scope="col" width="4%" align="left"><label>{{$t('business.workload.read_only')}}</label></th>
        <th align="left"></th>
      </tr>
      <tr v-for="(row, index) in volumeMounts" v-bind:key="index">
        <td>
          <ko-form-item itemType="select2" v-model="row.name" @change="changeVolume(row)" :selections="volume_name_list" />
        </td>
        <td>
          <ko-form-item itemType="input" disabled v-model="row.type" />
        </td>
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
          <el-button type="text" style="font-size: 10px" @click="handleMountDelete(index)">
            {{ $t("commons.button.delete") }}
          </el-button>
        </td>
      </tr>
      <tr>
        <td align="left">
          <el-button @click="handleMountAdd()">{{$t('business.workload.add')}}</el-button>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoStorage",
  components: { KoFormItem },
  props: {
    mountParentObj: Object,
    volumeList: Array,
  },
  watch: {
    currentContainer: {
      handler(newName) {
        this.currentContainer = newName
      },
      immediate: true,
      deep: true,
    },
    volumeList: {
      handler(newName) {
        this.volume_name_list = []
        this.volume_list = []
        if (newName) {
          this.volume_list = newName
          for (const v of newName) {
            this.volume_name_list.push(v.name)
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      volume_list: [],
      volume_name_list: [],
      volumeMounts: [],
      optional_list: [
        { label: this.$t("business.workload.yes"), value: true },
        { label: this.$t("business.workload.no"), value: false },
      ],
    }
  },
  methods: {
    handleMountAdd() {
      this.volumeMounts.push({
        name: "",
        mountPath: "",
        subPath: "",
        readOnly: false,
      })
    },
    changeVolume(row) {
      for (const v of this.volume_list) {
        if (v.name === row.name) {
          row.type = v.type
          return
        }
      }
    },
    loadVolumeMount() {
      for (const row of this.volumeMounts) {
        this.changeVolume(row)
      }
    },
    handleMountDelete(index) {
      this.volumeMounts.splice(index, 1)
    },
    transformation(parentFrom) {
      parentFrom.volumeMounts = []
      for (const mount of this.volumeMounts) {
        if (mount.name !== "") {
          parentFrom.volumeMounts.push({
            name: mount.name,
            mountPath: mount.mountPath || undefined,
            subPath: mount.subPath || undefined,
            readOnly: mount.readOnly || undefined,
          })
        }
      }
    },
  },
  mounted() {
    if (this.mountParentObj.volumeMounts) {
      for (const mount of this.mountParentObj.volumeMounts) {
        this.changeVolume(mount)
      }
      this.volumeMounts = this.mountParentObj.volumeMounts
    }
  },
}
</script>