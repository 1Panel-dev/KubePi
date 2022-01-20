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
                  <h3>{{ $t("business.common.system") }}</h3>
                </th>
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
      </el-row>
      <br>
      <el-row class="row-box" v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
        <el-card class="el-card">
          <el-row :gutter="20">
            <h4>{{ $t('business.node.allocata') }}</h4>
            <el-form label-position="top">
              <el-col :span="4">
                <el-form-item :label="'CPU ' + $t('business.workload.reservation')">
                  <span>{{ cpuResource.requests }} core ({{ cpuResource.requestsUsage }}%)</span>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item :label="'CPU ' + $t('business.workload.limit')">
                  <span>{{ cpuResource.limits }} core ({{ cpuResource.limitsUsage }}%)</span>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item :label="$t('business.workload.memory') + $t('business.workload.reservation')">
                  <span>{{ memResource.requests }} Mi ({{ memResource.requestsUsage }}%)</span>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item :label="$t('business.workload.memory') + $t('business.workload.limit')">
                  <span>{{ memResource.limits }} Mi ({{ memResource.limitsUsage }}%)</span>
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-form-item label="Pods">
                  <span>{{ podsData.podsCount }} ({{ podsData.usage }}%)</span>
                </el-form-item>
              </el-col>
            </el-form>
          </el-row>
        </el-card>
      </el-row>
      <br>
      <el-row class="row-box" v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
        <el-card class="el-card">
          <el-row :gutter="20">
            <h4>{{$t('business.node.health_statu')}}</h4>
            <el-form label-position="top">
              <div v-for="(cond, index) in item.status.conditions" v-bind:key="index">
                <el-col :span="8" v-if="cond.type === 'NetworkUnavailable'">
                  <div>
                    <div style="font-size: 30px;"><i class="iconfont iconnetwork" style="float: left"></i></div>
                    <div v-if="cond.status === 'False'" style="font-size: 10px; color: green; "><i class="el-icon-success" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                    <div v-if="cond.status === 'True'" style="font-size: 10px; color: red; "><i class="el-icon-warning" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                  </div>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div><span style="font-wight: blod">{{ cond.reason }}</span></div>
                      <div style="margin-top: 5px"><span>{{ cond.message }}</span></div>
                      <div style="margin-top: 2px"><span>{{ $t('business.pod.lastUpdateTime' )}}: {{ cond.lastTransitionTime | datetimeFormat }}</span></div>
                    </div>
                    <el-form-item style="float: left; margin-left: 8px" :label="$t('business.node.network_statu')">
                      <span>{{ $t('business.node.network_statu_help') }}</span>
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="8" v-if="cond.type === 'MemoryPressure'">
                  <div>
                    <div style="font-size: 30px;"><i class="iconfont iconxingzhuang" style="float: left"></i></div>
                    <div v-if="cond.status === 'False'" style="font-size: 10px; color: green; "><i class="el-icon-success" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                    <div v-if="cond.status === 'True'" style="font-size: 10px; color: red; "><i class="el-icon-warning" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                  </div>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div><span style="font-wight: blod">{{ cond.reason }}</span></div>
                      <div style="margin-top: 5px"><span>{{ cond.message }}</span></div>
                      <div style="margin-top: 2px"><span>{{ $t('business.pod.lastUpdateTime' )}}: {{ cond.lastTransitionTime | datetimeFormat }}</span></div>
                    </div>
                    <el-form-item style="float: left; margin-left: 8px" :label="$t('business.node.memory_statu')">
                      <span>{{ $t('business.node.memory_statu_help') }}</span>
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="8" v-if="cond.type === 'DiskPressure'">
                  <div>
                    <div style="font-size: 30px;"><i class="iconfont iconcipan" style="float: left"></i></div>
                    <div v-if="cond.status === 'False'" style="font-size: 10px; color: green; "><i class="el-icon-success" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                    <div v-if="cond.status === 'True'" style="font-size: 10px; color: red; "><i class="el-icon-warning" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                  </div>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div><span style="font-wight: blod">{{ cond.reason }}</span></div>
                      <div style="margin-top: 5px"><span>{{ cond.message }}</span></div>
                      <div style="margin-top: 2px"><span>{{ $t('business.pod.lastUpdateTime' )}}: {{ cond.lastTransitionTime | datetimeFormat }}</span></div>
                    </div>
                    <el-form-item style="float: left; margin-left: 8px" :label="$t('business.node.disk_statu')">
                      <span>{{ $t('business.node.disk_statu_help') }}</span>
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="8" v-if="cond.type === 'PIDPressure'">
                  <div>
                    <div style="font-size: 30px;"><i class="iconfont iconchakanjincheng" style="float: left"></i></div>
                    <div v-if="cond.status === 'False'" style="font-size: 10px; color: green; "><i class="el-icon-success" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                    <div v-if="cond.status === 'True'" style="font-size: 10px; color: red; "><i class="el-icon-warning" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                  </div>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div><span style="font-wight: blod">{{ cond.reason }}</span></div>
                      <div style="margin-top: 5px"><span>{{ cond.message }}</span></div>
                      <div style="margin-top: 2px"><span>{{ $t('business.pod.lastUpdateTime' )}}: {{ cond.lastTransitionTime | datetimeFormat }}</span></div>
                    </div>
                    <el-form-item style="float: left; margin-left: 8px" :label="$t('business.node.pid_statu')">
                      <span>{{ $t('business.node.pid_statu_help') }}</span>
                    </el-form-item>
                  </el-tooltip>
                </el-col>
                <el-col :span="8" v-if="cond.type === 'Ready'">
                  <div>
                    <div style="font-size: 30px;"><i class="iconfont iconjiedian" style="float: left"></i></div>
                    <div v-if="cond.status === 'True'" style="font-size: 10px; color: green; "><i class="el-icon-success" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                    <div v-if="cond.status === 'False'" style="font-size: 10px; color: red; "><i class="el-icon-warning" style="float: left;margin-top: 25px; margin-left: -5px"></i></div>
                  </div>
                  <el-tooltip class="item" effect="dark" placement="top">
                    <div slot="content">
                      <div><span style="font-wight: blod">{{ cond.reason }}</span></div>
                      <div style="margin-top: 5px"><span>{{ cond.message }}</span></div>
                      <div style="margin-top: 2px"><span>{{ $t('business.pod.lastUpdateTime' )}}: {{ cond.lastTransitionTime | datetimeFormat }}</span></div>
                    </div>
                    <el-form-item style="float: left; margin-left: 8px" :label="$t('business.node.node_statu')">
                      <span>{{ $t('business.node.node_statu_help') }}</span>
                    </el-form-item>
                  </el-tooltip>
                </el-col>
              </div>
            </el-form>
          </el-row>
        </el-card>
      </el-row>
      <el-row>
        <br>
        <el-tabs type="border-card" v-if="Object.keys(item.metadata).length > 0">
          <el-tab-pane label="Pods" v-has-permissions="{scope:'namespace',apiGroup:'',resource:'pods',verb:'list'}">
            <ko-detail-node-pods :cluster="cluster" :allocatable="item.status.allocatable" :field-selector="'spec.nodeName='+item.metadata.name"></ko-detail-node-pods>
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
                  <div v-for="(address,index) in item.status.addresses" v-bind:key="index">
                    <el-tag type="success">{{ address.type }} : {{ address.address }}</el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>Allocatable</td>
                <td>
                  <div>
                    <el-tag type="success">CPU : {{item.status.allocatable.cpu}}</el-tag>
                  </div>
                  <div>
                    <el-tag type="success">Memory : {{item.status.allocatable.memory}}</el-tag>
                  </div>
                  <div>
                    <el-tag type="success">Pods : {{item.status.allocatable.pods}}</el-tag>
                  </div>
                </td>
              </tr>
              <tr>
                <td>Capacity</td>
                <td>
                  <div>
                    <el-tag type="success">CPU : {{item.status.capacity.cpu}}</el-tag>
                  </div>
                  <div>
                    <el-tag type="success">Memory : {{item.status.capacity.memory}}</el-tag>
                  </div>
                  <div>
                    <el-tag type="success">Pods : {{item.status.capacity.pods}}</el-tag>
                  </div>
                </td>
              </tr>
            </table>
          </el-tab-pane>
        </el-tabs>
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
import { getNode } from "@/api/nodes"
import { listPodsWithNsSelector } from "@/api/pods"
import YamlEditor from "@/components/yaml-editor"
import KoDetailNodePods from "@/components/detail/detail-node-pods"
import KoDetailBasic from "@/components/detail/detail-basic"
import { checkPermissions } from "@/utils/permission"

