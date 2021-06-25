<template>
  <div style="margin-top: 20px">
    <el-button @click="handleAdd">Add Port</el-button>
    <el-checkbox style="margin-left: 20px; margin-top: 10px" v-model="isExpose">isExpose</el-checkbox>
    <el-table v-if="ports.length !== 0" :data="ports">
      <el-table-column min-width="40" label="Name">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="input" v-model="row.name" />
        </template>
      </el-table-column>
      <el-table-column min-width="40" label="Private Container Port">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="number" v-model.number="row.containerPort" />
        </template>
      </el-table-column>
      <el-table-column v-if="isExpose" min-width="20" label="Public Host Port">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="number" v-model.number="row.hostPort" />
        </template>
      </el-table-column>
      <el-table-column min-width="20" label="Protocol">
        <template v-slot:default="{row}">
          <ko-form-item :withoutLabel="true" itemType="select" v-model="row.protocol" :selections="protocol_list" />
        </template>
      </el-table-column>
      <el-table-column  width="60px">
        <template v-slot:default="{row}">
          <el-button type="text" style="font-size: 10px" @click="handleDelete(row)">{{ $t("commons.button.delete") }}</el-button>
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
  props: {
    portParentObj: Object,
  },
  data() {
    return {
      ports: [],
      protocol: "",
      isExpose: false,
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
        containerPort: "",
      }
      this.ports.unshift(item)
    },
    transformation(parentFrom) {
      if (this.ports.length !== 0) {
        parentFrom.ports = this.ports
        for (const po of parentFrom.ports) {
          po.expose = this.isExpose
        }
      }
    },
  },
  mounted() {
    if (this.portParentObj) {
      if (this.portParentObj.ports) {
        this.ports = this.portParentObj.ports
      }
    }
  },
}
</script>