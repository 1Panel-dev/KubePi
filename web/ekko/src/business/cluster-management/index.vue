<template>
  <layout-content :header="$t('business.cluster.cluster')">

    <div style="float: right">
      <el-switch
          v-model="hiddenUnAccessCluster"
          @change="onHiddenUnAccessClusterChange"
          :active-text="$t('business.cluster.hidden_cluster')">
      </el-switch>
    </div>
    <br>
    <el-row>
      <el-col :span="4" v-for="(item,index) in items" :key="index">
        <el-card class="card_header box-card cluster-card" style=""
                 shadow="hover">
          <div>
            <div slot="header" class="clearfix">
              <b style="font-size: 20px">{{ item.name }}</b>
              <div style="float: right">
                <el-dropdown trigger="click" @command="handleCommand">
                  <el-button type="text" size="large" class="bottom-button"><i
                      class="el-icon-more"></i></el-button>

                  <el-dropdown-menu slot="dropdown">
                    <el-dropdown-item v-has-permissions="{resource:'clusters',verb:'update'}"
                                      :command="{name:item.name,action:'manage'}">
                      {{ $t('business.cluster.management') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-has-permissions="{resource:'clusters',verb:'delete'}"
                                      :command="{name:item.name,action:'delete'}">
                      {{ $t('commons.button.delete') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </el-dropdown>

              </div>
              <hr style="border-color: gray"/>
            </div>
            <div>
              <el-row>
                <el-col :span="4" style="padding-top: 10px">
                  <img height="48px" src="~@/assets/kubernetes.png" alt="kubernetes.png">
                </el-col>
                <el-col :span="20">
                  <div style="margin-top: 20px"
                       v-if="item.status.phase!=='Initializing' && item.status.phase!=='Failed' ">
                    <ul>
                      <li style="list-style-type: none;margin: 5px">
                        {{ $t('business.cluster.cluster_version') }}:{{ item.status.version }}
                      </li>
                    </ul>
                  </div>
                  <div v-if="item.status.phase==='Initializing'" style="margin-left: 30px;margin-top: 28px">
                    <i class="el-icon-loading"></i>初始化中...
                  </div>
                </el-col>
              </el-row>
            </div>
            <div class="bottom clearfix">
              <el-button type="text" size="large" class="bottom-button"
                         v-has-permissions="{resource:'clusters',verb:'get'}"
                         @click="onGotoDashboard(item.name)">
                {{ $t('business.cluster.open_dashboard') }}>
              </el-button>
            </div>

          </div>

        </el-card>


      </el-col>
      <el-col :span="4">
        <el-card v-has-permissions="{resource:'clusters',verb:'update'}" class="card_header box-card cluster-card"
                 shadow="hover">
          <div @click="onCreate" style="text-align: center;padding-top: 45px;font-size: 36px;cursor: pointer">
            <i class="el-icon-plus"></i>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog
        :title="$t('commons.header.guide')"
        :visible.sync="guideDialogVisible"
        width="30%">

      <div>
        <div style="font-size: 14px">
          {{ $t('commons.header.guide_text') }}
        </div>

        <div style="font-size: 14px;margin-top: 50px">
          <i class="el-icon-document"></i>{{ $t('commons.header.help_doc') }}：&nbsp;
          <el-link href="https://doc.kubeoperator.io">https://doc.kubeoperator.io</el-link>
          <br/>
          <br/>
          <i class="el-icon-chat-square"></i>{{ $t('commons.header.support') }}：&nbsp;
          <el-link class="el-link" href="mailto:support@fit2cloud.com">support@fit2cloud.com</el-link>
        </div>
      </div>
      <div slot="footer" class="dialog-footer">
        <el-button @click="guideDialogVisible = false">{{ $t('commons.button.skip') }}</el-button>
        <el-button type="primary" @click="onGuildSubmit">{{ $t('commons.button.confirm') }}</el-button>
      </div>
    </el-dialog>


  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listClusters, deleteCluster} from "@/api/clusters"
import {checkPermissions} from "@/utils/permission";


export default {
  name: "ClusterList",
  components: {LayoutContent},
  data() {
    return {
      guideDialogVisible: false,
      items: [],
      hiddenUnAccessCluster: false,
      timer: null,
    }
  },
  methods: {

    handleCommand(command) {
      const name = command.name
      switch (command.action) {
        case "manage":
          this.onDetail(name)
          break
        case "delete":
          this.onDelete(name)
          break
      }
    },

    onCreate() {
      this.$router.push({name: "ClusterCreate"})
    },
    onDetail(name) {
      this.$router.push({name: "ClusterMembers", params: {name: name}})
    },
    onDelete(name) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteCluster(name).then(() => {
          this.$message({
            type: 'success',
            message: this.$t("commons.msg.delete_success"),
          });
          this.onVueCreated()
        })
      });
    },
    onGotoDashboard(name) {
      window.open(`/dashboard?cluster=${name}`, "_self")
    },
    onHiddenUnAccessClusterChange() {
      this.onVueCreated()
    },
    onGuildSubmit() {
      this.$router.push({name: "ClusterCreate"})
    },
    pullingClusterStatus() {
      this.timer = setInterval(() => {
        listClusters(!this.hiddenUnAccessCluster).then(data => {
          this.items = data.data;
        })
      }, 3000)
    },
    onVueCreated() {
      listClusters(!this.hiddenUnAccessCluster).then(data => {
        this.items = data.data;
        if (this.items.length === 0 && checkPermissions({resource: 'clusters', verb: 'create'})) {
          this.guideDialogVisible = true
        }
        this.pullingClusterStatus()
      })
    }
  },
  created() {
    this.onVueCreated()
  }, destroyed() {
    clearInterval(this.timer)
  }
}
</script>

<style scoped>
.clearfix:before,
.clearfix:after {
  display: table;
  content: "";
}

.clearfix:after {
  clear: both
}

.bottom {
  margin-top: 13px;
  line-height: 12px;
}

.bottom-button {
  padding: 0;
  float: right;
}

.cluster-card {
  margin-left: 10px;
  margin-top: 20px;
  min-height: 169px;
}
</style>