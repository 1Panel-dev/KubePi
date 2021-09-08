<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.configuration.data')">
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
        <tr v-for="label in data" v-bind:key="label.index">
          <td>
            <ko-form-item  placeholder="e.g. foo" clearable itemType="input" v-model="label.key"
                          @change.native="transform"/>
          </td>
          <td>
            <ko-form-item  placeholder="e.g. bar" clearable itemType="textarea"
                          v-model="label.value" @change.native="transform"/>
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
            <el-upload :before-upload="readFile" action="" style="display: inline-block;margin-left: 5px">
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
  components: { KoCard, KoFormItem },
  props: {
    dataObj: Object,
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
      this.transform()
    },
    handleAdd () {
      const item = {
        index: Math.random(),
        key: "",
        value: "",
      }
      this.data.push(item)
      this.transform()
    },
    readFile (file) {
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
        this.transform()
      }
      return false
    },
    transform () {
      if (this.data && this.data.length > 0) {
        let obj = {}
        for (let i = 0; i < this.data.length; i++) {
          if (this.data[i].key !== "") {
            const { Base64 } = require("js-base64")
            obj[this.data[i].key] = Base64.encode(this.data[i].value)
          }
        }
        this.$emit("update:dataObj", obj)
      }
    }
  },
  mounted () {
    if (this.dataObj) {
      for (const key in this.dataObj) {
        if (Object.prototype.hasOwnProperty.call(this.dataObj, key)) {
          const { Base64 } = require("js-base64")
          const value = Base64.decode(this.dataObj[key])
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
