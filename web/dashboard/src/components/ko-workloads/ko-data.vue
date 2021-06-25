<template>
  <div class="tab-content">
    <div class="tab-title">
      <span>Data</span>
    </div>
    <table style="width: 90%" class="tab-table">
      <tr>
        <th scope="col" width="48%" align="left">
          <label>key</label>
        </th>
        <th scope="col" width="48%" align="left">
          <label>value</label>
        </th>
        <th></th>
      </tr>
      <tr v-for="label in data" v-bind:key="label.index">
        <td >
          <ko-form-item :withoutLabel="true" placeholder="e.g. foo" clearable itemType="input" v-model="label.key"/>
        </td>
        <td >
          <ko-form-item :withoutLabel="true" placeholder="e.g. bar" clearable itemType="input" v-model="label.value"/>
        </td>
        <td>
          <el-button type="text" style="font-size: 10px" @click="handleDelete(label)">
            {{ $t("commons.button.delete") }}
          </el-button>
        </td>
      </tr>
      <tr>
        <td>
          <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoData",
  components: {KoFormItem},
  props: {
    dataObj: Object
  },
  data () {
    return {
      data: []
    }
  },
  methods: {
    handleDelete (row) {
      for (let i = 0; i < this.data.length; i++) {
        if (this.data[i] === row) {
          this.data.splice(i, 1)
        }
      }
    },
    handleAdd () {
      const item = {
        index: Math.random(),
        key: "",
        value: "",
      }
      this.data.push(item)
    },
    transformation (parentFrom) {
      if (this.data) {
        let obj = {}
        for (let i = 0; i < this.data.length; i++) {
          if (this.data[i].key !== "") {
            obj[this.data[i].key] = this.data[i].value
          }
        }
        parentFrom.data = obj
      }
    },
  },
  mounted () {
    if (this.dataObj && this.dataObj.data) {
      for (const key in this.dataObj.data) {
        if (Object.prototype.hasOwnProperty.call(this.dataObj.data, key)) {
          this.data.push({
            index: Math.random(),
            key: key,
            value: this.dataObj.data[key],
          })
        }
      }
    }
  }
}
</script>

<style scoped>
</style>