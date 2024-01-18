<template>
  <layout-content :header="$t('business.chart.app')" v-loading="loading"  :back-to="{ name: 'Charts' }">
    <div v-if="!installStep">
      <el-row type="flex" :gutter="20">
        <el-col :span="16">
          <div>
            <el-image :src="current.metadata.icon" fit="contain" style="height: 65px;width: 65px;position:relative">
              <div slot="error" class="image-slot">
                <img :src="require('@/assets/apps.svg')"/>
              </div>
            </el-image>
          </div>
        </el-col>
        <el-col :span="8">
          <div style="text-align: right">
            <el-button  v-has-permissions="{scope:'cluster',apiGroup:'kubepi.org',resource:'appmarkets',verb:'create'}" style="margin-right: 150px" type="primary" @click="install">{{ $t("business.chart.install") }}
            </el-button>
          </div>
        </el-col>
      </el-row>
      <el-row type="flex" :gutter="20">
        <el-col :span="18">
          <div>
            <div v-html="compiledMarkdown">
            </div>
          </div>
        </el-col>
        <el-col :span="6">
          <div class="detail">
            <span class="title">Chart Versions</span>
            <br>
            <table>
              <tr>
                <th width="33%"></th>
                <th width="33%"></th>
                <th width="33%"></th>
              </tr>
              <tr v-for="(version,index) in versions" :key="index">
                <td colspan="2">
                  <el-link v-if="version.version !== current.metadata.version"
                           @click="getDetailByVersion(version.version)">{{ version.version }}
                  </el-link>
                  <span v-else>{{ version.version }}</span>
                </td>
                <td>{{ version.date | dateFormat }}</td>
              </tr>
            </table>
          </div>
          <div class="detail">
            <span class="title">Application Version</span>
            <br>
            <table>
              <tr>
                <td><span>{{ current.metadata.appVersion }}</span></td>
              </tr>
            </table>
          </div>
          <div class="detail">
            <span class="title">Maintainers</span>
            <br>
            <table>
              <tr v-for="(maintainer,index) in current.metadata.maintainers" :key="index">
                <el-link :href="'mailto:'+maintainer.email">{{ maintainer.name }}</el-link>
              </tr>
            </table>
          </div>
          <div class="detail">
            <span class="title">Related</span>
            <br>
            <table>
              <tr v-for="(source,index) in current.metadata.sources" :key="index">
                <el-link :href="source">{{ source }}</el-link>
              </tr>
            </table>
          </div>
          <div class="detail">
            <span class="title">Keywords</span>
            <br>
            <table>
              <tr>
                <td>
                <span v-for="(keyword,index) in current.metadata.keywords" :key="index">
                  <span>{{ keyword }}&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;&emsp;</span>
                </span>
                </td>
              </tr>
            </table>
          </div>
        </el-col>
      </el-row>
    </div>
    <div v-if="installStep" style="display: block">
      <el-form ref="form" label-position='left' label-width="220px" :model="form" :rules="rules">
        <fu-steps ref="steps" footerAlign="flex" finish-status="success" @finish="onSubmit" @cancel="onCancel" :beforeLeave="beforeLeave"
                  :isLoading="loading"
                  showCancel>
          <fu-step id="metadata" :title="'Metadata'">
            <div class="example">
              <el-scrollbar style="height:100%">
                <el-row>
                  <el-col :span="12">
                    <el-form-item :label="$t('business.chart.name')" prop="name">
                      <el-input v-model="form.name" clearable></el-input>
                    </el-form-item>
                    <el-form-item :label="$t('business.chart.namespace')" prop="namespace" required>
                      <el-select v-model="form.namespace">
                        <el-option v-for="(ns,index) in namespaces"
                                   :label="ns" :key="index" :value="ns"></el-option>
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
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getChart, getChartByVersion, installChart} from "@/api/charts"
import { marked } from 'marked'
import DOMPurify from "dompurify"
import {getNamespaces} from "@/api/auth"
import YamlEditor from "@/components/yaml-editor"
import Rule from "@/utils/rules"

export default {
  name: "ChartDetail",
  components: { YamlEditor, LayoutContent },
  props: {
    name: String,
    repo: String,
  },
  data () {
    return {
      loading: false,
      current: {
        metadata: {}
      },
      versions: [],
      cluster: "",
      chartMap: null,
      installStep: false,
      form: {},
      rules: {
        name: [Rule.RequiredRule, Rule.CommonNameRule],
        namespace: [Rule.RequiredRule],
      },
      namespaces: [],
      validate: false
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getChart(this.cluster, this.repo, this.name).then(res => {
        this.loading = false
        this.current = res.data.chart
        this.versions = res.data.versions
        this.chartMap.set(this.current.metadata.version, res.data.chart)
      })
    },
    getDetailByVersion (version) {
      if (this.chartMap.has(version)) {
        this.current = this.chartMap.get(version)
        return
      }
      this.loading = true
      getChartByVersion(this.cluster, this.repo, this.name, version).then(res => {
        this.current = res.data
        this.chartMap.set(this.current.metadata.version, res.data)
        this.loading = false
      })
    },
    install () {
      this.installStep = true
      this.form = {
        name: "",
        namespace: ""
      }
      this.form.values = this.current.values
    },
    onSubmit () {
      this.loading = true
      const installForm = this.form
      installForm.cluster = this.cluster
      installForm.chartVersion = this.current.metadata.version
      installForm.chartName = this.name
      installForm.repo = this.repo
      installForm.values = this.$refs.yaml_editor.getValue()
      installChart(this.cluster,installForm).then(() => {
        this.loading = false
        this.installStep = false
        this.$message({
          type: "success",
          message: this.$t("commons.msg.operation_success"),
        })
        this.$router.push({name: "Apps"})
      }).finally(() => {
        this.loading = false
      })
    },
    onCancel () {
      this.installStep = false
    },
    beforeLeave(step, isNext){
      if (!isNext) {
        return
      }
      this.checkForm()
      if (!this.validate) {
        return  false
      }
      return this.validate
    },
    checkForm() {
      this.validate = false
      this.$refs["form"].validate((valid) => {
          this.validate = valid
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.chartMap = new Map()
    this.getDetail()
    getNamespaces(this.cluster).then((res) => {
      this.namespaces = res.data
    })
  },
  computed: {
    compiledMarkdown () {
      if (this.current.readme === undefined || this.current.readme === "") {
        return
      }
      const html = marked(this.current.readme, {
        gfm: true,
        breaks: true,
        tables: true,
        smartLists: true,
        smartypants: true,
        mangle: true,
      })
      return DOMPurify.sanitize(html)
    }
  }
}
</script>

<style scoped>
    .title {
        font-size: 20px;
    }

    .detail {
        margin-top: 20px;
    }

    .example {
        margin: 0 5%;
        height: 600px;
    }

</style>
