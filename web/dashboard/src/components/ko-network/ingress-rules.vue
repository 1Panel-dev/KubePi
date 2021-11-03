<template>
  <div style="margin-top: 20px">
    <ko-card :title="$t('business.network.rule')">
      <el-form label-position="top" ref="form" :model="form">
        <el-card v-for="(row,ind) in form.rules" v-bind:key="ind"
                 style="background-color: #292a2e;margin-top: 10px;position: relative">
          <div style="float: right; padding: 3px 0;position: relative;z-index: 1">
            <el-button type="text" v-if="form.rules.length > 1"
                       @click="removeRule(ind)">{{ $t("commons.button.delete") }}
            </el-button>
          </div>
          <div>
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item :label="$t('business.network.host')" prop="host">
                  <el-input placeholder="e.g. example.com" v-model="row.host"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="24">
                <table style="width: 100%;padding: 0" class="tab-table">
                  <tr>
                    <th scope="col" width="10%" align="left">
                      <label>{{ $t("business.network.path") }}</label>
                    </th>
                    <th scope="col" width="38.5%" align="left"></th>
                    <th scope="col" width="20%" align="left">
                      <label>{{ $t("business.network.target") }} Service</label>
                    </th>
                    <th scope="col" width="10%" align="left">
                      <label>{{ $t("business.network.port") }}</label>
                    </th>
                    <th></th>
                  </tr>
                  <tr v-for="(row,index) in row.http.paths" v-bind:key="index">
                    <td>
                      <el-select v-model="row.pathType">
                        <el-option label="Prefix" value="Prefix"></el-option>
                        <el-option label="Exact" value="Exact"></el-option>
                        <el-option label="ImplementationSpecific" value="ImplementationSpecific"></el-option>
                      </el-select>
                    </td>
                    <td>
                      <el-input v-model="row.path" @change="transformation"></el-input>
                    </td>
                    <td v-if="row.backend.serviceName!==undefined">
                      <el-select v-model="row.backend.serviceName"
                                 @change="changeService(row.backend.serviceName,index)"
                                 style="width: 100%">
                        <el-option v-for="service in services" :key="service.metadata.name"
                                   :label="service.metadata.name" :value="service.metadata.name">
                        </el-option>
                      </el-select>
                    </td>
                    <td v-else>
                      <el-select v-model="row.backend.service.name"
                                 @change="changeService(row.backend.service.name,index)"
                                 style="width: 100%">
                        <el-option v-for="service in services" :key="service.metadata.name"
                                   :label="service.metadata.name" :value="service.metadata.name">
                        </el-option>
                      </el-select>
                    </td>

                    <td v-if="row.backend.servicePort !== undefined">
                      <el-select v-model="row.backend.servicePort">
                        <el-option v-for="port in servicePorts[index]" :key="port.name"
                                   :label="port.port" :value="port.port">
                        </el-option>
                      </el-select>
                    </td>
                    <td v-else>
                      <el-select v-model="row.backend.service.port.number">
                        <el-option v-for="port in servicePorts[index]" :key="port.name"
                                   :label="port.port" :value="port.port">
                        </el-option>
                      </el-select>
                    </td>
                    <td>
                      <el-button type="text" style="font-size: 10px" @click="deletePath(ind,index)">
                        {{ $t("commons.button.delete") }}
                      </el-button>
                    </td>
                  </tr>
                  <tr>
                    <td align="left">
                      <el-button @click="addPath(row)">{{
                          $t("commons.button.add")
                        }}{{ $t("business.network.path") }}
                      </el-button>
                    </td>
                  </tr>
                </table>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-form>
      <el-button @click="addRule" style="margin-top: 10px">{{
          $t("commons.button.add")
        }}{{ $t("business.network.rule") }}
      </el-button>
    </ko-card>
  </div>
</template>

<script>

import KoCard from "@/components/ko-card"
import {listServicesWithNs} from "@/api/services"
import {set} from "@/utils/object"

export default {
  name: "KoIngressRule",
  components: { KoCard },
  props: {
    rulesArray: Array,
    namespace: String,
    cluster: String,
    mode: String,
    newVersion: Boolean
  },
  data () {
    return {
      form: {
        rules: [],
      },
      services: [],
      servicePorts: {},
      servicePath: "",
      portPath: ""
    }
  },
  methods: {
    removeRule (index) {
      this.form.rules.splice(index, 1)
    },
    transformation () {
      this.$emit("update:rulesArray", this.form.rules)
    },
    changeService (serviceName, index) {
      for (const service of this.services) {
        if (service.metadata.name === serviceName) {
          this.servicePorts[index] = service.spec.ports
          break
        }
      }
    },
    addRule () {
      this.form.rules.push(this.getDefaultRule())
    },
    addPath (row) {
      const path = {
        backend: {},
        path: "",
        pathType: ""
      }
      set(path.backend, this.servicePath, "")
      set(path.backend, this.portPath, 0)
      row.http.paths.push(path)
    },
    deletePath (ind, index) {
      this.form.rules[ind].http.paths.splice(index, 1)
    },
    getServices () {
      listServicesWithNs(this.cluster, this.namespace).then(res => {
        this.services = res.items
      })
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
    getDefaultRule () {
      const rule = {
        host: "",
        http: {
          paths: [{
            backend: {},
            path: "",
            pathType: ""
          }]
        }
      }
      set(rule.http.paths[0].backend, this.servicePath, "")
      set(rule.http.paths[0].backend, this.portPath, 0)
      return rule
    }
  },
  watch: {
    namespace: function (old, newValue) {
      if (old !== newValue) {
        for (let i = 0; i < this.form.rules.length; i++) {
          for (let j = 0; j < this.form.rules[i].http.paths.length; j++) {
            set(this.form.rules[i].http.paths[j].backend, this.servicePath, "")
            set(this.form.rules[i].http.paths[j].backend, this.portPath, 0)
          }
        }
        this.getServices()
      }
    }
  },
  created () {
    this.servicePath = this.serviceNamePath()
    this.portPath = this.servicePortPath()

    if (this.rulesArray) {
      this.form.rules = this.rulesArray
    } else {
      this.form.rules = [this.getDefaultRule()]
    }
    this.transformation()
    this.getServices()
  }
}
</script>

<style scoped>

</style>
