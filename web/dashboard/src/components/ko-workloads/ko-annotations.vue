<template>
  <div style="margin-top: 20px">
    <ko-card title="Annotations">
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
        <tr v-for="label in annotations" v-bind:key="label.index">
          <td>
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo" itemType="input" v-model="label.key" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" placeholder="e.g. bar" itemType="input" v-model="label.value" />
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
    annotationsParentObj: Object,
  },
  data() {
    return {
      annotations: [],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.annotations.length; i++) {
        if (this.annotations[i] === row) {
          this.annotations.splice(i, 1)
        }
      }
    },
    handleAdd() {
      const item = {
        index: Math.random(),
        key: "",
        value: "",
      }
      this.annotations.push(item)
    },
    transformation(parentFrom) {
      if (this.annotations) {
        let obj = {}
        for (let i = 0; i < this.annotations.length; i++) {
          if (this.annotations[i].key !== "") {
            obj[this.annotations[i].key] = this.annotations[i].value
          }
        }
        parentFrom.annotations = obj
      }
    },
  },
  mounted() {
    if (this.annotationsParentObj && this.annotationsParentObj.annotations) {
      for (const key in this.annotationsParentObj.annotations) {
        if (Object.prototype.hasOwnProperty.call(this.annotationsParentObj.annotations, key)) {
          this.annotations.push({
            index: Math.random(),
            key: key,
            value: this.annotationsParentObj.annotations[key],
          })
        }
      }
    }
  },
}
</script>

<style scoped>
.tab-content {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>