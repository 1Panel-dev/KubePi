<template>
  <div style="margin-top: 20px">
    <ko-card title="Ports">
      <el-checkbox style="margin-left: 20px; margin-top: 10px" v-model="isExpose">isExpose</el-checkbox>
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="40%" align="left"><label>name</label></th>
          <th scope="col" width="20%" align="left"><label>Private Container Port</label></th>
          <th scope="col" v-if="isExpose" width="20%" align="left"><label>Public Host Port</label></th>
          <th scope="col" width="16%" align="left"><label>Protocol</label></th>
          <th align="left" width="1%"></th>
        </tr>
        <tr v-for="row in ports" v-bind:key="row.index">
          <td>
            <ko-form-item :withoutLabel="true" itemType="input" v-model="row.name" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="number" v-model.number="row.containerPort" />
          </td>
          <td v-if="isExpose">
            <ko-form-item :withoutLabel="true" itemType="number" v-model.number="row.hostPort" />
          </td>
          <td>
            <ko-form-item :withoutLabel="true" itemType="select" v-model="row.protocol" :selections="protocol_list" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(row.index)">
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
import KoFormItem from "@/components/ko-form-item/index"
import KoCard from "@/components/ko-card/index"

export default {
  name: "KoPorts",
  components: { KoFormItem, KoCard },
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
    handleDelete(index) {
      this.ports.splice(index, 1)
    },
    handleAdd() {
      var item = {
        name: "",
        containerPort: "",
      }
      this.ports.push(item)
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