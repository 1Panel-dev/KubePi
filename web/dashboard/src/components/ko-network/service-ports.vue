<template>
  <div style="margin-top: 20px">
    <ko-card title="Service Ports">
      <table style="width: 100%;padding: 0" class="tab-table">
        <tr>
          <th scope="col" width="30%" align="left">
            <label>Port Name</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>Listening Port</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>Protocol</label>
          </th>
          <th scope="col" width="20%" align="left">
            <label>Target Port</label>
          </th>
          <th></th>
        </tr>
        <tr v-for="(row,index) in servicePorts" v-bind:key="index">
          <td>
            <el-input v-model="row.name" @change="transformation"></el-input>
          </td>
          <td>
            <el-input v-model.number="row.port" @change="transformation"></el-input>
          </td>
          <td>
            <el-select v-model="row.protocol" style="width: 100%" @change="transformation">
              <el-option label="TCP" value="TCP"></el-option>
              <el-option label="UDP" value="UDP"></el-option>
            </el-select>
          </td>
          <td>
            <el-input type="text" v-model="row.targetPort" @change="transformation(index)"></el-input>
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
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"

export default {
  name: "KoServicePorts",
  components: { KoCard },
  props: {
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
        port: 8080,
        protocol: "TCP",
        targetPort: ""
      }
      this.servicePorts.push(item)
    },
    handleDelete (index) {
      this.servicePorts.splice(index, 1)
      this.transformation()
    },
    transformation (index) {
      if (index !== undefined && !Number.isNaN(Number(this.servicePorts[index].targetPort))) {
        this.servicePorts[index].targetPort = parseInt(this.servicePorts[index].targetPort)
      }
      this.$emit("update:ports", this.servicePorts)
    }
  },
  created () {
    if (!this.ports) {
      this.servicePorts.push({
        name: "",
        port: 8080,
        protocol: "TCP",
        targetPort: ""
      })
    } else {
      this.servicePorts = this.ports
    }
  }
}
</script>

<style scoped>

</style>
