<template>
  <layout-content :header="`${instance}${$t('business.monitor.metrics.name')}`" :back-to="{ name: 'Metrics' }">
    <el-row>
      <el-col  style="top: auto;" :span="22">
        <div class="promql-input">
          <el-autocomplete v-model="promql" placeholder="Please input PromQL" class="promql-input" size="large" :fetch-suggestions="filteredPromql" :trigger-on-focus="true" @keyup.enter.native="onExecute()">
            <template #prepend>
              <el-button icon="iconfont iconsearch" />
            </template>
            <template #append>
              <el-button type="primary" icon="iconfont iconwangluo" @click="onExporterBox()" />
            </template>
          </el-autocomplete>
        </div>
      </el-col>
      <el-button type="primary" :span="2" size="large" class="square-button" @click="onExecute()" >{{ $t('business.monitor.metrics.execute') }}</el-button>
    </el-row>

    <!-- 只保留 Table 标签，去掉按钮 -->
    <div class="table-container">
      <h4>Table</h4>
      <el-table :data="promql_query_data"  v-if="promql_query_data.length > 0">
        <el-table-column prop="metrics" width="auto"/>
        <el-table-column prop="value" align="right" width="100px"/>
      </el-table>
      <p v-if="promql_query_data.length === 0">No data queried yet</p>
    </div>

    <!-- Dialog Component -->
    <el-dialog title="Metrics Explorer" :visible.sync="dialogVisible" :close-on-click-modal="true" width="60%">
      <el-input
        v-model="explorer_search"
        placeholder="Search"
        size="large"
        style="margin-bottom: 20px"
      />
      <div class="scrollable-content">
        <div v-for="(metric, index) in filteredMetrics" :key="index" class="metric-item">
          {{ metric }}
        </div>
      </div>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {getExplorer, testConnectMetrics, queryMetrics} from "@/api/monitor"

export default {
  components: { LayoutContent },
  name: "MetricsDetail",
  props: {
    repo: String,
  },
  data () {
    return {
      promql: "",
      promql_list: [],
      dialogVisible: false,
      loading: false,
      explorer_list: [],
      explorer_search: "",
      promql_query_data: [],
    }
  },
  methods: {
    search () {
      getExplorer(this.$route.params.instance).then(res => {
        this.explorer_list = res.data.data
        this.promql_list = res.data.data.map(item => ({ value: item }));
      });
    },
    connectTest () {
      this.loading = true
      testConnectMetrics(this.$route.params.instance).then(() => {
        this.search()
      }).catch(() => {
        this.$router.push({ name: "Metrics" })
      })
    },
    onExporterBox () {
      // getExplorer(this.$route.params.instance).then(res => {
      //   this.explorer_list = res.data.data
      // })
      this.dialogVisible = true;
    },
    filteredPromql(queryString, callback) {
      if (queryString) {
        const results = this.promql_list.filter(metric =>
          metric.value.toLowerCase().includes(queryString.toLowerCase())
        );
        callback(results); // 回传过滤结果
      } else {
        callback(this.promql_list); // 如果没有输入内容，显示全部建议
      }
    },
    onExecute () {
      const nowTime = Math.floor(Date.now() / 1000);
      queryMetrics(this.$route.params.instance,this.promql,nowTime).then(res => {
        this.promql_query_data = res.data;
        console.log("1233")
        console.log("promql_query_data",this.promql_query_data)
      }).catch(() => {
        //this.$router.push({ name: "Metrics" })
      })
    },
  },
  created () {
    this.connectTest()
    //this.search()
  },
  computed: {
    instance() {
      return this.$route.params.instance+` `+this.$t("business.monitor.metrics.instance");
    },
    filteredMetrics() {
      return this.explorer_list.filter(metric =>
        metric.toLowerCase().includes(this.explorer_search.toLowerCase())
      );
    },
  },
}
</script>

<style scoped>
  .square-button {
    border-radius: 0 !important; /* 直角边框 */
    margin-left: -1px;
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
  }
  .scrollable-content {
    max-height: 400px;        /* 限制内容的最大高度 */
    overflow-y: auto;         /* 垂直滚动条 */
    padding-right: 10px;      /* 防止滚动条遮住内容 */
  }
  .metric-item {
    padding: 10px 0;          /* 每行内容的上下间距 */
  }
  .promql-input {
    flex: 1;
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
    width: 100%;
  }
  .metrics-panel {
    margin-top: 20px;
  }
  .table-container {
    border: 2px solid #444444; /* 灰色边框 */
    border-radius: 0 !important; /* 直角边框 */
    padding: 10px; /* 内边距 */
    margin-top: 10px;
  }
  /* 防止表格内容换行 */
  .el-table-column {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>