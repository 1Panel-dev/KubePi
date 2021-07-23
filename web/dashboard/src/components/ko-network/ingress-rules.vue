<template>
  <div style="margin-top: 20px">
    <ko-card title="Rules">
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
                <el-form-item label="Request Host" prop="host">
                  <el-input placeholder="e.g. example.com" v-model="row.host"></el-input>
                </el-form-item>
              </el-col>
            </el-row>
            <el-row :gutter="20">
              <el-col :span="24">
                <table style="width: 100%;padding: 0" class="tab-table">
                  <tr>
                    <th scope="col" width="10%" align="left">
                      <label>Path</label>
                    </th>
                    <th scope="col" width="38.5%" align="left"></th>
                    <th scope="col" width="20%" align="left">
                      <label>Target Service</label>
                    </th>
                    <th scope="col" width="10%" align="left">
                      <label>Port</label>
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
                    <td>
                      <el-select v-model="row.backend.service.name"
                                 @change="changeService(row.backend.service.name,index)"
                                 style="width: 100%">
                        <el-option v-for="service in services" :key="service.metadata.name"
                                   :label="service.metadata.name" :value="service.metadata.name">
                        </el-option>
                      </el-select>
                    </td>
                    <td>
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
                      <el-button @click="addPath(row)">{{ $t("commons.button.add") }} Path</el-button>
                    </td>
                  </tr>
                </table>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-form>
      <el-button @click="addRule" style="margin-top: 10px">{{ $t("commons.button.add") }} Rule</el-button>
    </ko-card>
  </div>
</template>

<script>

import KoCard from "@/components/ko-card"
import {listNsServices} from "@/api/services"

export default {
  name: "KoIngressRule",
  components: { KoCard },
  props: {
    rulesArray: Array,
    namespace: String,
    cluster: String
  },
  data () {
    return {
      form: {
        rules: [],
      },
      services: [],
      servicePorts: {}
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
      this.form.rules.push({
        host: "",
        http: {
          paths: [{
            backend: {
              service: {
                port: {
                  number: 0
                },
                name: ""
              }
            },
            path: "",
            pathType: ""
          }]
        }
      })
    },
    addPath (row) {
      row.http.paths.push({
        backend: {
          service: {
            port: {
              number: 0
            },
            name: ""
          }
        },
        path: "",
        pathType: ""
      })
    },
    deletePath (ind, index) {
      this.form.rules[ind].http.paths.splice(index, 1)
    },
    getServices () {
      listNsServices(this.cluster, this.namespace).then(res => {
        this.services = res.items
      })
    }
  },
  watch: {
    namespace: function (old, newValue) {
      if (old !== newValue) {
        this.getServices()
      }
    }
  },
  created () {
    this.form.rules = this.rulesArray ? this.rulesArray : [{
      host: "",
      http: {
        paths: [{
          backend: {
            service: {
              port: {
                number: 0
              },
              name: ""
            }
          },
          path: "",
          pathType: ""
        }]
      }
    }]
    this.transformation()
    this.getServices()
  }
}
</script>

<style scoped>

</style>
