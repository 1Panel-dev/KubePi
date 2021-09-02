<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'HPA'}" v-loading="loading">
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
                  <h3>{{ $t("business.common.config") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("business.workload.workload") }}</td>
                <td>{{ item.spec.scaleTargetRef.kind }}/{{ item.spec.scaleTargetRef.name }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.min_replicas") }}</td>
                <td>{{ item.spec.minReplicas }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.max_replicas") }}</td>
                <td>{{ item.spec.maxReplicas }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.configuration.current_replicas") }}</td>
                <td>{{ item.status.currentReplicas }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card">
            <el-tab-pane label="Metrics" v-if="item.spec.metrics.length > 0">
              <div v-for="(metric,index) in item.spec.metrics" v-bind:key="index"
                   style="border:1px solid #383c42;margin-top: 5px">
                <el-row :gutter="20">
                  <el-col :span="12">
                    <table style="width: 100%" class="myTable" v-if="metric.type">
                      <tr>
                        <th scope="col" width="30%" align="left">
                          <h3>{{ metric.type }} Metric</h3>
                        </th>
                        <th scope="col"></th>
                      </tr>
                      <tr>
                        <td>Target Name</td>
                        <td>{{ metric[metric.type.toLowerCase()].target.type }}</td>
                      </tr>
                      <tr>
                        <td>Value</td>
                        <td>{{ getTargetValue(metric[metric.type.toLowerCase()].target) }}</td>
                      </tr>
                      <tr v-if="metric[metric.type.toLowerCase()].name">
                        <td>Resource Name</td>
                        <td>{{ metric[metric.type.toLowerCase()].name }}</td>
                      </tr>
                      <tr v-if="metric[metric.type.toLowerCase()].metric">
                        <td>Name</td>
                        <td>{{ metric[metric.type.toLowerCase()].metric.name }}</td>
                      </tr>
                    </table>
                  </el-col>
                  <el-col :span="12">
                    <table style="width: 100%" class="myTable">
                      <tr>
                        <th scope="col" width="30%" align="left">
                          <h3>Current Metrics</h3>
                        </th>
                        <th scope="col"></th>
                      </tr>
                      <tr v-if="item.status.currentMetrics && item.status.currentMetrics[index].type !== '' ">
                        <td>Average Utilization</td>
                        <td v-if="metric.type">
                          {{ item.status.currentMetrics[index][metric.type.toLowerCase()].current.averageUtilization }}
                        </td>
                      </tr>
                      <tr v-if="item.status.currentMetrics && item.status.currentMetrics[index].type !== '' ">
                        <td>Average Value</td>
                        <td v-if="metric.type">
                          {{ item.status.currentMetrics[index][metric.type.toLowerCase()].current.averageValue }}
                        </td>
                      </tr>
                      <span v-if="item.status.currentMetrics && item.status.currentMetrics[index].type === '' ">No Current Metrics</span>
                    </table>
                  </el-col>
                </el-row>
              </div>
            </el-tab-pane>
            <el-tab-pane :label="$t('business.common.conditions')">
              <ko-detail-conditions :conditions="item.status.conditions"></ko-detail-conditions>
            </el-tab-pane>
            <el-tab-pane label="Events">
              <ko-detail-events :namespace="namespace" :cluster="cluster"
                                :selector="'involvedObject.kind=HorizontalPodAutoscaler,involvedObject.name='+name"></ko-detail-events>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>
    <el-row>
      <div v-if="yamlShow">
        <yaml-editor :value="item" :read-only="true"></yaml-editor>
        <div class="bottom-button">
          <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.back_detail") }}</el-button>
        </div>
      </div>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getHpa} from "@/api/hpa"
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailConditions from "@/components/detail/detail-conditions"
import KoDetailEvents from "@/components/detail/detail-events"

export default {
  name: "HPADetail",
  components: { KoDetailEvents, KoDetailConditions, KoDetailBasic, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String,
  },
  data () {
    return {
      loading: false,
      item: {
        metadata: {},
        spec: {
          scaleTargetRef: {},
          metrics: [{
            resource: {
              target: {}
            }
          }]
        },
        status: {
          currentMetrics: [{
            resource: {
              current: {}
            }
          }]
        }
      },
      cluster: "",
      yamlShow: false,
      events: []
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getHpa(this.cluster, this.namespace, this.name).then(res => {
        this.item = res
      }).finally(() => {
        this.loading = false
      })
    },
    getTargetValue (target) {
      let type = target.type
      if (type === "Utilization") {
        return target.averageUtilization
      } else {
        const resourceName = type.charAt(0).toLowerCase() + type.slice(1)
        return target[resourceName]
      }
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.yamlShow = this.$route.query.yamlShow === "true"
    this.getDetail()
  }
}
</script>

<style scoped>

</style>
