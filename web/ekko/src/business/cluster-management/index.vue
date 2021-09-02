<template>
  <layout-content :header="$t('business.cluster.cluster')">
    <complex-table v-loading="loading" :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig" @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'clusters',verb:'create'}" type="primary" size="small"
                     @click="onCreate">
            {{ $t("commons.button.add") }}
          </el-button>
        </el-button-group>
      </template>

      <el-table-column min-width="60px" fix>
        <template v-slot:default="{row}">
          <el-tag type="success" v-if="row.extraClusterInfo.health">{{ $t('business.cluster.ready') }}</el-tag>
          <el-tag type="danger" v-if="!row.extraClusterInfo.health">{{ $t('business.cluster.not_ready') }}</el-tag>
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.name }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.cluster.node')" min-width="100" fix>
        <template v-slot:default="{row}">
          <el-tag>
            {{ row.extraClusterInfo.readyNodeNum }} / {{ row.extraClusterInfo.totalNodeNum }}
          </el-tag>
        </template>
      </el-table-column>


      <el-table-column min-width="200" fix>
        <template v-slot:default="{row}">
          cpu
          <el-progress type="dashboard" :width="60" :color="colors" :percentage="getCpuUsed(row)"></el-progress>
          <el-progress type="dashboard" :width="60" :color="colors" :percentage="getMemoryUsed(row)"></el-progress>
          memory
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.cluster.version')" min-width="100" fix>
        <template v-slot:default="{row}">
          <el-tag>
            {{ row.status.version }}
          </el-tag>
        </template>
      </el-table-column>


      <el-table-column :label="$t('commons.table.creat_by')" prop="createdBy" min-width="100"
                       fix/>


      <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | datetimeFormat }}
        </template>
      </el-table-column>

      <el-table-column label=" " width="90">
        <template v-slot:default="{row}">
          <el-button @click="onGotoDashboard(row)">{{ $t("business.cluster.open_dashboard") }}</el-button>
        </template>
      </el-table-column>


      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>

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
import {listClusters, deleteCluster, searchClusters} from "@/api/clusters"
import {checkPermissions} from "@/utils/permission";
import ComplexTable from "@/components/complex-table";


export default {
  name: "ClusterList",
  components: {LayoutContent, ComplexTable},
  data() {
    return {
      loading: false,
      guideDialogVisible: false,
      items: [],
      timer: null,
      data: [],
      selects: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      colors: [
        {color: '#1989fa', percentage: 0},
      ],
      searchConfig: {
        keywords: ""
      },
      buttons: [
        {
          label: this.$t("commons.button.manage"),
          icon: "el-icon-user",
          click: (row) => {
            this.onDetail(row.name)
          },
          disabled: () => {
            return !checkPermissions({resource: "clusters", verb: "update"})
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row.name)
          },
          disabled: () => {
            return !checkPermissions({resource: "clusters", verb: "delete"})
          },
        },
      ],
    }
  },
  computed: {},
  methods: {
    search() {
      this.loading = true
      const {currentPage, pageSize} = this.paginationConfig
      searchClusters(currentPage, pageSize, this.searchConfig.keywords).then(data => {
        this.loading = false
        this.data = data.data.items
        this.paginationConfig.total = data.data.total
      })
    },
    onCreate() {
      this.$router.push({name: "ClusterCreate"})
      this.search()
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
          this.search()
        })
      });
    },
    onGotoDashboard(row) {
      if (row.accessable) {
        window.open(`/dashboard?cluster=${row.name}`, "_self")
      } else {
        this.$message.error(this.$t('business.cluster.user_not_in_cluster'))
      }
    },

    getCpuUsed(item) {
      const r = Math.round((item.extraClusterInfo.cpuRequested / item.extraClusterInfo.cpuAllocatable) * 100)
      return r > 100 ? 100 : r
    },

    getMemoryUsed(item) {
      const r = Math.round((item.extraClusterInfo.memoryRequested / item.extraClusterInfo.memoryAllocatable) * 100)
      return r > 100 ? 100 : r
    },

    onGuildSubmit() {
      this.$router.push({name: "ClusterCreate"})
    },
    pullingClusterStatus() {
      this.timer = setInterval(() => {
        listClusters().then(data => {
          this.items = data.data;
        })
      }, 3000)
    },
    onVueCreated() {
      listClusters().then(data => {
        this.items = data.data.filter((item) => {
          if (this.$store.getters.isAdmin) {
            return true
          }
          if (!checkPermissions({resource: "clusters", verb: "create"}) && !checkPermissions({
            resource: "clusters",
            verb: "delete"
          })) {
            return item.accessable
          }
        });
        if (this.items.length === 0 && checkPermissions({resource: 'clusters', verb: 'create'})) {
          this.guideDialogVisible = true
        }
        this.pullingClusterStatus()
      })
    }
  },
  created() {
    this.search()
    // this.onVueCreated()
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