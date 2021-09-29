<template>
  <layout-content :header="$t('business.chart.app')" v-loading="loading">
    <div>
      <el-row :gutter="20" v-infinite-scroll="load" :infinite-scroll-disabled="loading || moreLoading">
        <el-col v-for="(d,index) in shows" v-bind:key="index" :xs="8" :sm="8" :lg="6">
          <div>
            <el-card :body-style="{padding: '0px'}" class="d-card el-card">
              <el-row :gutter="24">
                <el-col :span="8">
                  <div style="margin-top: 32px;margin-left: 25px;margin-bottom: 15px">
                    <el-image :src="d.icon" fit="contain" style="height: 56px;width: 56px;position:relative">
                      <div slot="error" class="image-slot">
                        <img :src="require('@/assets/apps.svg')"/>
                      </div>
                    </el-image>
                  </div>
                </el-col>
                <el-col :span="16">
                  <div style="text-align: right">
                    <el-tag>{{ d.repo }}</el-tag>
                  </div>
                  <div class="card-content">
                    <span style="font-size: large;">{{ d.name }}</span>
                    <p class="chart-description">{{ d.description }}</p>
                  </div>
                </el-col>
              </el-row>
            </el-card>
          </div>
        </el-col>
      </el-row>
    </div>
    <div style="height: 20px;text-align: center">
      <i class="el-icon-loading" v-if="moreLoading"></i>
    </div>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listRepos, searchCharts} from "@/api/charts"

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
        pageSize: 400,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      repos: [],
      colors: [],
      shows: [],
      sliceNum: 24,
      currentSlice: 0,
      moreLoading: false
    }
  },
  methods: {
    listAll () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchCharts(this.cluster, currentPage, pageSize, this.searchConfig.keywords).then(res => {
        this.data = res.data.items
        this.shows = this.data.slice(0, 24)
        this.currentSlice = this.currentSlice + this.sliceNum
        this.paginationConfig.total = res.data.total
        this.loading = false
      })
      listRepos(this.cluster).then(res => {
        this.repos = res.data
      })
    },
    load () {
      if (this.loading || this.moreLoading) {
        return
      }
      if (this.shows.length < this.data.length) {
        this.moreLoading = true
        setTimeout(() => {
          this.moreLoading = false
          const addData = this.data.slice(this.currentSlice, this.currentSlice + this.sliceNum)
          this.shows = this.shows.concat(addData)
          this.currentSlice = this.currentSlice + addData.length
        }, 500)
      } else {
        this.paginationConfig.currentPage++
        this.moreLoading = true
        const { currentPage, pageSize } = this.paginationConfig
        searchCharts(this.cluster, currentPage, pageSize, this.searchConfig.keywords).then(res => {
          this.data = this.data.concat(res.data.items)
          this.moreLoading = false
          this.load()
        })
      }
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.listAll()
  }
}
</script>

<style scoped>
    .d-card {
        height: 120px;
        background-color: #1f2224;
        margin-top: 10px;
    }

    .d-card:hover {
        cursor: pointer;
        border: 1px solid #e92322;
        z-index: 1;
    }

    .card-content {
        text-align: left;
        float: left;
    }

    .chart-description {
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 3;
        overflow: hidden;
        color: #a1a9ae;
        font-size: small;
    }
</style>
