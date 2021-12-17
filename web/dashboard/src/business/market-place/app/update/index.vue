<template>
  <layout-content :header="$t('business.chart.upgrade')" v-loading="loading" :back-to="{ name: 'Apps' }">
    <el-form ref="form" label-position='left' label-width="220px" :model="form" :rules="rules">
      <fu-steps ref="steps" footerAlign="flex" finish-status="success" @finish="onSubmit" @cancel="onCancel"
                showCancel>
        <fu-step id="version" :title="'Version'">
          <div class="example">
            <el-scrollbar style="height:100%">
              <el-row>
                <el-col :span="12">
                  <el-form-item :label="$t('business.chart.version')" prop="chartVersion">
                    <el-select v-model="form.chartVersion" @change="getDetailByVersion(form.chartVersion)">
                      <el-option v-for="(vs,index) in versions"
                                 :label="vs.version" :key="index" :value="vs.version">
                        <span style="float: left">{{ vs.version }}</span>
                        <span style="float: right; color: #8492a6; font-size: 13px">AppVersion:{{ vs.appVersion }}</span>
                        <span style="float: right;" v-if="vs.version === oldChart.metadata.version">(current)</span>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-col>
              </el-row>
            </el-scrollbar>
          </div>
        </fu-step>
        <fu-step id="values" :title="'Values'">
          <div class="example">
            <el-scrollbar style="height:100%">
              <yaml-editor :value="form.values" ref="yaml_editor"></yaml-editor>
            </el-scrollbar>
          </div>
        </fu-step>
      </fu-steps>
    </el-form>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getApp, getChartByVersion, getChartUpdate, updateChart} from "@/api/charts"
import YamlEditor from "@/components/yaml-editor"

export default {
  name: "AppUpgrade",
  components: { LayoutContent, YamlEditor },
  props: {
    name: String
  },
  data () {
    return {
      loading: false,
      form: {},
      rules: {},
      oldValues: null,
      cluster: "",
      versions: [],
      chartMap: null,
      repo: "",
      oldChart: "",
      namespace: ""
    }
  },
  methods: {
    onSubmit () {
      const installForm = this.form
      installForm.cluster = this.cluster
      installForm.chartName = this.oldChart.metadata.name
      installForm.name = this.name
      installForm.repo = this.repo
      installForm.namespace = this.namespace
      installForm.values = this.$refs.yaml_editor.getValue()
      this.loading = true
      updateChart(this.cluster, this.name, installForm).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.operation_success"),
        })
        this.$router.push({ name: "Apps" })
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.$router.push({ name: "Apps" })
    },
    getDetail () {
      this.loading = true
      getApp(this.cluster, this.name).then(res => {
        this.oldValues = res.data.config
        this.oldChart = res.data.chart
        this.form.values = res.data.config
        this.form.chartVersion = res.data.chart.metadata.version
        this.chartMap.set(res.data.chart.metadata.version, {values:res.data.config} )
        this.namespace = res.data.namespace
        getChartUpdate(this.cluster, res.data.chart.metadata.name, this.name).then(res => {
          this.versions = res.data.versions
          this.repo = res.data.repo
        }).finally(() => {
          this.loading = false
        })
      }).finally(() => {
        this.loading = false
      })

    },
    getDetailByVersion (version) {
      this.$forceUpdate()
      if (this.chartMap.has(version)) {
        this.form.values = this.chartMap.get(version).values
        return
      }
      this.loading = true
      getChartByVersion(this.cluster, this.repo, this.oldChart.metadata.name, version).then(res => {
        this.form.values = res.data.values
        this.chartMap.set(res.data.metadata.version, res.data)
        this.loading = false
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.chartMap = new Map()
    this.getDetail()
  }
}
</script>

<style scoped>
    .example {
        margin: 0 5%;
        height: 600px;
    }
</style>
