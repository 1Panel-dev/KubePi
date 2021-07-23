<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'HPA'}" v-loading="loading">
    <el-row :gutter="20">
      <div v-if="!yamlShow">
        <el-col :span="15">
          <el-card>
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" width="30%" align="left">
                  <h3>{{ $t("business.common.basic") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>{{ $t("commons.table.name") }}</td>
                <td>{{ item.metadata.name }}</td>
              </tr>
              <tr>
                <td>{{ $t("business.namespace.namespace") }}</td>
                <td>{{ item.metadata.namespace }}</td>
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
                  <div v-for="(value,key,index) in item.metadata.annotations" v-bind:key="index" class="myTag">
                    <el-tag type="info" size="small" v-if="value.length < 50">
                      {{ key }} = {{ value }}
                    </el-tag>
                    <el-tooltip v-if="value.length > 50" :content="value" placement="top">
                      <el-tag type="info" size="small" v-if="value.length >= 50">
                        {{ key }} = {{ value.substring(0, 50) + "..." }}
                      </el-tag>
                    </el-tooltip>
                  </div>
                </td>
              </tr>
              <tr>
                <td>{{ $t("commons.table.created_time") }}</td>
                <td>{{ item.metadata.creationTimestamp | age }}</td>
              </tr>
            </table>
            <div class="bottom-button">
              <el-button @click="yamlShow=!yamlShow">{{ $t("commons.button.view_yaml") }}</el-button>
            </div>
          </el-card>
        </el-col>
        <el-col :span="9">
          <el-card>
            <table style="width: 100%" class="myTable">
              <tr>
                <th scope="col" width="30%" align="left">
                  <h3>{{ $t("business.common.config") }}</h3>
                </th>
                <th scope="col"></th>
              </tr>
              <tr>
                <td>Workload</td>
                <td>{{ item.spec.scaleTargetRef.kind }}/{{ item.spec.scaleTargetRef.name }}</td>
              </tr>
              <tr>
                <td>Min Replicas</td>
                <td>{{ item.spec.minReplicas }}</td>
              </tr>
              <tr>
                <td>Max Replicas</td>
                <td>{{ item.spec.maxReplicas }}</td>
              </tr>
              <tr>
                <td>Current Replicas</td>
                <td>{{ item.status.currentReplicas }}</td>
              </tr>
            </table>
          </el-card>
        </el-col>
        <el-col :span="24">
          <br>
          <el-tabs type="border-card">
            <el-tab-pane label="Metrics" v-if="item.spec.metrics.length > 0">
              <div v-for="(metric,index) in item.spec.metrics" v-bind:key="index"  style="border:1px solid #383c42;margin-top: 5px">
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
                      <tr v-if="item.status.currentMetrics[index].type !== '' ">
                        <td>Average Utilization</td>
                        <td v-if="metric.type">
                          {{ item.status.currentMetrics[index][metric.type.toLowerCase()].current.averageUtilization }}
                        </td>
                      </tr>
                      <tr v-if="item.status.currentMetrics[index].type !== '' ">
                        <td>Average Value</td>
                        <td v-if="metric.type">
                          {{ item.status.currentMetrics[index][metric.type.toLowerCase()].current.averageValue }}
                        </td>
                      </tr>
                      <span v-if="item.status.currentMetrics[index].type === '' ">No Current Metrics</span>
                    </table>
                  </el-col>
                </el-row>
              </div>
            </el-tab-pane>
            <el-tab-pane label="Conditions">
              <complex-table :data="item.status.conditions">
                <el-table-column label="Condition" prop="type">
                  <template v-slot:default="{row}">
                    {{ row.type }}
                  </template>
                </el-table-column>
                <el-table-column label="Status" prop="status">
                  <template v-slot:default="{row}">
                    {{ row.status }}
                  </template>
                </el-table-column>
                <el-table-column label="Reason" prop="reason">
                  <template v-slot:default="{row}">
                    {{ row.reason }}
                  </template>
                </el-table-column>
                <el-table-column label="Message" prop="message" min-width="200px">
                  <template v-slot:default="{row}">
                    {{ row.message }}
                  </template>
                </el-table-column>
              </complex-table>
            </el-tab-pane>
            <el-tab-pane label="Events">
              <complex-table :data="events">
                <el-table-column :label="$t('business.event.reason')" prop="reason" fix max-width="50px">
                  <template v-slot:default="{row}">
                    {{ row.reason }}
                  </template>
                </el-table-column>
                <el-table-column :label="$t('business.event.type')" prop="type" fix max-width="50px">
                  <template v-slot:default="{row}">
                    <el-tag v-if="row.type ==='Normal'">{{ row.type }}</el-tag>
                    <el-tag v-else type="danger">{{ row.type }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" fix max-width="50px">
                  <template v-slot:default="{row}">
                    {{ row.metadata.namespace }}
                  </template>
                </el-table-column>
                <el-table-column :label="$t('business.event.message')" prop="resource" fix min-width="200px"
                                 show-overflow-tooltip>
                  <template v-slot:default="{row}">
                    {{ row.message }}
                  </template>
                </el-table-column>
              </complex-table>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </div>
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
import ComplexTable from "@/components/complex-table"
import {listEventsWithNsSelector} from "@/api/events"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "HPADetail",
  components: { ComplexTable, LayoutContent, YamlEditor },
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
      const selectors = "involvedObject.kind=HorizontalPodAutoscaler,involvedObject.name=" + this.name
      listEventsWithNsSelector(this.cluster, this.namespace, selectors).then(res => {
        this.events = res.items
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
