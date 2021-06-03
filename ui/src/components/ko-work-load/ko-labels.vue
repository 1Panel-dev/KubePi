<template>
  <div style="margin-top: 20px">
      <el-button @click="handleAdd">Add</el-button>
      <el-table v-if="labels.length !== 0" :data="labels">
        <el-table-column min-width="80" label="Key">
          <template v-slot:default="{row}">
            <ko-form-item placeholder="e.g. foo" clearable itemType="input" v-model="row.key" />
          </template>
        </el-table-column>
        <el-table-column min-width="80" label="Value">
          <template v-slot:default="{row}">
            <ko-form-item placeholder="e.g. bar" clearable itemType="input" v-model="row.value" />
          </template>
        </el-table-column>
        <el-table-column width="120px">
          <template v-slot:default="{row}">
            <el-button type="text" style="font-size: 20px" @click="handleDelete(row)">REMOVE</el-button>
          </template>
        </el-table-column>
      </el-table>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoLabels",
  components: { KoFormItem },
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
      var item = {
        key: "",
        value: "",
      }
      this.labels.unshift(item)
    },
  },
}
</script>