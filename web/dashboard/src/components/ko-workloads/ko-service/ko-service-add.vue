<template>
  <div style="margin-left: 20px">
    <el-form label-position="top" :model="form" ref="form">
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item :label="$t('business.configuration.type')">
            <ko-form-item itemType="radio" v-model="serviceType" :radios="service_type_list" />
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20" v-if="serviceType !== 'None'">
        <el-col :span="12">
          <el-form-item :label="$t('business.workload.annotations')">
            <table style="width: 100%" class="tab-table">
              <tr v-for="(row,index) in annotations" :key="index">
                <td width="49%">
                  <ko-form-item placeholder="e.g. foo" clearable itemType="input" v-model="row.key" />
                </td>
                <td width="49%">
                  <ko-form-item clearable itemType="input" v-model="row.value" />
                </td>
                <td width="2%">
                  <el-button type="text" style="font-size: 10px" @click="handleAnnoDelete(index)">
                    {{ $t("commons.button.delete") }}
                  </el-button>
                </td>
              </tr>
              <tr>
                <td align="left">
                  <el-button @click="handleAnnoAdd()">{{ $t("commons.button.add") }}</el-button>
                </td>
              </tr>
            </table>
          </el-form-item>
        </el-col>
      </el-row>

      <div v-if="hasExternalIPs()">
        <el-row :gutter="20">
          <table style="width: 100%;" class="tab-table">
            <tr>
              <th scope="col" width="20%" align="left"><label>{{$t('business.network.port_name')}}</label></th>
              <th scope="col" width="20%" align="left"><label>{{$t('business.workload.service_port')}}</label></th>
              <th scope="col" width="10%" align="left"><label>{{$t('business.network.protocol')}}</label></th>
              <th scope="col" width="20%" align="left"><label>{{$t('business.workload.container_port')}}</label></th>
              <th scope="col" width="20%" align="left"><label>{{$t('business.workload.node_port')}}</label></th>
              <th />
            </tr>
            <tr v-for="(row,index) in ports" v-bind:key="index">
              <td>
                <ko-form-item itemType="input" v-model="row.name" />
              </td>
              <td>
                <ko-form-item :placeholder="$t('business.workload.service_port_help')" itemType="number" v-model.number="row.port" />
              </td>
              <td>
                <ko-form-item itemType="select" v-model="row.protocol" :selections="protocol_list" />
              </td>
              <td>
                <ko-form-item :placeholder="$t('business.workload.application_port')" itemType="number" v-model.number="row.targetPort" />
              </td>
              <th>
                <ko-form-item v-if="serviceType==='NodePort'" itemType="number" v-model.number="row.nodePort" />
              </th>
              <td>
                <el-button type="text" style="font-size: 10px" @click="handleDelete(index)">{{ $t("commons.button.delete") }}</el-button>
              </td>
            </tr>
            <tr>
              <td align="left">
                <el-button @click="handleAdd">{{ $t("commons.button.add") }}</el-button>
              </td>
            </tr>
          </table>
        </el-row>

        <el-row :gutter="20" style="margin-top: 20px">
          <el-col :span="6">
            <el-form-item label="SessionAffinity">
              <ko-form-item itemType="radio" v-model="form.spec.sessionAffinity" :radios="session_type_list" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="6" v-if="form.spec.sessionAffinity === 'ClientIP'">
            <el-form-item label="TimeoutSeconds">
              <ko-form-item itemType="number" v-model.number="form.spec.sessionAffinityConfig.clientIP.timeoutSeconds" :deviderName="$t('business.workload.seconds')" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span=12>
            <span>{{$t('business.network.external_ip')}}</span>
            <table style="width: 100%;margin-top:10px" class="tab-table">
              <tr v-for="(row,index) in externalIPs" :key="index">
                <td width="97%">
                  <ko-form-item placeholder="e.g. 1.1.1.1" itemType="input" v-model="row.value" />
                </td>
                <td>
                  <el-button type="text" style="font-size: 10px" @click="handleExternalIPsDelete(index)">
                    {{ $t("commons.button.delete") }}
                  </el-button>
                </td>
              </tr>
              <tr>
                <td align="left">
                  <el-button @click="handleExternalIPsAdd">{{$t('business.workload.add')}}{{$t('business.network.external_ip')}}</el-button>
                </td>
              </tr>
            </table>
          </el-col>
        </el-row>
      </div>
    </el-form>
  </div>
