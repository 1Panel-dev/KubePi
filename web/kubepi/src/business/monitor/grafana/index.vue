<template>
    <layout-content class="full-layout-content">
      <el-col :span="2"><br/></el-col>
      <el-col :span="20">
        <div v-if="iframeVisible" class="iframe-container">
          <!-- 当 iframeVisible 为 true 时，嵌入 iframe -->
          <iframe :src="form.address+'/dashboards/f/nErXDvCkzzkubepi?kiosk=tv'" class="full-screen-iframe"></iframe>
          <!-- 添加关闭按钮 -->
          <el-button class="close-iframe-btn" type="danger" @click="closeIframe">
            {{ $t('business.monitor.grafana.iframe_close') }}
          </el-button>
        </div>
        <div v-else>
          <!-- 当 iframeVisible 为 false 时，显示配置表单 -->
          <el-form ref="form" v-loading="loading" label-position="left" :rules="rules" :model="form" label-width="150px" class="config-form">
            <el-form-item style="width: 100%" :label="$t('business.monitor.grafana.address')" prop="address">
              <el-input v-model="form.address" :placeholder="'http://192.168.56.101:30020'"></el-input>
            </el-form-item>
            <el-form-item style="width: 100%" :label="$t('business.monitor.grafana.service_account_token')" prop="service_account_token">
              <el-input v-model="form.service_account_token" type="password" :placeholder="'https://grafana.com/docs/grafana/latest/administration/service-accounts/'"></el-input>
            </el-form-item>
            <el-form-item style="width: 100%" prop="enable">
              <el-checkbox v-model="form.enable">{{ $t("commons.enable.true") }}</el-checkbox>
              <el-checkbox v-model="form.default_dashboard">{{ $t("business.monitor.grafana.default_dashboard") }}</el-checkbox>
            </el-form-item>
            <div style="float: right">
              <el-form-item>
                <el-button @click="connectTest" :disabled="isSubmitGoing"
                           v-has-permissions="{resource:'monitor',verb:'create'}">{{ $t("business.monitor.grafana.test") }}
                </el-button>
                <el-button @click="remake" :disabled="isSubmitGoing"
                           v-has-permissions="{resource:'monitor',verb:'create'}">{{ $t("business.monitor.grafana.remake") }}
                </el-button>
                <el-button type="primary" @click="onSubmit" :disabled="isSubmitGoing"
                           v-has-permissions="{resource:'monitor',verb:'create'}">{{ $t("commons.button.confirm") }}
                </el-button>
              </el-form-item>
            </div>
          </el-form>
        </div>
      </el-col>
    </layout-content>
  </template>
  
<script>
import LayoutContent from "@/components/layout/LayoutContent";
import Rule from "@/utils/rules";
import { getGrafana, createGrafana, updateGrafana, testConnectGrafana, importGrafanaDashboard  } from "@/api/monitor";
  
export default {
    components: { LayoutContent },
    data() {
      return {
        address: "",
        service_account_token: "",
        form: {
          enable: false,
          default_dashboard: false,
        },
        loading: false,
        iframeVisible: false, // 控制 iframe 的可见性
        rules: {
          address: [Rule.RequiredRule],
          service_account_token: [Rule.RequiredRule],
        },
        isSubmitGoing: false,
        oldConfig: {},
      };
    },
    methods: {
      connectTest() {
        let isFormReady = false;
        this.$refs["form"].validate((valid) => {
          if (valid) {
            isFormReady = true;
          }
        });
        if (!isFormReady) {
          return;
        }
        this.loading = true;
        testConnectGrafana(this.form)
          .then((res) => {
            this.$message({
              type: "success",
              message: this.$t("business.monitor.grafana.test_result", {
                count: res.data,
              }),
            });
          })
          .finally(() => {
            this.loading = false;
          });
      },
      remake() {
        this.form = JSON.parse(JSON.stringify(this.oldConfig));
      },
      onSubmit() {
        if (this.isSubmitGoing) {
          return;
        }
        let isFormReady = false;
        this.$refs["form"].validate((valid) => {
          if (valid) {
            isFormReady = true;
          }
        });
        if (!isFormReady) {
          return;
        }
        this.isSubmitGoing = true;
        this.loading = true;
        if (!this.form.uuid) {
          createGrafana(this.form)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.create_success"),
              });
              this.list();
            })
            .finally(() => {
              this.isSubmitGoing = false;
              this.loading = false;
            });
        } else {
          updateGrafana(this.form)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("commons.msg.update_success"),
              });
              this.list();
            })
            .finally(() => {
              this.isSubmitGoing = false;
              this.loading = false;
            });
        }

        // 如果勾选导入仪表盘
        if (this.form.default_dashboard) {
            importGrafanaDashboard(this.form)
            .then(() => {
              this.$message({
                type: "success",
                message: this.$t("business.monitor.grafana.default_dashboard_success"),
              });
              this.list();
            })
            .finally(() => { 
              this.isSubmitGoing = false;
              this.loading = false;
            });
        }
      },
      list() {
        this.loading = true;
        getGrafana()
          .then((res) => {
            this.form = res.data;
            this.oldConfig = res.data;
            this.iframeVisible = this.form.enable;
          })
          .finally(() => {
            this.loading = false;
          });
      },
      closeIframe() {
        this.iframeVisible = false;
        this.form.enable = false;
        updateGrafana(this.form)
      },
    },
    created() {
      this.list();
    },
  };
</script>
  
<style scoped>
  .full-layout-content {
    position: relative;
    width: 100%;
  }
  
  .iframe-container {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .full-screen-iframe {
    width: 100%;
    height: 100%;
    border: none;
  }
  
  .close-iframe-btn {
    position: absolute;
    top: 45px;
    right: 20px;
  }
  
  .config-form {
    width: 100%; /* 表单宽度 100% */
    padding: 20px; /* 添加适当的边距 */
  }
</style>