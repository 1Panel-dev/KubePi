<template>
  <layout-content :header="$t('business.chart.app')" v-loading="loading">
    ss
    <div>
      <el-row :gutter="20">
        <el-col v-for="(d,index) in data" v-bind:key="index">
          <el-card>{{d}}</el-card>
        </el-col>
      </el-row>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {searchCharts} from "@/api/charts"

export default {
  name: "Charts",
  components: { LayoutContent },
  props: {},
  data () {
    return {
      cluster: "",
      data: [],
      loading: false,
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
    }
  },
  methods: {
    listAll () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchCharts(this.cluster,currentPage, pageSize,this.searchConfig.keywords).then(res => {
        this.data = res.data
        this.loading = false
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.listAll()
  }
}
</script>

<style scoped>

</style>
