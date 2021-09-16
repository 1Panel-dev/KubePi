<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Nodes' }" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="15">
          <el-card class="el-card">
            <ko-detail-basic :item="item" :yaml-show.sync="yamlShow"></ko-detail-basic>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card class="el-card">
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" width="30%" align="left">
                  <h3>{{ $t("business.common.system") }}</h3></th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("business.node.arch") }}</td>
                <td>{{ item.status.nodeInfo.architecture }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.node.os") }}</td>
                <td>{{ item.status.nodeInfo.operatingSystem }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.node.osImage") }}</td>
                <td>{{ item.status.nodeInfo.osImage }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.node.kernel") }}</td>
                <td>{{ item.status.nodeInfo.kernelVersion }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.node.container") }}</td>
                <td>{{ item.status.nodeInfo.containerRuntimeVersion }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.node.kubelet_version") }}</td>
                <td>{{ item.status.nodeInfo.kubeletVersion }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
        <el-col :span="24" v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
          <br>
          <el-row :gutter="20">
            <el-col :span="9">
              <el-card>
                <div slot="header" class="clearfix" style="text-align: center;">
                  <span>CPU</span>
                </div>
                <div style="text-align: center;">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div>Requests</div>
                      <br>
                      <el-progress type="circle" :percentage="cpuResource.requestsUsage"></el-progress>
                      <br>
                      <span>Cores:{{ cpuResource.requests }}</span>
                    </el-col>
                    <el-col :span="12">
                      <div>Limits</div>
                      <br>
                      <el-progress type="circle"
                                   :percentage="cpuResource.limitsUsage  >= 100 ? 100: cpuResource.limitsUsage"
                                   :format="cpuFormat"></el-progress>
                      <br>
                      <span>Cores:{{ cpuResource.limits }}</span>
                    </el-col>
                  </el-row>
                </div>
              </el-card>
            </el-col>
            <el-col :span="9">
              <el-card>
                <div slot="header" class="clearfix" style="text-align: center;">
                  <span>Memory</span>
                </div>
                <div style="text-align: center">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <div>Requests</div>
                      <br>
                      <el-progress type="circle" :percentage="memResource.requestsUsage"></el-progress>
                      <br>
                      <span>Mib:{{ memResource.requests }}</span>
                    </el-col>
                    <el-col :span="12">
                      <div>Limits</div>
                      <br>
                      <el-progress type="circle"

                                   :percentage="memResource.limitsUsage  >= 100 ? 100: memResource.limitsUsage"
                                   :format="memFormat"></el-progress>
                      <br>
                      <span>Mib:{{ memResource.limits }}</span>
                    </el-col>
                  </el-row>
                </div>
              </el-card>
            </el-col>
            <el-col :span="6">
              <el-card>
                <div slot="header" class="clearfix" style="text-align: center;">
                  <span>Pods</span>
                </div>
                <div style="text-align: center">
                  <el-row :gutter="20">
                    <el-col>
                      <div>Allocation</div>
                      <br>
                      <el-progress type="circle" :percentage="podsData.usage"></el-progress>
                      <br>
                      <span>Pods:{{ podsData.podsCount }}</span>
                    </el-col>
                  </el-row>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card" v-if="Object.keys(item.metadata).length > 0">
            <el-tab-pane label="Pods" v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
              <ko-detail-pods :cluster="cluster" :field-selector="'spec.nodeName='+item.metadata.name"></ko-detail-pods>
            </el-tab-pane>
            <el-tab-pane :label="$t('business.common.conditions')">
              <ko-detail-conditions :conditions="item.status.conditions"></ko-detail-conditions>
            </el-tab-pane>
            <el-tab-pane :label="$t('business.pod.image')">
              <table style="width: 90%" class="myTable">
                <tr>
                  <th scope="col" width="70%" align="left">{{ $t("commons.table.name") }}</th>
                  <th scope="col" align="left">{{ $t("business.pod.size") }}</th>
                </tr>
                <tr v-for="(image,index) in item.status.images" v-bind:key="index">
                  <td>{{ image.names[1] ? image.names[1] : image.names[0] }}</td>
                  <td>{{ Math.round(image.sizeBytes / (1024 * 1024)) }} MB</td>
                </tr>
              </table>
            </el-tab-pane>
            <el-tab-pane :label="$t('business.pod.resource')">
              <table style="width: 100%" class="myTable">
                <tr>
                  <td>Pod CIDR</td>
                  <td>{{ item.spec.podCIDR }}</td>
                </tr>
                <tr>
                  <td>{{ $t("business.pod.address") }}</td>
                  <td>
                  <span v-for="(address,index) in item.status.addresses" v-bind:key="index">
                      {{ address.type }} : {{ address.address }}
                  </span>
                  </td>
                </tr>
              </table>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="yaml" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getNode} from "@/api/nodes"
import {listPodsWithNsSelector} from "@/api/pods"
import YamlEditor from "@/components/yaml-editor"
import KoDetailConditions from "@/components/detail/detail-conditions"
import KoDetailPods from "@/components/detail/detail-pods"
import KoDetailBasic from "@/components/detail/detail-basic"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "NodeDetail",
  components: {
    KoDetailBasic,
    KoDetailPods,
    KoDetailConditions,
    YamlEditor,
    LayoutContent
  },
  props: {
    name: String,
  },
  data () {
    return {
      loading: false,
      item: {
        metadata: {},
        status: {
          nodeInfo: {}
        },
        spec: {}
      },
      showItem: false,
      pods: [],
      cpuResource: {
        total: 0,
        limitsUsage: 0,
        limits: 0,
        requestsUsage: 0,
        requests: 0,
      },
      memResource: {
        total: 0,
        limitsUsage: 0,
        limits: 0,
        requestsUsage: 0,
        requests: 0
      },
      podsData: {
        usage: 0,
        limit: 110,
        podsCount: 0
      },
      page: {
        pageSize: 20,
        nextToken: "",
      },
      yaml: {},
      yamlShow: false,
      cluster: ""
    }
  },
  methods: {
    getNodeByName () {
      this.loading = true
      getNode(this.cluster, this.name).then(res => {
        this.item = res
        this.loading = false
        if (this.yamlShow) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
        this.listPodsByNodeName()
        console.log(this.item)
        if (this.item.status.allocatable.cpu.indexOf("m") > 0) {
          this.cpuResource.total = parseInt(this.item.status.allocatable.cpu)
        } else {
          this.cpuResource.total = parseInt(this.item.status.allocatable.cpu) / 1000
        }

        this.memResource.total = parseInt(this.item.status.allocatable.memory) / 1000
        this.podsData.limit = parseInt(this.item.status.allocatable.pods)
      })
    },
    listPodsByNodeName () {
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "list" })) {
        return
      }
      listPodsWithNsSelector(this.cluster, "", "", "spec.nodeName=" + this.item.metadata.name).then(res => {
        this.podsData.usage = Math.round(parseInt(res.items.length) / this.podsData.limit * 100)
        this.podsData.podsCount = res.items.length

        let cpuLimits = 0
        let memLimits = 0
        let cpuRequests = 0
        let memRequests = 0
        for (const pod of res.items) {
          if (pod.spec && pod.spec.containers && pod.spec.containers.length > 0) {
            for (const container of pod.spec.containers) {
              if (container.resources.limits) {
                if (container.resources.limits.cpu) {
                  const value = container.resources.limits.cpu
                  const cpu = value.toString().indexOf("m") > -1 ? parseInt(value) : parseInt(value) * 1000
                  cpuLimits = cpuLimits + cpu
                }
                if (container.resources.limits.memory) {
                  memLimits = memLimits + parseInt(container.resources.limits.memory)
                }
              }
              if (container.resources.requests) {
                if (container.resources.requests.cpu) {
                  const value = container.resources.requests.cpu
                  const cpu = value.indexOf("m") > -1 ? parseInt(value) : parseInt(value) * 1000
                  cpuRequests = cpuRequests + cpu
                }
                if (container.resources.requests.memory) {
                  memRequests = memRequests + parseInt(container.resources.requests.memory)
                }
              }
            }
          }
        }
        this.cpuResource.limits = (cpuLimits / 1000).toFixed(3)
        this.cpuResource.requests = (cpuRequests / 1000).toFixed(3)
        this.memResource.limits = memLimits
        this.memResource.requests = memRequests
        this.cpuResource.limitsUsage = Math.round(cpuLimits / this.cpuResource.total * 100)
        this.cpuResource.requestsUsage = Math.round(cpuRequests / this.cpuResource.total * 100)
        this.memResource.limitsUsage = Math.round(memLimits / this.memResource.total * 100)
        this.memResource.requestsUsage = Math.round(memRequests / this.memResource.total * 100)
      })
    },
    cpuFormat () {
      return this.cpuResource.limitsUsage + "\n %"
    },
    memFormat () {
      return this.memResource.limitsUsage + "%"
    },
  },
  watch: {
    yamlShow: function (newValue) {
      if (newValue) {
        this.yaml = JSON.parse(JSON.stringify(this.item))
      }
      this.$router.push({
        path: "/nodes/detail/" + this.name,
        query: { yamlShow: newValue }
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getNodeByName()
  }
}
</script>

<style scoped>

</style>
