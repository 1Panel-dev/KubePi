<template>
  <div style="margin-top: 20px">
      <el-button @click="handleAdd">Add Port</el-button>
      <el-table v-if="ports.length !== 0" :data="ports">
        <el-table-column min-width="40" label="Name">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.name" />
          </template>
        </el-table-column>
        <el-table-column min-width="40" label="Private Container Port">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.port" />
          </template>
        </el-table-column>
        <el-table-column min-width="20" label="Public Host Port">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="input" v-model="row.targetPort" />
          </template>
        </el-table-column>
        <el-table-column min-width="20" label="Protocol">
          <template v-slot:default="{row}">
            <ko-form-item :withoutLabel="true" clearable itemType="select" v-model="row.protocol" :selections="protocol_list" />
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
  name: "KoPorts",
  components: { KoFormItem },
  data() {
    return {
      ports: [],
      protocol:"",
      protocol_list: [
        { label: "TCP", value: "TCP" },
        { label: "UDP", value: "UDP" },
        { label: "SCTP", value: "SCTP" },
      ],
    }
  },
  methods: {
    handleDelete(row) {
      for (let i = 0; i < this.ports.length; i++) {
        if (this.ports[i] === row) {
          this.ports.splice(i, 1)
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
      this.ports.unshift(item)
    },
  },
}
</script>