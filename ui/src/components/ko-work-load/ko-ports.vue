<template>
  <div>
      <el-button @click="handleAdd">Add Port</el-button>
      <el-table style="color: black; background-color: black" :show-header="false" :data="data">
        <el-table-column min-width="80">
          <template v-slot:default="{row}">
            <ko-form-item labelName="Name" clearable itemType="input" v-model="row.name" />
          </template>
        </el-table-column>
        <el-table-column min-width="80">
          <template v-slot:default="{row}">
            <ko-form-item labelName="Private Container Port" clearable itemType="input" v-model="row.port" />
          </template>
        </el-table-column>
        <el-table-column min-width="30">
          <template v-slot:default="{row}">
            <ko-form-item labelName="Public Host Port" clearable itemType="input" v-model="row.targetPort" />
          </template>
        </el-table-column>
        <el-table-column min-width="30">
          <template v-slot:default="{row}">
            <ko-form-item labelName="Protocol" clearable itemType="select" v-model="row.protocol" :selections="protocol_list" />
          </template>
        </el-table-column>
        <el-table-column min-width="20">
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
  name: "KoPorts",
  components: { KoFormItem },
  data() {
    return {
      data: [],
      protocol_list: [
        { label: "TCP", value: "TCP" },
        { label: "UDP", value: "UDP" },
        { label: "SCTP", value: "SCTP" },
      ],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.data.length; i++) {
        if (this.data[i] === row) {
          this.data.splice(i, 1)
        }
      }
    },
    handleAdd() {
      var item = {
        name: "",
        port: "",
        targetPort: "",
        protocol: "",
      }
      this.data.unshift(item)
    },
  },
}
</script>