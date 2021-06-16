<template>

  <div class="tab-content">
    <el-row :gutter="20">
      <el-col :span="11">
        <span>key</span>
      </el-col>
      <el-col :span="11">
        <span>value</span>
      </el-col>
      <div v-for="label in labels" v-bind:key="label.index">
        <div class="grid-content tab-content">
          <br>
          <br>
          <el-col :span="11">
            <ko-form-item :withoutLabel="true" placeholder="e.g. foo" clearable itemType="input" v-model="label.key" />
          </el-col>
          <el-col :span="11">
            <ko-form-item :withoutLabel="true" placeholder="e.g. bar" clearable itemType="input" v-model="label.value" />
          </el-col>
          <el-col :span="2">
            <el-button type="text" style="font-size: 10px" @click="handleDelete(label)">
              {{ $t("commons.button.delete") }}
            </el-button>
          </el-col>
        </div>
      </div>
    </el-row>
    <div style="margin-top: 10px;">
      <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
    </div>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item"

export default {
  name: "KoLabels",
  components: { KoFormItem },
  props: {
    labelParentObj: Object,
  },
  data() {
    return {
      labels: [],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.labels.length; i++) {
        if (this.labels[i] === row) {
          this.labels.splice(i, 1)
        }
      }
    },
    handleAdd() {
      const item = {
        index: Math.random(),
        key: "",
        value: "",
      }
      this.labels.push(item)
    },
    transformation(parentFrom) {

      if (this.labels&&parentFrom) {
        if (parentFrom.labels === undefined) {
          parentFrom.labels = {}
        }
        for (let i = 0; i < this.labels.length; i++) {
          if (this.labels[i].key !== "") {
            parentFrom.labels[this.labels[i].key] = this.labels[i].value
          }
        }
      }
    },
  },
  mounted() {
    if (this.labelParentObj) {
      for(const key in this.labelParentObj.labels) {
        if (Object.prototype.hasOwnProperty.call(this.labelParentObj.labels, key)) {
          this.labels.push({
            index: Math.random(),
            key: key,
            value: this.labelParentObj.labels[key],
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