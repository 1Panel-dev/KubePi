<template>
  <div style="margin-top: 20px">
    <ko-card :title="annotationsTitle">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="48%" align="left">
            <label>{{$t('business.workload.key')}}</label>
          </th>
          <th scope="col" width="48%" align="left">
            <label>{{$t('business.workload.value')}}</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in annotations" v-bind:key="index">
          <td>
            <ko-form-item placeholder="e.g. foo" itemType="input" v-model="row.key" @change.native="transformation" />
          </td>
          <td>
            <ko-form-item placeholder="e.g. bar" itemType="textarea" v-model="row.value" @change.native="transformation" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
          </td>
        </tr>
      </table>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoAnnotations",
  components: { KoFormItem, KoCard },
  props: {
    annotationsObj: Object,
    annotationsTitle: String,
  },
  data() {
    return {
      annotations: [],
    }
  },
  methods: {
    handleAdd() {
      const item = {
        key: "",
        value: "",
      }
      this.annotations.push(item)
    },
    handleDelete(index) {
      this.annotations.splice(index, 1)
      this.transformation()
    },
    transformation() {
      if (this.annotations) {
        let obj = {}
        for (let i = 0; i < this.annotations.length; i++) {
          if (this.annotations[i].key !== "") {
            obj[this.annotations[i].key] = this.annotations[i].value
          }
        }
        this.$emit("update:annotationsObj", obj)
      }
    },
    initData(obj) {
      if (obj) {
        for (const key in obj) {
          if (Object.prototype.hasOwnProperty.call(obj, key)) {
            this.annotations.push({
              index: Math.random(),
              key: key,
              value: this.annotationsObj[key],
            })
          }
        }
      }
    },
  },
  mounted() {
    this.initData(this.annotationsObj)
  },
}
</script>

<style scoped>
</style>
