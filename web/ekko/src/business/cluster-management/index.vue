<template xmlns:el-col="http://www.w3.org/1999/html">
  <layout-content :header="$t('business.cluster.cluster')">

    <div style="float: right">
      <el-switch
          v-model="hiddenUnAccessCluster"
          @change="onHiddenUnAccessClusterChange"
          active-text="隐藏不可访问集群">
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
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {listClusters, deleteCluster} from "@/api/clusters"


export default {
  name: "ClusterList",
  components: {LayoutContent},
  data() {
    return {
      items: [],
      hiddenUnAccessCluster: false
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
    onVueCreated() {
      listClusters(!this.hiddenUnAccessCluster).then(data => {
        this.items = data.data;
      })
    }
  },
  created() {
    this.onVueCreated()
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