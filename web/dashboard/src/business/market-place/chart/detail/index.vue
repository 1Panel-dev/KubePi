<template>
  <layout-content :header="$t('business.chart.app')" v-loading="loading" :back-to="{ name: 'Charts' }">
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
          <el-button style="margin-right: 150px" type="primary">安装</el-button>
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
                <el-link v-if="version.version === current.metadata.version">{{ version.version }}</el-link>
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
    .title {
        font-size: 20px;
    }

    .detail {
        margin-top: 20px;
    }
</style>
