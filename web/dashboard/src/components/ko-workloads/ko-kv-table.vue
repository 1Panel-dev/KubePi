<template>
  <div>
    <table style="width: 100%" class="tab-table">
      <tr v-for="(label,index) in data" v-bind:key="label.index">
        <td width="48%">
          <ko-form-item placeholder="e.g. foo" clearable itemType="input" v-model="label.key" @change.native="changeValues" />
        </td>
        <td width="48%">
          <ko-form-item placeholder="e.g. bar" clearable itemType="input" v-model="label.value" @change.native="changeValues" />
        </td>
        <td width="4%">
          <el-button :disabled="kvType === 'label' && label.key === 'k8s.kubepi.cn/name'" type="text" style="font-size: 10px" @click="handleDelete(index)">
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
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoKvTable",
  components: { KoFormItem },
  props: {
    dataObj: Object,
    kvType: String,
    resourceName: String,
  },
  watch: {
    resourceName: {
      handler(newVal) {
        if (newVal) {
          for (const item of this.data) {
            if (item.key === "k8s.kubepi.cn/name") {
              item.value = newVal
              break
            }
          }
          this.changeValues()
        }
      },
      immediate: true,
    },
  },
  data() {
    return {
      data: [{ "k8s.kubepi.cn/name": "" }],
    }
  },
  methods: {
    handleDelete(index) {
      this.data.splice(index, 1)
    },
    handleAdd() {
      const item = {
        key: "",
        value: "",
      }
      this.data.push(item)
    },
    changeValues() {
      if (this.data) {
        let obj = {}
        for (let i = 0; i < this.data.length; i++) {
          if (this.data[i].key !== "") {
            obj[this.data[i].key] = this.data[i].value
          }
        }
        this.$emit("updateKvDatas", this.kvType, obj)
      }
    },
  },
  mounted() {
    this.data = []
    if (this.dataObj) {
      for (const key in this.dataObj) {
        if (Object.prototype.hasOwnProperty.call(this.dataObj, key)) {
          this.data.push({
            key: key,
            value: this.dataObj[key],
          })
        }
      }
    }
  },
}
</script>

<style scoped>
</style>
