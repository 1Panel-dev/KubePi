<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'CronJobs'}" v-loading="loading">
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
            <td colspan="2">{{ form.metadata.namespace }}</td>
          </tr>
          <tr>
            <td>{{ $t("commons.table.created_time") }}</td>
            <td colspan="2">{{ form.metadata.creationTimestamp | datetimeFormat }}</td>
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
        <el-tab-pane label="Jobs" name="Jobs">
          <complex-table :data="jobs">
            <el-table-column sortable :label="$t('commons.table.status')" prop="status.succeeded" min-width="30">
              <template v-slot:default="{row}">
                <el-button v-if="row.status.succeeded ===1" type="success" size="mini" plain round>
                  Succeeded
                </el-button>
                <el-button v-if="row.status.succeeded ===2" type="warning" size="mini" plain round>
                  Failed
                </el-button>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.name')" prop="metadata.name" min-width="90" />
            <el-table-column sortable :label="$t('commons.table.status')" min-width="40">
              <template v-slot:default="{row}">
                {{ row.spec.completions }} / {{ row.spec.parallelism }}
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('business.workload.duration')" min-width="40">
              <template v-slot:default="{row}">
                {{ getDuration(row) }}S
              </template>
            </el-table-column>
            <el-table-column :label="$t('commons.table.created_time')" min-width="60" prop="metadata.creationTimestamp" fix>
              <template v-slot:default="{row}">
                {{ row.metadata.creationTimestamp | datetimeFormat }}
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane label="Events" name="Events">
          <complex-table :data="events">
            <el-table-column sortable :label="$t('commons.table.status')">
              <template v-slot:default="{row}">
                <span>{{ row.firstTimestamp | datetimeFormat }} - {{ row.lastTimestamp | datetimeFormat }}</span>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.message')" min-width=250>
              <template v-slot:default="{row}">
                [{{ row.reason }} ]: {{ row.message }}
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
import { getCronJobByName } from "@/api/cronjobs"
import { listJobsWithNsSelector } from "@/api/jobs"
import { listEventsWithNsSelector } from "@/api/events"
import YamlEditor from "@/components/yaml-editor"

import ComplexTable from "@/components/complex-table"

export default {
  name: "CronJobDetail",
  components: { LayoutContent, YamlEditor, ComplexTable },
  data() {
    return {
      form: {
        metadata: {
          name: "",
          namespace: "",
          creationTimestamp: "",
          labels: [],
          annotations: [],
        },
      },
      yamlShow: false,
      activeName: "Jobs",
      loading: false,
      clusterName: "",
      jobs: [],
      events: [],
    }
  },
  methods: {
    getDetail() {
      this.loading = true
      this.events = []
      getCronJobByName(this.clusterName, this.$route.params.namespace, this.$route.params.name).then((res) => {
        this.form = res
        if (this.form.spec.jobTemplate.spec.template.metadata.labels) {
          let selectors = ""
          for (const key in this.form.spec.jobTemplate.spec.template.metadata.labels) {
            if (Object.prototype.hasOwnProperty.call(this.form.spec.jobTemplate.spec.template.metadata.labels, key)) {
              selectors += key + "=" + this.form.spec.jobTemplate.spec.template.metadata.labels[key] + ","
            }
          }
          selectors = selectors.length !== 0 ? selectors.substring(0, selectors.length - 1) : ""
          listJobsWithNsSelector(this.clusterName, this.$route.params.namespace, selectors).then((res) => {
            this.jobs = res.items
          })
          let eventSelectors = "involvedObject.name=" + res.metadata.name + ",involvedObject.namespace=" + res.metadata.namespace + ",involvedObject.uid=" + res.metadata.uid
          listEventsWithNsSelector(this.clusterName, this.$route.params.namespace, eventSelectors).then((res) => {
            for (const e of res.items) {
              this.events.unshift(e)
            }
          })
          this.events.sort(function (a, b) {
            return Date.parse(b.firstTimestamp.replace(/-/g, "/")) - Date.parse(a.firstTimestamp.replace(/-/g, "/"))
          })
        }
        this.loading = false
      })
    },
    getDuration(row) {
      let startTime = new Date(row.status.startTime)
      let endTime = new Date(row.status.completionTime)
      return Math.floor((endTime - startTime) / 1000)
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