export default {
  name: "NodeDetail",
  components: { KoDetailBasic, KoDetailNodePods, YamlEditor, LayoutContent },
  props: {
    name: String,
  },
  data() {
    return {
      loading: false,
      item: {
        metadata: {},
        status: {
          nodeInfo: {},
        },
        spec: {},
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
        requests: 0,
      },
      podsData: {
        usage: 0,
        limit: 110,
        podsCount: 0,
      },
      hasMetric: null,
      metricResource: {
        cpu: "",
        cpuUsage: "",
        memory: "",
        memoryUsage: "",
      },
      page: {
        pageSize: 20,
        nextToken: "",
      },
      yaml: {},
      yamlShow: false,
      cluster: "",
    }
  },
  methods: {
    getNodeByName() {
      this.loading = true
      getNode(this.cluster, this.name).then((res) => {
        this.item = res
        this.loading = false
        if (this.yamlShow) {
          this.yaml = JSON.parse(JSON.stringify(this.item))
        }
        this.listPodsByNodeName()
        this.cpuResource.total = this.item.status.allocatable.cpu.indexOf("m") > -1 ? parseInt(this.item.status.allocatable.cpu) : parseInt(this.item.status.allocatable.cpu) * 1000
        this.memResource.total = parseInt(this.item.status.allocatable.memory) / 1000
        this.podsData.limit = parseInt(this.item.status.allocatable.pods)
      })
    },

    listPodsByNodeName() {
      if (!checkPermissions({ scope: "namespace", apiGroup: "", resource: "pods", verb: "list" })) {
        return
      }
      listPodsWithNsSelector(this.cluster, "", "", "spec.nodeName=" + this.item.metadata.name).then((res) => {
        this.podsData.usage = Math.round((parseInt(res.items.length) / this.podsData.limit) * 100)
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
        this.cpuResource.limitsUsage = Math.floor((cpuLimits / this.cpuResource.total) * 100)
        this.cpuResource.requestsUsage = Math.floor((cpuRequests / this.cpuResource.total) * 100)
        this.memResource.limitsUsage = Math.floor((memLimits / this.memResource.total) * 100)
        this.memResource.requestsUsage = Math.floor((memRequests / this.memResource.total) * 100)
      })
    },
  },
  watch: {
    yamlShow: function (newValue) {
      if (newValue) {
        this.yaml = JSON.parse(JSON.stringify(this.item))
      }
      this.$router.push({
        path: "/nodes/detail/" + this.name,
        query: { yamlShow: newValue },
      })
    },
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getNodeByName()
  },
}
</script>

<style scoped>
</style>
