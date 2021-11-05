<template>
  <layout-content :header="$t('commons.form.detail')" :back-to="{name: 'Jobs'}" v-loading="loading">
    <div v-if="!yamlShow">
      <el-row :gutter="20" class="row-box">
        <el-col :span="12">
          <el-card class="el-card">
            <ko-detail-basic :item="form" :yaml-show.sync="yamlShow" />
          </el-card>
        </el-col>
        <el-col :span="12">
          <el-card class="el-card">
            <h3>{{ $t("business.common.conditions") }}</h3>
            <ko-detail-conditions :conditions="form.status.conditions" />
          </el-card>
        </el-col>
      </el-row>

      <el-row>
        <el-tabs style="margin-top:20px" v-model="activeName" type="border-card">
          <el-tab-pane label="Pods" name="Pods">
            <ko-detail-pods :cluster="clusterName" :namespace="namespace" :selector="selectors" />
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
import YamlEditor from "@/components/yaml-editor"
import KoDetailBasic from "@/components/detail/detail-basic"
import KoDetailPods from "@/components/detail/detail-pods"
import KoDetailConditions from "@/components/detail/detail-conditions"
import KoDetailEvents from "@/components/detail/detail-events"

export default {
  name: "JobDetail",
  components: { KoDetailPods, KoDetailConditions, KoDetailBasic, KoDetailEvents, LayoutContent, YamlEditor },
  props: {
    name: String,
    namespace: String,
  },
  data() {
    return {
      form: {
        metadata: {},
        status: {
          conditions: [],
        },
      },
      yamlShow: false,
      activeName: "Pods",
      loading: false,
      clusterName: "",
      selectors: "",
      eventSelectors: "",
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
      getWorkLoadByName(this.clusterName, "jobs", this.namespace, this.name).then((res) => {
        this.form = res
        if (this.form.spec.template.metadata.labels) {
          this.selectors = ""
          for (const key in this.form.spec.template.metadata.labels) {
            if (Object.prototype.hasOwnProperty.call(this.form.spec.template.metadata.labels, key)) {
              this.selectors += key + "=" + this.form.spec.template.metadata.labels[key] + ","
            }
          }
          this.selectors = this.selectors.length !== 0 ? this.selectors.substring(0, this.selectors.length - 1) : ""
          this.eventSelectors = "involvedObject.name=" + res.metadata.name + ",involvedObject.namespace=" + res.metadata.namespace + ",involvedObject.uid=" + res.metadata.uid
        }
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
