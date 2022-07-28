<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Endpoints'}" v-loading="loading">
    <div v-if="!showYaml">
      <el-form label-position="top" :model="form" ref="form" :rules="rules">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-form-item :label="$t('commons.table.name')" prop="metadata.name">
              <el-input clearable v-model="form.metadata.name"></el-input>
            </el-form-item>
          </el-col>
          <el-col :span="3">
            <el-form-item :label="$t('business.namespace.namespace')" prop="metadata.namespace">
              <ko-select :namespace.sync="form.metadata.namespace"></ko-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-col :span="24">
          <div v-for="(subset,index) in form.subsets" v-bind:key="index">
            <el-tabs v-model="activeName" tab-position="top" type="border-card" @tab-click="handleClick" ref=tabs>
              <el-tab-pane name="Address" :label="$t('business.pod.address')">
                <ko-card :title="$t('business.pod.address')">
                  <table style="width: 100%;padding: 0" class="tab-table">
                    <tr>
                      <th scope="col" width="30%" align="left">
                        <label>ip</label>
                      </th>
                      <th>
                      </th>
                    </tr>
                    <tr v-for="(address,index2) in subset.addresses" v-bind:key="index2">
                      <td>
                        <el-input v-model="address.ip"></el-input>
                      </td>
                      <td>
                        <el-button type="text" style="font-size: 10px" @click="handleDelete(subset, index2)">
                          {{ $t("commons.button.delete") }}
                        </el-button>
                      </td>
                    </tr>
                    <tr>
                      <td align="left">
                        <el-button @click="handleAdd(subset)">{{ $t("commons.button.add") }}</el-button>
                      </td>
                    </tr>
                  </table>
                </ko-card>
              </el-tab-pane>
              <el-tab-pane name="Port" :label="$t('business.network.port')">
                <ko-card :title="$t('business.network.port')">
                  <table style="width: 100%;padding: 0" class="tab-table">
                    <tr>
                      <th scope="col" width="20%" align="left"><label>{{$t('business.network.port_name')}}</label></th>
                      <th scope="col" width="20%" align="left"><label>port</label></th>
                      <th scope="col" width="10%" align="left"><label>protocol</label></th>
                      <th>
                      </th>
                    </tr>
                    <tr v-for="(row,index) in subset.ports" v-bind:key="index">
                      <td>
                        <el-input v-model="row.name"></el-input>
                      </td>
                      <td>
                        <el-input v-model.number="row.port" placeholder="8080"></el-input>
                      </td>
                      <td>
                        <el-select v-model="row.protocol" style="width: 100%">
                          <el-option label="TCP" value="TCP"></el-option>
                          <el-option label="UDP" value="UDP"></el-option>
                        </el-select>
                      </td>
                      <td>
                        <el-button type="text" style="font-size: 10px" @click="handlePortDelete(subset, index)">
                          {{ $t("commons.button.delete") }}
                        </el-button>
                      </td>
                    </tr>
                    <tr>
                      <td align="left">
                        <el-button @click="handlePortAdd(subset)">{{ $t("commons.button.add") }}</el-button>
                      </td>
                    </tr>
                  </table>
                </ko-card>
              </el-tab-pane>
            </el-tabs>
          </div>
        </el-col>
      </el-form>
    </div>
    <div v-if="showYaml">
      <yaml-editor ref="yaml_editor" :value="form"></yaml-editor>
    </div>
    <div class="bottom-button">
      <el-button @click="onCancel()">{{ $t("commons.button.cancel") }}</el-button>
      <el-button v-if="!showYaml" @click="onEditYaml()">{{ $t("commons.button.yaml") }}</el-button>
      <el-button v-if="showYaml" @click="backToForm()">{{ $t("commons.button.back_form") }}</el-button>
      <el-button v-loading="loading" @click="onSubmit" type="primary">
        {{ $t("commons.button.submit") }}
      </el-button>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import YamlEditor from "@/components/yaml-editor"
import { createEndPoint } from "@/api/endpoints"
import KoSelect from "@/components/ko-select"
import Rule from "@/utils/rules"
import KoCard from "@/components/ko-card"
import { listNodes } from "@/api/nodes"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "EndpointCreate",
  components: { KoSelect, YamlEditor, LayoutContent, KoCard },
  props: {},
  data() {
    return {
      loading: false,
      form: {
        kind: "Endpoints",
        apiVersion: "v1",
        metadata: {
          name: "",
          namespace: "",
        },
        subsets: [{ addresses: [], ports: [] }],
      },
      rules: {
        metadata: {
          name: [Rule.RequiredRule],
          namespace: [Rule.RequiredRule],
        },
      },
      showYaml: false,
      yaml: undefined,
      clusterName: "",
      activeName: "Address",
      node_list: [],
    }
  },
  methods: {
    handleClick(tab) {
      this.activeName = tab.name
    },
    onCancel() {
      this.$router.push({ name: "Endpoints" })
    },
    onEditYaml() {
      this.$refs["form"].validate((valid) => {
        if (!valid) {
          this.$confirm(this.$t("commons.validate.params_not_complete") + this.$t("commons.confirm_message.open_yaml"), this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "warning",
          }).then(() => {
            this.showYaml = true
            this.yaml = this.transformYaml()
          })
        } else {
          this.showYaml = true
          this.yaml = this.transformYaml()
        }
      })
    },
    backToForm() {
      this.$confirm(this.$t("commons.confirm_message.back_form"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.showYaml = false
      })
    },
    loadNodes() {
      this.node_list = []
      if (!checkPermissions({ scope: "cluster", apiGroup: "", resource: "nodes", verb: "list" })) {
        return
      }
      listNodes(this.clusterName).then((res) => {
        for (const node of res.items) {
          this.node_list.push(node.metadata.name)
        }
      })
    },
    handleAdd(subset) {
      const item = { ip: "" }
      subset.addresses.push(item)
    },
    handleDelete(subset, index) {
      subset.addresses.splice(index, 1)
    },
    handlePortAdd(subset) {
      const item = {
        name: "",
        port: "",
        protocol: "TCP",
      }
      subset.ports.push(item)
    },
    handlePortDelete(subset, index) {
      subset.ports.splice(index, 1)
    },
    transformYaml() {
      for (const subset of this.form.subsets) {
        for (const port of subset.ports) {
          if (!port.port) {
            delete port.port
          }
          if (!port.tar) {
            delete port.hostname
          }
        }
      }
      return JSON.parse(JSON.stringify(this.form))
    },
    onSubmit() {
      if (this.showYaml) {
        this.onCreate(this.$refs.yaml_editor.getValue())
      } else {
        this.$refs["form"].validate((valid) => {
          if (valid) {
            this.onCreate(this.transformYaml())
          }
        })
      }
    },
    onCreate(data) {
      this.loading = true
      createEndPoint(this.clusterName, data.metadata.namespace, data)
        .then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.create_success"),
          })
          this.$router.push({ name: "Endpoints" })
        })
        .finally(() => {
          this.loading = false
        })
    },
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.loadNodes()
    this.showYaml = this.$route.query.yamlShow === "true"
  },
}
</script>

<style scoped>
</style>