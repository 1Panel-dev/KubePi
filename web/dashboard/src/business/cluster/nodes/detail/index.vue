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
                <el-link @click="showItem=!showItem"> {{showItem? $t("business.common.pack_up"):$t("business.common.expand")}}</el-link>
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
        <el-card></el-card>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getNode} from "@/api/nodes"

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
          nodeInfo:{}
        }
      },
      showItem: false
    }
  },
  methods: {
    getNodeByName () {
      this.loading = true
      getNode(this.cluster, this.name).then(res => {
        this.item = res
        this.loading = false
      })
    },
  },
  created () {
    this.getNodeByName()
  }
}
</script>

<style scoped>

</style>