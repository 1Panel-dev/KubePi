<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.workload.port')">
      <el-form :disabled="isReadOnly">
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="16%" align="left"><label>{{$t('business.workload.service_type')}}</label></th>
            <th scope="col" width="16%" align="left"><label>{{$t('business.workload.name')}}</label></th>
            <th scope="col" width="10%" align="left"><label>{{$t('business.workload.private_port')}}</label></th>
            <th scope="col" width="8%" align="left"><label>{{$t('business.workload.protocol')}}</label></th>
            <th scope="col" width="5%" align="left"><label>{{$t('business.workload.expose')}}</label></th>
            <th scope="col" width="10%" align="left"><label>{{$t('business.workload.public_port')}}</label></th>
            <th scope="col" width="10%" align="left"><label>{{$t('business.workload.listening_port')}}</label></th>
            <th scope="col" width="15%" align="left"><label>{{$t('business.workload.host_ip')}}</label></th>
            <th align="left"></th>
          </tr>
          <tr v-for="(row, index) in ports" v-bind:key="index">
            <td>
              <ko-form-item itemType="select" v-model="row._serviceType" :selections="service_type_list" />
            </td>
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
              <el-switch v-model="row.expose" />
            </td>
            <td>
              <ko-form-item :disabled="row._serviceType === 'NodePort' || row._serviceType === 'LoadBalancer' || !row.expose" placeholder="e.g. 80" itemType="number" v-model.number="row.hostPort" />
            </td>
            <td>
              <ko-form-item :disabled="row._serviceType === 'ClusterIP' || row._serviceType === '' || !row.expose" placeholder="e.g. 80" itemType="number" v-model.number="row._listeningPort" />
            </td>
            <td>
              <ko-form-item :disabled="row._serviceType === 'NodePort' || row._serviceType === 'LoadBalancer' || !row.expose" placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.hostIP" />
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
    isReadOnly: Boolean,
  },
  data() {
    return {
      ports: [],
      protocol_list: [
        { label: "TCP", value: "TCP" },
        { label: "UDP", value: "UDP" },
      ],
      service_type_list: [
        { label: this.$t("business.workload.not_create"), value: "" },
        { label: this.$t("business.workload.cluster_ip"), value: "ClusterIP" },
        { label: this.$t("business.workload.node_port"), value: "NodePort" },
        { label: this.$t("business.workload.load_balancer"), value: "LoadBalancer" },
      ],
    }
  },
  methods: {
    handleAdd() {
      var item = {
        name: "",
        expose: false,
        protocol: "TCP",
        containerPort: 80,
        hostPort: "",
        hostIP: "",
        _serviceType: "",
        _listeningPort: "",
      }
      this.ports.push(item)
    },
    handleDelete(index) {
      this.ports.splice(index, 1)
    },
    transformation(parentFrom) {
      if (this.ports.length !== 0) {
        parentFrom.ports = []
        for (const po of this.ports) {
          var itemPo = {}
          itemPo.name = po.name || undefined
          itemPo.expose = po.expose || undefined
          itemPo.protocol = po.protocol || undefined
          itemPo.containerPort = po.containerPort || undefined
          itemPo._serviceType = po._serviceType || undefined
          if (po.expose) {
            switch (po._serviceType) {
              case "":
              case "ClusterIP":
                itemPo.hostPort = po.hostPort || undefined
                itemPo.hostIP = po.hostIP || undefined
                break
              case "NodePort":
              case "LoadBalancer":
                itemPo._listeningPort = po._listeningPort || undefined
                break
            }
          }
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
          if (po.expose !== undefined) {
            itemPo.expose = po.expose
          }
          itemPo.protocol = po.protocol
          itemPo.containerPort = po.containerPort
          itemPo._serviceType = po._serviceType ? po._serviceType : ""
          if (po.expose) {
            switch (po._serviceType) {
              case "":
              case "ClusterIP":
                itemPo.hostPort = po.hostPort
                itemPo.hostIP = po.hostIP
                break
              case "NodePort":
              case "LoadBalancer":
                itemPo._listeningPort = po._listeningPort
                break
            }
          }
          this.ports.push(itemPo)
        }
      }
    }
  },
}
</script>