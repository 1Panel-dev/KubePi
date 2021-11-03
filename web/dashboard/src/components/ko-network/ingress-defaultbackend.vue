<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.network.default_backend')">
      <el-form-item label-position="top" ref="form" :model="defaultBackend">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="$t('business.network.target')+'Service'">
              <el-select v-model="serviceName"
                         @change="changeService(serviceName)"
                         style="width: 100%">
                <el-option v-for="service in services" :key="service.metadata.name"
                           :label="service.metadata.name" :value="service.metadata.name">
                </el-option>
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="$t('business.network.port')">
              <el-select v-model="servicePort">
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
import {listServicesWithNs} from "@/api/services"
import {set} from "@/utils/object"

export default {
  name: "KoIngressDefaultBackend",
  components: { KoCard },
  props: {
    namespace: String,
    cluster: String,
    defaultBackendObj: Object,
    newVersion: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {
      defaultBackend: {
        service: {}
      },
      services: [],
      servicePorts: [],
      servicePath: "",
      portPath: "",
      serviceName: "",
      servicePort: 0,
    }
  },
  methods: {
    getServices () {
      listServicesWithNs(this.cluster, this.namespace).then(res => {
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
    serviceNamePath () {
      const nestedPath = "service.name"
      const flatPath = "serviceName"
      return this.newVersion ? nestedPath : flatPath
    },
    servicePortPath () {
      const nestedPath = "service.port.number"
      const flatPath = "servicePort"
      return this.newVersion ? nestedPath : flatPath
    },
    getDefaultBackend () {
      const backend = {}
      set(backend, this.servicePath, "")
      set(backend, this.portPath, 0)
      return backend
    }
  },
  watch: {
    namespace: function (old, newValue) {
      if (old !== newValue) {
        set(this.defaultBackend, this.servicePath, "")
        set(this.defaultBackend, this.portPath, 0)
        this.getServices()
      }
    },
    serviceName:function (value) {
      set(this.defaultBackend, this.servicePath, value)
    },
    servicePort: function (value) {
      set(this.defaultBackend, this.portPath, value)
    }
  },
  created () {
    this.servicePath = this.serviceNamePath()
    this.portPath = this.servicePortPath()
    this.defaultBackend = this.defaultBackendObj ? this.defaultBackendObj : this.getDefaultBackend()
    this.$emit("update:defaultBackendObj", this.defaultBackend)
    this.getServices()
    console.log(this.servicePath)
    console.log(this.newVersion)
    console.log(this.portPath)
  }
}
</script>

<style scoped>

</style>
