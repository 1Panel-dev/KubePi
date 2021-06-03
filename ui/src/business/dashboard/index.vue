<template>
  <layout-content :header="$t('business.dashboard.dashboard')">
    <el-row :gutter="24">
      <el-col :span="24">
        <el-card shadow="always" style="background-color: #243441;height: 120px">
          <el-row :gutter="24">
            <el-col :span="6">
              <span class="title">{{ $t("commons.table.name") }}</span>
              <div style="text-align: center">
                <h1>测试</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("business.cluster.version") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>v1.18.0</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("business.cluster.nodes") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>2</h1>
              </div>
            </el-col>
            <el-col :span="6">
              <span class="title">{{ $t("commons.table.created_time") }}</span>
              <div class="line"></div>
              <div style="text-align: center">
                <h1>16 days ago</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <br>
    <el-row :gutter="24" class="resources">
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}" class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}"  class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}"  class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}" class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}"  class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}"  class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card :body-style="{padding: '0px'}" class="d-card">
          <el-row :gutter="24">
            <el-col :span="10">
              测试
            </el-col>
            <el-col :span="14">
              <div class="card-content">
                <span>Namespaces</span>
                <h1>20</h1>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <h2>Events</h2>
      <complex-table class="d-table">

      </complex-table>
    </el-row>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"

export default {
  name: "Dashboard",
  components: { ComplexTable, LayoutContent }
}
</script>

<style scoped>
    .line {
        float: left;
        margin-top: 10px;
        width: 1px;
        height: 80px;
        background: #3884c5;
    }

    .title {
        margin-left: 10px;
        color: #a1a9ae;
    }

    .d-card {
        height: 90px;
        background-color: #1d3e4d;
        margin-top: 10px;
    }
    .card-content {
        margin-top: 10px;
        margin-right: 10px;
        text-align: center;
        float: right;
    }

    .card-content > span:first-child {
        color: #a1a9ae;
    }

    .resources {

    }

    .d-table {
        margin-top: 10px;
    }
</style>