<template>
  <div style="margin-top: 20px">
    <ko-card title="Data">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="48%" align="left">
            <label>key</label>
          </th>
          <th scope="col" width="48%" align="left">
            <label>value</label>
          </th>
          <th align="left"></th>
        </tr>
        <tr v-for="label in data" v-bind:key="label.index">
          <td>
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo" clearable itemType="input" v-model="label.key"/>
          </td>
          <td>
            <ko-form-item :withoutLabel="true" placeholder="e.g. bar" clearable itemType="textarea" v-model="label.value"/>
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(label)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </td>
        </tr>
        <tr>
          <td align="left">
            <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
            <el-upload  :before-upload="readFile" action=""  style="display: inline-block;margin-left: 5px">
              <el-button>{{ $t("commons.button.upload") }}</el-button>
            </el-upload>
          </td>
        </tr>
      </table>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoSecretData",
  components: { KoCard ,KoFormItem},
  props: {},
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
    readFile(file) {
      const reader = new FileReader()
      reader.readAsText(file)
      reader.onerror = e => {
        console.log("error" + e)
      }
      reader.onload = () => {
        const item = {
          index: Math.random(),
          key: file.name,
          value: reader.result,
        }
        this.data.push(item)
      }
      return false
    },
    transformation (parentFrom) {
      if (this.data) {
        let obj = {}
        for (let i = 0; i < this.data.length; i++) {
          if (this.data[i].key !== "") {
            const { Base64 } = require("js-base64")
            obj[this.data[i].key] = Base64.encode( this.data[i].value)
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
          const { Base64 } = require("js-base64")
          const value = Base64.decode( this.data[i].value)
          this.data.push({
            index: Math.random(),
            key: key,
            value: value,
          })
        }
      }
    }
  }
}
</script>

<style scoped>

</style>