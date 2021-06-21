<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{ name: 'Nodes' }" v-loading="loading">
    <el-row :gutter="20">
      <el-col :span="15">
        <el-card>
          <table style="width: 90%" class="myTable">
            <tr>
              <th scope="col" width="30%" align="left">
                <h3>{{ $t("business.common.basic") }}</h3></th>
              <th scope="col"></th>
            </tr>
            <tr>
              <td>{{ $t("commons.table.name") }}</td>
              <td>{{ item.metadata.name }}</td>
            </tr>
            <tr>
              <td>{{ $t("commons.table.created_time") }}</td>
              <td>{{ item.metadata.creationTimestamp | datetimeFormat }}</td>
            </tr>
            <tr>
              <td>{{ $t("business.common.label") }}</td>
              <td>
                <div v-for="(value,key,index) in item.metadata.labels" v-bind:key="index" class="myTag">
                  <el-tag type="info" size="small">
                    {{ key }} = {{ value }}
                  </el-tag>
                </div>
              </td>
            </tr>
            <tr>
              <td>{{ $t("business.common.annotation") }}</td>
              <td>
                <el-link @click="showItem=!showItem">
                  {{ showItem ? $t("business.common.pack_up") : $t("business.common.expand") }}
                </el-link>
                <div v-if="showItem">
                  <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small">
                      {{ key }} = {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                    </el-tag>
                  </div>
                </div>
              </td>
            </tr>
          </table>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
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
            <tr>
              <td>{{ $t("business.node.kubeProxy_version") }}</td>
              <td>{{ item.status.nodeInfo.kubeProxyVersion }}</td>
            </tr>
          </table>
        </el-card>
      </el-col>
      <el-col :span="23">
        <br>
        <el-card>
          <div slot="header" class="clearfix">
            <span>{{ $t("business.node.allocation") }}</span>
          </div>
          <el-row :gutter="20">
            <div style="text-align: center">
              <el-col :span="9">
                <h4>CPU</h4>
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
              </el-col>
            </div>
            <div style="text-align: center">
              <el-col :span="9">
                <h4>Memory</h4>
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
              </el-col>
            </div>
            <el-col :span="6">
              <div style="text-align: center">
                <h4>Pods</h4>
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
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getNode} from "@/api/nodes"
import {listPodsByNode} from "@/api/pods"

export default {
  name: "NodeDetail",
  components: { LayoutContent },
  props: {
    name: String,
    cluster: String,
  },
  data () {
    return {
      loading: false,
      item: {
        metadata: {},
        status: {
          nodeInfo: {}
        }
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
      }
    }
  },
  methods: {
    getNodeByName () {
      this.loading = true
      getNode(this.cluster, this.name).then(res => {
        this.item = res
        this.loading = false
        this.listPodsByNodeName(this.item.metadata.name)

        this.cpuResource.total = parseInt(this.item.status.allocatable.cpu)
        this.memResource.total = parseInt(this.item.status.allocatable.memory) / 1000
        this.podsData.limit = parseInt(this.item.status.allocatable.pods)
      })
    },
    listPodsByNodeName (nodeName) {
      listPodsByNode(this.cluster, nodeName).then(res => {
        this.pods = res.items
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
      return this.cpuResource.limitsUsage + "%"
    },
    memFormat () {
      return this.memResource.limitsUsage + "%"
    }
  },
  created () {
    this.getNodeByName()
  }
}
</script>

<style scoped>

</style>