<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.network.default_backend')">
      <el-form-item label-position="top" ref="form" :model="defaultBackend">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.network.target')+'Service'">
              <el-select v-model="defaultBackend.service.name"
                         @change="changeService(defaultBackend.service.name)"
                         style="width: 100%">
                <el-option v-for="service in services" :key="service.metadata.name"
                           :label="service.metadata.name" :value="service.metadata.name">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.network.port')">
              <el-select v-model="defaultBackend.service.port.number">
                <el-option v-for="port in servicePorts" :key="port.name"
                           :label="port.port" :value="port.port">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form-item>
    </ko-card>
  </div>
</template>

<script>
import KoCard from "@/components/ko-card"
import {listNsServices} from "@/api/services"

export default {
  name: "KoIngressDefaultBackend",
  components: { KoCard },
  props: {
    namespace: String,
    cluster: String,
    defaultBackendObj: Object,
  },
  data () {
    return {
      defaultBackend: {
        service: {}
      },
      services: [],
      servicePorts: []
    }
  },
  methods: {
    getServices () {
      listNsServices(this.cluster, this.namespace).then(res => {
        this.services = res.items
      })
    },
    changeService (serviceName) {
      for (const service of this.services) {
        if (service.metadata.name === serviceName) {
          this.servicePorts = service.spec.ports
          break
        }
      }
    },
  },
  watch: {
    namespace: function (old, newValue) {
      if (old !== newValue) {
        this.defaultBackend.service.name = ""
        this.defaultBackend.service.port.number = 0
        this.getServices()
      }
    }
  },
  created () {
    this.defaultBackend = this.defaultBackendObj ? this.defaultBackendObj : {
      service: {
        name: "",
        port: {
          number: 0
        }
      }
    }
    this.$emit("update:defaultBackendObj", this.defaultBackend)
    this.getServices()
  }
}
</script>

<style scoped>

</style>
