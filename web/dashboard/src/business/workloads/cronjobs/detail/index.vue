<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'CronJobs'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row class="row-box">
        <el-card class="el-card">
          <ko-detail-basic :item="form" :yaml-show.sync="yamlShow"></ko-detail-basic>
        </el-card>
      </el-row>
      <el-row>
        <el-tabs style="margin-top:20px" v-model="activeName" type="border-card">
          <el-tab-pane label="Jobs" name="Jobs">
            <complex-table :data="jobs">
              <el-table-column :label="$t('commons.table.status')" prop="status.succeeded" min-width="30">
                <template v-slot:default="{row}">
                  <el-button v-if="row.status.succeeded" type="success" size="mini" plain round>
                    Succeeded
                  </el-button>
                  <el-button v-if="row.status.failed" type="warning" size="mini" plain round>
                    Failed
                  </el-button>
                  <el-button v-if="row.status.active" type="success" size="mini" plain round>
                    Active
                  </el-button>
                </template>
              </el-table-column>
              <el-table-column :label="$t('commons.table.name')" prop="metadata.name" min-width="90" show-overflow-tooltip>
                <template v-slot:default="{row}">
                  <span class="span-link" @click="toResource('Job', row.metadata.namespace, row.metadata.name)">{{
                      row.metadata.name
                    }}
                  </span>
                </template>
              </el-table-column>

              <el-table-column :label="$t('business.workload.completions')" min-width="30">
                <template v-slot:default="{row}">
                  {{ row.spec.completions }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('business.workload.completions_status')" min-width="60">
                <template v-slot:default="{row}">
                  <el-tag style="margin-left: 5px" v-if="row.status.active" type="info">active: *{{row.status.active}}</el-tag>
                  <el-tag style="margin-left: 5px" v-if="row.status.succeeded" type="success">succeeded: {{row.status.succeeded}}</el-tag>
                  <el-tag style="margin-left: 5px" v-if="row.status.failed" type="danger">failed: {{row.status.failed}}</el-tag>
                </template>
              </el-table-column>
              <el-table-column :label="$t('business.workload.duration')" min-width="30">
                <template v-slot:default="{row}">
                  {{ getDuration(row) }}
                </template>
              </el-table-column>
              <el-table-column :label="$t('commons.table.created_time')" min-width="60" prop="metadata.creationTimestamp" fix>
                <template v-slot:default="{row}">
                  {{ row.metadata.creationTimestamp | age }}
                </template>
              </el-table-column>
            </complex-table>
          </el-tab-pane>
          <el-tab-pane :label="$t('business.event.event')" name="Events">
            <ko-detail-events :cluster="clusterName" :namespace="namespace" :selector="eventSelectors" />
          </el-tab-pane>
        </el-tabs>
      </el-row>
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
import { getWorkLoadByName } from "@/api/workloads"
import { listJobsWithNsSelector } from "@/api/jobs"
import YamlEditor from "@/components/yaml-editor"
import ComplexTable from "@/components/complex-table"
import { mixin } from "@/utils/resourceRoutes"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailEvents from "@/components/detail/detail-events"

export default {
  name: "CronJobDetail",
  components: { KoDetailBasic, KoDetailEvents, LayoutContent, YamlEditor, ComplexTable },
  mixins: [mixin],
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      form: {
        metadata: {
          name: "",
          namespace: "",
          creationTimestamp: "",
          labels: {},
          annotations: {},
        },
      },
      yamlShow: false,
      activeName: "Jobs",
      loading: false,
      clusterName: "",
      jobs: [],
      eventSelectors: "",
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "CronJobDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue },
      })
    },
  },
  methods: {
    getDetail() {
      this.loading = true
      this.events = []
      getWorkLoadByName(this.clusterName, "cronjobs", this.namespace, this.name).then((res) => {
        this.form = res
        this.loading = false
      })
      listJobsWithNsSelector(this.clusterName, this.namespace).then((res) => {
        this.jobs = []
        for (const job of res.items) {
          if (job.metadata.name.indexOf("-") !== -1) {
            if (job.metadata.name.substring(0, job.metadata.name.lastIndexOf("-")) === this.name) {
              this.jobs.push(job)
            }
          }
        }
        this.loading = false
      })
    },
    getDuration(row) {
      let startTime = new Date(row.status.startTime)
      if (!row.status.completionTime) {
        return "/"
      } else {
        let endTime = new Date(row.status.completionTime)
        let t = Math.floor((endTime - startTime) / 1000)
        if (t % 60 !== 0) {
          return (t % 60) + " mins"
        }
        if (t % 3600 !== 0) {
          return (t % 60) + " hours ago"
        }
        return Math.floor((endTime - startTime) / 1000) + "S"
      }
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