</template>

<script>
import KoFormItem from "@/components/ko-form-item/index"
import { parseArryToObj, parseObjToArry } from "@/utils/objArryParse"

export default {
  name: "KoServiceAdd",
  components: { KoFormItem },
  props: {
    serviceObj: Object,
  },
  watch: {
    serviceObj: {
      handler(newObj) {
        if (newObj) {
          this.form = newObj
          if (this.form.spec?.type) {
            this.serviceType = this.form.spec.type
          }
          if (this.form.metadata?.annotations) {
            this.annotations = parseObjToArry(this.form.metadata.annotations)
          }
          this.externalIPs = []
          if (this.form.spec?.externalIPs) {
            for (const ip of this.form.spec.externalIPs) {
              this.externalIPs.push({ value: ip })
            }
          }
          this.ports = []
          if (this.form.spec?.ports) {
            for (const item of this.form.spec.ports) {
              var itemPo = {}
              itemPo.name = item.name || undefined
              itemPo.port = item.port || undefined
              itemPo.protocol = item.protocol || undefined
              itemPo.targetPort = item.targetPort || undefined
              itemPo.nodePort = item.nodePort || undefined
              this.ports.push({
                name: item.name,
                port: item.port,
                protocol: item.protocol,
                targetPort: item.targetPort,
                nodePort: item.nodePort,
              })
            }
          }
        }
      },
      immediate: true,
      deep: true,
    },
  },
  data() {
    return {
      cluster: "",
      loading: false,
      annotations: [],
      externalIPs: [],
      ports: [],
      serviceType: "None",
      form: {
        apiVersion: "v1",
        kind: "Service",
        metadata: {
          name: "",
          namespace: "",
          labels: {},
          annotations: {},
        },
        spec: {
          type: "None",
          ports: [],
          sessionAffinity: "None",
          sessionAffinityConfig: {
            clientIP: {
              timeoutSeconds: null,
            },
          },
          externalIPs: [],
        },
      },
      service_type_list: [
        { label: this.$t("business.workload.not_create"), value: "None" },
        { label: "Headless", value: "Headless" },
        { label: "ClusterIP", value: "ClusterIP" },
        { label: "NodePort", value: "NodePort" },
      ],
      session_type_list: [
        { label: "None", value: "None" },
        { label: "ClientIP", value: "ClientIP" },
      ],
      protocol_list: [
        { label: "TCP", value: "TCP" },
        { label: "UDP", value: "UDP" },
      ],
    }
  },
  methods: {
    handleExternalIPsAdd() {
      const item = { value: "" }
      this.externalIPs.push(item)
    },
    handleExternalIPsDelete(index) {
      this.externalIPs.splice(index, 1)
    },
    handleAnnoAdd() {
      const item = {
        key: "",
        value: "",
      }
      this.annotations.push(item)
    },
    handleAnnoDelete(index) {
      this.annotations.splice(index, 1)
    },
    handleAdd() {
      const item = {
        name: "",
        port: "",
        protocol: "TCP",
        targetPort: "",
        nodePort: "",
      }
      this.ports.push(item)
    },
    handleDelete(index) {
      this.ports.splice(index, 1)
    },
    hasExternalIPs() {
      return this.serviceType === "ClusterIP" || this.serviceType === "NodePort"
    },
    isServiceCreate() {
      return !this.serviceType === "None"
    },
    transformation(metadata) {
      if (this.serviceType === "None") {
        return null
      }
      this.form.metadata.annotations = parseArryToObj(this.annotations)
      for (const item of this.externalIPs) {
        this.form.spec.externalIPs.push(item.value)
      }
      this.form.spec.ports = []
      for (const item of this.ports) {
        var itemPo = {}
        itemPo.name = item.name || undefined
        itemPo.port = item.port || undefined
        itemPo.protocol = item.protocol || undefined
        itemPo.targetPort = item.targetPort || undefined
        itemPo.nodePort = item.nodePort || undefined
        this.form.spec.ports.push(itemPo)
      }
      if (this.serviceType !== "None") {
        this.form.spec.type = this.serviceType
      }
      this.form.metadata.name = metadata.name
      this.form.metadata.namespace = metadata.namespace
      this.form.spec.selector = metadata.labels
      return this.form
    },
  },
}
</script>

<style scoped>
</style>
