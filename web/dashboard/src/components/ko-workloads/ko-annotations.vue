<template>

  <div class="tab-content">
    <el-row :gutter="20">
      <el-col :span="11">
        <span>key</span>
      </el-col>
      <el-col :span="11">
        <span>value</span>
      </el-col>
      <div v-for="label in annotations" v-bind:key="label.index">
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
  name: "KoAnnotations",
  components: { KoFormItem },
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
    if (this.annotationsParentObj.annotations) {
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