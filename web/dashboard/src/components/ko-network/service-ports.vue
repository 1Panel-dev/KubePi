<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.network.service_port')">
      <table style="width: 100%;padding: 0" class="tab-table">
        <tr>
          <th scope="col" width="20%" align="left">
            <label>{{$t('business.network.port_name')}}</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>{{$t('business.network.listening_port')}}</label>
          </th>
          <th scope="col" width="10%" align="left">
            <label>{{$t('business.network.protocol')}}</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>{{$t('business.network.target_port')}}</label>
          </th>
          <th scope="col"  width="20%" align="left" v-if="type==='NodePort'">
            <label>{{$t('business.network.node_port')}}</label>
          </th>
          <th>
          </th>
        </tr>
        <tr v-for="(row,index) in servicePorts" v-bind:key="index">
          <td>
            <el-input v-model="row.name" ></el-input>
          </td>
          <td>
            <el-input v-model.number="row.port" placeholder="8080"></el-input>
          </td>
          <td>
            <el-select v-model="row.protocol" style="width: 100%" >
              <el-option label="TCP" value="TCP"></el-option>
              <el-option label="UDP" value="UDP"></el-option>
            </el-select>
          </td>
          <td>
            <el-input type="text" v-model="row.targetPort" @change="transformation(index)"></el-input>
          </td>
          <th  v-if="type==='NodePort'">
            <el-input type="number" v-model="row.nodePort" @change="transformation(index)"></el-input>
          </th>
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
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"

export default {
  name: "KoServicePorts",
  components: { KoCard },
  props: {
    type: String,
    ports: Array
  },
  data () {
    return {
      servicePorts: []
    }
  },
  methods: {
    handleAdd () {
      const item = {
        name: "",
        port: "",
        protocol: "TCP",
        targetPort: ""
      }
      this.servicePorts.push(item)
    },
    handleDelete (index) {
      this.servicePorts.splice(index, 1)
    },
    transformation (index) {
      if (index !== undefined && this.servicePorts[index].targetPort !== "" &&!Number.isNaN(Number(this.servicePorts[index].targetPort))) {
        this.servicePorts[index].targetPort = parseInt(this.servicePorts[index].targetPort)
      }
      if (index !== undefined &&this.servicePorts[index].nodePort !== "" && !Number.isNaN(Number(this.servicePorts[index].nodePort))) {
        this.servicePorts[index].nodePort = parseInt(this.servicePorts[index].nodePort)
      } else {
        delete this.servicePorts[index].nodePort
      }
      this.$emit("update:ports", this.servicePorts)
    }
  },
  watch: {
    type: function (newValue) {
      if (newValue !== 'NodePort') {
        for (let i =0;i<this.servicePorts.length;i++) {
          delete this.servicePorts[i].nodePort
        }
      }
    }
  },
  created () {
    if (!this.ports) {
      this.servicePorts.push({
        name: "",
        port: "",
        protocol: "TCP",
        targetPort: ""
      })
      this.$emit("update:ports", this.servicePorts)
    } else {
      this.servicePorts = this.ports
    }
  }
}
</script>

<style scoped>

</style>
