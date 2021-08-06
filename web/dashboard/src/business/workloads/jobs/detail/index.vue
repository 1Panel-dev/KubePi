<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Jobs'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-card>
        <ko-detail-basic :item="form" :yaml-show.sync="yamlShow"></ko-detail-basic>
      </el-card>
      <el-tabs style="margin-top:20px" v-model="activeName" type="border-card">
        <el-tab-pane label="Pods" name="Pods">
          <complex-table :data="jobs">
            <el-table-column sortable :label="$t('commons.table.status')" prop="status.phase" min-width="35">
              <template v-slot:default="{row}">
                <el-button v-if="row.status.phase ==='Succeeded'" type="success" size="mini" plain round>
                  {{row.status.phase}}
                </el-button>
                <el-button v-if="row.status.phase !=='Succeeded'" type="warning" size="mini" plain round>
                  {{row.status.phase}}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('commons.table.name')" prop="metadata.name" min-width="120">
              <template v-slot:default="{row}">
                <el-link @click="toResource('Pod', row.metadata.namespace, row.metadata.name)">{{ row.metadata.name }}</el-link>
              </template>
            </el-table-column>
            <el-table-column sortable :label="$t('business.cluster.nodes')" min-width="50" prop="spec.nodeName" />
            <el-table-column sortable :label="$t('business.pod.image')" min-width="150">
              <template v-slot:default="{row}">
                <div v-for="(item,index) in row.spec.containers" v-bind:key="index" class="myTag">
                  <el-tag type="info" size="small">
                    {{ item.image }}
                  </el-tag>
                </div>
              </template>
            </el-table-column>
          </complex-table>
        </el-tab-pane>
        <el-tab-pane label="Conditions" name="Conditions">
          <ko-detail-conditions :conditions="form.status.conditions"></ko-detail-conditions>
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
import { getJobByName } from "@/api/jobs"
import { listPodsWithNsSelector } from "@/api/pods"
import { listEventsWithNsSelector } from "@/api/events"
import YamlEditor from "@/components/yaml-editor"
import { mixin } from "@/utils/resourceRoutes"

import ComplexTable from "@/components/complex-table"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailConditions from "@/components/detail/detail-conditions"

export default {
  name: "JobDetail",
  components: { KoDetailConditions, KoDetailBasic, LayoutContent, YamlEditor, ComplexTable },
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
          labels: [],
          annotations: [],
        },
        status: [],
      },
      yamlShow: false,
      activeName: "Pods",
      loading: false,
      clusterName: "",
      jobs: [],
      events: [],
    }
  },
  watch: {
    yamlShow: function (newValue) {
      this.$router.push({
        name: "JobDetail",
        params: { name: this.name, namespace: this.namespace },
        query: { yamlShow: newValue },
      })
    },
  },
  methods: {
    getDetail() {
      this.loading = true
      this.events = []
      getJobByName(this.clusterName, this.namespace, this.name).then((res) => {
        this.form = res
        if (this.form.spec.template.metadata.labels) {
          let selectors = ""
          for (const key in this.form.spec.template.metadata.labels) {
            if (Object.prototype.hasOwnProperty.call(this.form.spec.template.metadata.labels, key)) {
              selectors += key + "=" + this.form.spec.template.metadata.labels[key] + ","
            }
          }
          selectors = selectors.length !== 0 ? selectors.substring(0, selectors.length - 1) : ""
          listPodsWithNsSelector(this.clusterName, this.namespace, selectors).then((res) => {
            this.jobs = res.items
          })
          let eventSelectors = "involvedObject.name=" + res.metadata.name + ",involvedObject.namespace=" + res.metadata.namespace + ",involvedObject.uid=" + res.metadata.uid
          listEventsWithNsSelector(this.clusterName, this.namespace, eventSelectors).then((res) => {
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
