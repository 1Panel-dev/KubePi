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
      <el-tabs style="margin-top:20px" v-model="activeName">
        <el-tab-pane label="Containers" name="Containers">
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
                <i class="el-icon-check" v-if="row.ready" />
                <i class="el-icon-close" v-if="!row.ready" />
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.name')" prop="name" min-width="50" />
            <el-table-column sortable :label="$t('business.pod.image')" min-width="170">
              <template v-slot:default="{row}">
                <div class="myTag">
                  <el-tag type="info" size="small">
                    {{ row.image }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('business.workload.restarts')" prop="restartCount" min-width="30" />
            <el-table-column sortable :label="$t('commons.table.created_time')" min-width="70">
              <template v-slot:default="{row}">
                <span v-if="row.started">{{ row.state.running.startedAt | age }}</span>
                <span v-if="!row.started">-</span>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane label="Conditions" name="Conditions">
          <complex-table :data="form.status.conditions">
            <el-table-column sortable label="Condition" prop="type" />
            <el-table-column sortable :label="$t('commons.table.status')" prop="status" />
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
import { getPodByName } from "@/api/pods"
// import { listPodsWithNsSelector } from "@/api/pods"
import YamlEditor from "@/components/yaml-editor"

import ComplexTable from "@/components/complex-table"

export default {
  name: "PodDetail",
  components: { LayoutContent, ComplexTable, YamlEditor },
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
      yamlShow: false,
      activeName: "Containers",
      loading: false,
      clusterName: "",
      pods: [],
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      getPodByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        // if (this.form.spec.selector.matchLabels) {
        //   let selectors = ""
        //   for (const key in this.form.spec.selector.matchLabels) {
        //     if (Object.prototype.hasOwnProperty.call(this.form.spec.selector.matchLabels, key)) {
        //       selectors += key + "=" + this.form.spec.selector.matchLabels[key] + ","
        //     }
        //   }
        //   selectors = selectors.length !== 0 ? selectors.substring(0, selectors.length - 1) : ""
        //   listPodsWithNsSelector(this.clusterName, this.$route.params.namespace, selectors).then((res) => {
        //     this.pods = res.items
        //   })
        // }
        this.loading = false
      })
    },
  },
  created() {
    this.clusterName = this.$route.query.cluster
    this.getDetail()
  },
}
</script>

<style scoped>
</style>
