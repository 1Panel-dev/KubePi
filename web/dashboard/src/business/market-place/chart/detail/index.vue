<template>
  <layout-content :header="$t('business.chart.app')" v-loading="loading">
    <el-row type="flex" :gutter="20">
      <el-col :span="16">
        <div>
          <el-image :src="current.metadata.icon" fit="contain" style="height: 56px;width: 56px;position:relative">
            <div slot="error" class="image-slot">
              <img :src="require('@/assets/apps.svg')"/>
            </div>
          </el-image>
        </div>
      </el-col>
    </el-row>
    <el-row type="flex" :gutter="20">
      <el-col :span="20">
        <div>
          <div v-html="compiledMarkdown">
          </div>
        </div>
      </el-col>
      <el-col :span="4">
        <div>
          <h3>Chart Versions</h3>
          <table>
            <tr>
              <th width="33%"></th>
              <th width="33%"></th>
              <th width="33%"></th>
            </tr>
            <tr v-for="(version,index) in versions" :key="index">
              <td colspan="2">
                <el-link v-if="version.version === current.metadata.version">{{ version.version }}</el-link>
                <span v-else>{{ version.version }}</span>
              </td>
              <td>{{ version.date | dateFormat }}</td>
            </tr>
          </table>
        </div>
        <div>
          <h3>Application Version</h3>
          <span>{{current.metadata.appVersion}}</span>
        </div>
        <div>
          <h3>Maintainers</h3>
          <table>
            <tr v-for="(maintainer,index) in current.metadata.maintainers" :key="index">
              <el-link :href="'mailto:'+maintainer.email">{{maintainer.name}}</el-link>
            </tr>
          </table>
        </div>
        <div>
          <h3>Related</h3>
          <table>
            <tr v-for="(source,index) in current.metadata.sources" :key="index">
              <el-link :href="source">{{source}}</el-link>
            </tr>
          </table>
        </div>
      </el-col>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getChart} from "@/api/charts"
import marked from "marked"
import DOMPurify from "dompurify"

export default {
  name: "ChartDetail",
  components: { LayoutContent },
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
    }
  },
  methods: {
    getDetail () {
      this.loading = true
      getChart(this.cluster, this.repo, this.name).then(res => {
        this.loading = false
        this.current = res.data.chart
        this.versions = res.data.versions
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.getDetail()

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

</style>
