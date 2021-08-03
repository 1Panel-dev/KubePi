<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Pods'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-card>
        <table style="width: 100%" class="myTable">
          <tr>
            <th scope="col" align="left">
              <h3>{{ $t("business.common.basic") }}</h3>
            </th>
            <th scope="col"></th>
          </tr>
          <tr>
            <td>{{ $t("commons.table.name") }}</td>
            <td colspan="2">{{ form.metadata.name }}</td>
          </tr>
          <tr>
            <td>{{ $t("business.namespace.namespace") }}</td>
            <td>{{ form.metadata.namespace }}</td>
          </tr>
          <tr>
            <td>Pod IP</td>
            <td>{{ form.status.podIP }}</td>
          </tr>
          <tr>
            <td>{{ $t("business.cluster.nodes") }}</td>
            <td>{{ form.spec.nodeName }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.created_time") }}</td>
            <td>{{ form.metadata.creationTimestamp | age }}</td>
          </tr>
          <tr>
            <td>{{ $t("business.common.label") }}</td>
            <td colspan="4">
              <div v-for="(value,key,index) in form.metadata.labels" v-bind:key="index" class="myTag">
                <el-tag type="info" size="small">
                  {{ key }} = {{ value }}
                </el-tag>
              </div>
            </td>
          </tr>
          <tr>
            <td>{{ $t("business.common.annotation") }}</td>
            <td colspan="4">
              <div v-for="(value,key,index) in form.metadata.annotations" v-bind:key="index" class="myTag">
                <el-tag type="info" size="small">
                  {{ key }} = {{ value.length > 100 ? value.substring(0, 100) + "..." : value }}
                </el-tag>
              </div>
            </td>
          </tr>
        </table>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
        </div>
      </el-card>
      <el-tabs style="margin-top:20px" v-model="activeName" @tab-click="onTabChange">
        <el-tab-pane lazy label="Containers" name="Containers">
          <complex-table :data="form.status.containerStatuses">
            <el-table-column sortable :label="$t('commons.table.status')" min-width="30">
              <template v-slot:default="{row}">
                <el-button v-if="row.state.running" type="success" size="mini" plain round>
                  Running
                </el-button>
                <el-button v-if="!row.state.running" type="warning" size="mini" plain round>
                  Failed
                </el-button>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('business.pod.ready')" prop="ready" min-width="40">
              <template v-slot:default="{row}">
                <i class="el-icon-check" v-if="row.ready"/>
                <i class="el-icon-close" v-if="!row.ready"/>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="50"/>
            <el-table-column sortable :label="$t('business.pod.image')" min-width="170">
              <template v-slot:default="{row}">
                <div class="myTag">
                  <el-tag type="info" size="small">
                    {{ row.image }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('business.workload.restarts')" prop="restartCount" min-width="30"/>
            <el-table-column sortable :label="$t('commons.table.created_time')" min-width="70">
              <template v-slot:default="{row}">
                <span v-if="row.started">{{ row.state.running.startedAt | age }}</span>
                <span v-if="!row.started">-</span>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane lazy label="Conditions" name="Conditions">
          <complex-table :data="form.status.conditions">
            <el-table-column sortable label="Condition" prop="type"/>
            <el-table-column sortable :label="$t('commons.table.status')" prop="status"/>
            <el-table-column sortable :label="$t('commons.table.lastUpdateTime')" prop="lastUpdateTime">
              <template v-slot:default="{row}">
                {{ row.lastTransitionTime | age }}
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.message')" min-width="200">
              <template v-slot:default="{row}">
                <span v-if="row.message">[{{ row.reason }} ]: {{ row.message }}</span>
                <span v-if="!row.message">---</span>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane lazy label="Shell" name="Shell">
          <el-select v-model="selectedTerminalContainer" @change="onSelectedTerminalContainerChange" placeholder="选择容器">
            <el-option v-for="(item, index) in form.status.containerStatuses" :key="index" :label="item.name"
                       :value="item.name"/>
          </el-select>

          <div v-if="terminalOpened" style="margin-top: 5px">
            <iframe :src=getTerminalUrl style="width: 100%;min-height: 600px;border: 0"></iframe>
          </div>

        </el-tab-pane>
        <el-tab-pane lazy label="Logging" name="Logging">
          <el-select v-model="selectedLoggingContainer" placeholder="选择容器">
            <el-option v-for="(item, index) in form.status.containerStatuses" :key="index" :label="item.name"
                       :value="item.name"/>
          </el-select>
          <el-select v-model="tailLines" placeholder="显示条数">
            <el-option v-for="(item, index) in tailLinesOptions" :key="index" :label="item.label"
                       :value="item.value"/>
          </el-select>
          <el-switch
              v-model="follow"
              active-text="自动刷新">
          </el-switch>

          <div v-if="loggingOpened" style="margin-top: 5px">
            <iframe :src=getLoggingUrl style="width: 100%;min-height: 600px;border: 0"></iframe>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
    <div v-if="yamlShow">
      <yaml-editor :value="form" :read-only="true"></yaml-editor>
      <div class="bottom-button">
        <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
      </div>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getPodByName} from "@/api/pods"
import YamlEditor from "@/components/yaml-editor"

import ComplexTable from "@/components/complex-table"

export default {
  name: "PodDetail",
  components: {LayoutContent, ComplexTable, YamlEditor},
  data() {
    return {
      form: {
        metadata: {},
        spec: {
          nodeName: "",
        },
        status: {
          containerStatuses: [],
        },
      },
      pods: [],
      yamlShow: false,
      activeName: "Containers",
      loading: false,
      clusterName: "",

      terminalUrl: "",
      terminalOpened: false,
      selectedTerminalContainer: "",

      loggingOpened: false,
      loggingUrl: "",
      selectedLoggingContainer: "",
      follow: false,
      tailLines: 20,
      tailLinesOptions: [
        {label: "最后20行", value: 20,},
        {label: "最后100行", value: 100,},
        {label: "最后200行", value: 200,},
        {label: "最后500行", value: 500,},
      ]
    }
  },
  computed: {
    getTerminalUrl() {
      const namespace = this.form.metadata.namespace
      const podName = this.form.metadata.name
      const containerName = this.selectedTerminalContainer
      const clusterName = this.$route.query["cluster"]
      return `/terminal/app?cluster=${clusterName}&pod=${podName}&namespace=${namespace}&container=${containerName}`
    },
    getLoggingUrl() {
      const namespace = this.form.metadata.namespace
      const podName = this.form.metadata.name
      const containerName = this.selectedLoggingContainer
      const clusterName = this.$route.query["cluster"]
      const tailLines = this.tailLines
      const follow = this.follow

      let baseUrl = `/terminal/logging?cluster=${clusterName}&pod=${podName}&namespace=${namespace}`
      if (containerName) {
        baseUrl += `&&container=${containerName}`
      }
      if (tailLines) {
        baseUrl += `&&tailLines=${tailLines}`
      }
      if (follow) {
        baseUrl += `&&follow=${follow}`
      }
      return baseUrl
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getPodByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        this.loading = false
      })
    },
    onTabChange(tab) {
      if (this.terminalOpened) {
        this.terminalOpened = false
      }
      if (tab.label === "Shell") {
        this.selectedTerminalContainer = ""
        this.terminalContainer = ""
        let firstReadyContainer = ""
        for (const container of this.form.status.containerStatuses) {
          if (container.ready) {
            firstReadyContainer = container.name
            break
          }
        }
        if (!firstReadyContainer) {
          this.$message.error("暂无正在运行的容器")
        } else {
          this.selectedTerminalContainer = firstReadyContainer
          this.terminalContainer = firstReadyContainer
          this.terminalOpened = true
        }
      }
      if (tab.label === "Logging") {
        this.selectedLoggingContainer = ""
        this.loggingContainer = ""
        this.follow = false
        this.tailLines = 20
        this.loggingOpened = true
      }
    }
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
