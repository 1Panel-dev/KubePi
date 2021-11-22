<template>
  <div>
    <el-form :disabled="isReadOnly">
      <table style="width: 98%" class="tab-table">
        <tr>
          <th scope="col" width="45%" align="left"><label>{{$t('business.workload.name')}}</label></th>
          <th scope="col" width="15%" align="left"><label>{{$t('business.workload.container_port')}}</label></th>
          <th scope="col" width="15%" align="left"><label>{{$t('business.workload.protocol')}}</label></th>
          <th scope="col" width="15%" align="left"><label>{{$t('business.workload.host_port')}}</label></th>
          <th align="left"></th>
        </tr>
        <tr v-for="(row, index) in ports" v-bind:key="index">
          <td>
            <ko-form-item itemType="input" v-model="row.name" />
          </td>
          <td>
            <ko-form-item placeholder="e.g. 8080" itemType="number" v-model.number="row.containerPort" />
          </td>
          <td>
            <ko-form-item itemType="select" v-model="row.protocol" :selections="protocol_list" />
          </td>
          <td>
            <ko-form-item placeholder="e.g. 80" itemType="number" v-model.number="row.hostPort" />
          </td>
          <td>
            <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">
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
    </el-form>
  </div>
</template>
          
<script>
import KoFormItem from "@/components/ko-form-item/index"

export default {
  name: "KoPorts",
  components: { KoFormItem },
  props: {
    portParentObj: Object,
    isReadOnly: Boolean,
  },
  data() {
    return {
      ports: [],
      protocol_list: [
        { label: "TCP", value: "TCP" },
        { label: "UDP", value: "UDP" },
      ],
    }
  },
  methods: {
    handleAdd() {
      var item = {
        name: "",
        protocol: "TCP",
        containerPort: null,
        hostPort: "",
      }
      this.ports.push(item)
    },
    handleDelete(index) {
      this.ports.splice(index, 1)
    },
    transformation(parentFrom) {
      parentFrom.ports = []
      if (this.ports.length !== 0) {
        for (const po of this.ports) {
          var itemPo = {}
          itemPo.name = po.name || undefined
          itemPo.protocol = po.protocol || undefined
          itemPo.containerPort = po.containerPort || undefined
          itemPo.hostPort = po.hostPort || undefined
          parentFrom.ports.push(itemPo)
        }
      }
    },
  },
  mounted() {
    this.ports = []
    if (this.portParentObj) {
      if (this.portParentObj.ports) {
        for (const po of this.portParentObj.ports) {
          var itemPo = {}
          itemPo.name = po.name
          itemPo.protocol = po.protocol
          itemPo.containerPort = po.containerPort
          itemPo.hostPort = po.hostPort
          this.ports.push(itemPo)
        }
      }
    }
  },
}
</script>