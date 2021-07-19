<template>
  <div style="margin-top: 20px">
    <ko-card title="Ports">
      <el-form :disabled="isReadOnly">
        <table style="width: 98%" class="tab-table">
          <tr>
            <th scope="col" width="16%" align="left"><label>service type</label></th>
            <th scope="col" width="16%" align="left"><label>name</label></th>
            <th scope="col" width="10%" align="left"><label>private port</label></th>
            <th scope="col" width="8%" align="left"><label>protocol</label></th>
            <th scope="col" width="5%" align="left"><label>expose</label></th>
            <th scope="col" width="10%" align="left"><label>public port</label></th>
            <th scope="col" width="10%" align="left"><label>listening port</label></th>
            <th scope="col" width="15%" align="left"><label>host ip</label></th>
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
        { label: "Not create a service", value: "" },
        { label: "Cluster IP", value: "ClusterIP" },
        { label: "Node Port", value: "NodePort" },
        { label: "Load Balancer", value: "LoadBalancer" },
      ],
    }
  },
  methods: {
    handleAdd() {
      var item = {
        name: "",
        expose: false,
        protocol: "TCP",
        containerPort: "",
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
          itemPo.name = po.name
          itemPo.expose = po.expose
          itemPo.protocol = po.protocol
          itemPo.containerPort = po.containerPort
          itemPo._serviceType = po._serviceType
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