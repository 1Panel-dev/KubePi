<template>
  <div style="margin-top: 20px">
    <ko-card :title="labelTitle">
      <el-form :disabled="isReadOnly">
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
          <tr v-for="(row, index) in labels" v-bind:key="index">
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
      </el-form>
    </ko-card>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoLabels",
  components: { KoFormItem, KoCard },
  props: {
    labelObj: Object,
    labelTitle: String,
    isReadOnly: Boolean,
  },
  data() {
    return {
      labels: [],
      test: "",
    }
  },
  methods: {
    handleAdd() {
      const item = {
        index: Math.random(),
        key: "",
        value: "",
      }
      this.labels.push(item)
    },
    handleDelete(index) {
      this.labels.splice(index, 1)
      this.transformation()
    },
    transformation() {
      if (this.labels) {
        let obj = {}
        for (let i = 0; i < this.labels.length; i++) {
          if (this.labels[i].key !== "") {
            obj[this.labels[i].key] = this.labels[i].value
          }
        }
        this.$emit("update:labelObj", obj)
      }
    },
    initData(obj) {
      if (obj) {
        this.labels = []
        for (const key in obj) {
          if (Object.prototype.hasOwnProperty.call(obj, key)) {
            this.labels.push({
              index: Math.random(),
              key: key,
              value: obj[key],
            })
          }
        }
      }
    },
  },
  created() {
    this.initData(this.labelObj)
  },
}
</script>
