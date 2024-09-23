<template>
  <layout-content :header="$t('business.cluster.list')">
    <div style="float: left">
      <el-button v-has-permissions="{resource:'clusters',verb:'create'}" type="primary" size="small" @click="onCreate">{{ $t("commons.button.import") }}</el-button>
      <el-button v-has-permissions="{resource:'clusters',verb:'delete'}" :disabled="selects.length===0" type="primary" size="small" @click="onDelete()">{{ $t("commons.button.delete") }}</el-button>
      <el-button  :disabled="selects.length===0" type="primary" size="small" @click="onExportAllHelmReleases()">{{ $t("commons.button.export") }}</el-button>
    </div>
    <complex-table v-loading="loading" :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig" @search="search"
                   element-loading-background="rgba(0, 0, 0, 0.8)">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.status')" min-width="80px" fix>
        <template v-slot:default="{row}">
          <el-tag type="success" v-if="row.extraClusterInfo.health">{{ $t('business.cluster.ready') }}</el-tag>
          <el-tooltip class="item" effect="dark" :content="row.extraClusterInfo.message" placement="right">
            <el-tag type="danger" v-if="!row.extraClusterInfo.health">{{ $t('business.cluster.not_ready') }}</el-tag>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="140px" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span v-if="row.extraClusterInfo.health" class="span-link" @click="onGotoDashboard(row)">{{ row.name }}</span>
          <span v-else>{{ row.name }}</span>
        </template>
      </el-table-column>


      <el-table-column :label="$t('business.cluster.label')" align="left" min-width="160" fix>
        <template v-slot:default="{row}">
          <div style="font-size: 20px">
            <div v-if="row.labels" :key="row.k">
              <div v-for="(key,index) in row.labels" :key="index">

                <el-tag v-if="!checkPrem({resource: 'clusters', verb: 'update'})" type="info" size="mini"
                        @close="onDeleteLabel(row,key)">
                  {{ key }}
                </el-tag>

                <el-tag v-if="checkPrem({resource: 'clusters', verb: 'update'})" closable type="info" size="mini"
                        @close="onDeleteLabel(row,key)">
                  {{ key }}
                </el-tag>

                <br/>
              </div>
            </div>
            <div style="padding-top: 5px;">
              <el-popover
                placement="top"
                width="160"
                :title="$t('commons.button.add')+$t('business.cluster.label')"
                trigger="manual"
                v-model="row.showAddLabelVisible">
                <div style="text-align: right; margin: 0">
                  <el-form :ref="row.name" :model="row.form" :rules="labelRules" label-position="top">
                    <el-form-item size="mini" prop="key">
                      <el-input type="text" v-model="row.form.key" placeholder="key"></el-input>
                    </el-form-item>
                  </el-form>
                  <el-button size="mini" type="text" @click="row.showAddLabelVisible = false">
                    {{ $t('commons.button.cancel') }}
                  </el-button>
                  <el-button type="primary" size="mini" @click="onAddLabelSubmit(row)">
                    {{ $t('commons.button.confirm') }}
                  </el-button>
                </div>
                <i :id="row.id" class="el-icon-circle-plus-outline" slot="reference" style="cursor: pointer"
                  @click="onAddLabel(row)" v-has-permissions="{resource:'clusters',verb:'update'}"></i>
              </el-popover>
            </div>
          </div>
        </template>
      </el-table-column>


      <el-table-column :label="$t('business.cluster.nodes')" min-width="80" fix>
        <template v-slot:default="{row}">
          <el-tag>
            {{ row.extraClusterInfo.readyNodeNum }} / {{ row.extraClusterInfo.totalNodeNum }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column :label="$t('business.cluster.version')" min-width="100" fix>
        <template v-slot:default="{row}">
          <el-tag>
            {{ row.status.version }}
          </el-tag>
        </template>
      </el-table-column>


      <el-table-column :label="$t('commons.table.age')" min-width="120" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | ageFormat }}
        </template>
      </el-table-column>

      <el-table-column label=" " width="100">
        <template v-slot:default="{row}">
          <el-button @click="onGotoDashboard(row)" :disabled="!row.extraClusterInfo.health">
            {{ $t("business.cluster.open_dashboard") }}
          </el-button>
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
import {deleteCluster, listClusters, searchClusters, updateCluster} from "@/api/clusters"
import {checkPermissions} from "@/utils/permission";
import ComplexTable from "@/components/complex-table";
import Rule from "@/utils/rules"
import {downloadHelmReleases} from "@/utils/helm"

export default {
  name: "ClusterList",
  components: {LayoutContent, ComplexTable},
  data() {
    return {
      loading: false,
      guideDialogVisible: false,
      showAddLabelVisible: false,
      items: [],
      timer: null,
      data: [],
      selects: [],
      labelRules: {
        key: [
          Rule.RequiredRule,
          {
            min: 0,
            max: 10,
            message: this.$t("commons.validate.limit", [1, 10]),
            trigger: "blur"
          }
        ]
      },
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        quickPlaceholder: this.$t('commons.search.quickSearch'),
        components: [
          {field: "name", label: this.$t("commons.table.name"), component: "FuComplexInput", defaultOperator: "eq"},
          {
            field: "labels",
            label: this.$t("business.cluster.label"),
            component: "FuComplexInput",
            defaultOperator: "like"
          },
        ]
      },
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.onEdit(row.name)
          },
          disabled: () => {
            return !(checkPermissions({resource: "clusters", verb: "update"}))
          }
        },
        {
          label: this.$t("commons.button.rbac_manage"),
          icon: "el-icon-user",
          click: (row) => {
            this.onDetail(row.name)
          },
          disabled: (row) => {
            return !(checkPermissions({resource: "clusters", verb: "authorization"}) && row.extraClusterInfo.health)
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
  methods: {
    search(conditions) {
      this.loading = true
      const {currentPage, pageSize} = this.paginationConfig
      searchClusters(currentPage, pageSize, conditions).then(data => {
        this.loading = false
        this.data = data.data.items
        this.paginationConfig.total = data.data.total
        this.data.forEach(d => {
          this.$set(d, "showAddLabelVisible", false)
          this.$set(d, "k", 0)
          this.$set(d, "form", {
            key: "",
          },)
        })
      })
    },
    onCreate() {
      this.$router.push({name: "ClusterCreate"})
      this.search()
    },
    onAddLabel(row) {
      this.data.forEach(d => {
        if (d.showAddLabelVisible) {
          d.showAddLabelVisible = false
        }
      })
      row.form.key = ""
      row.showAddLabelVisible = true
    },

    onDeleteLabel(row, key) {
      row.labels.splice(row.labels.indexOf(key), 1)
      updateCluster(row.name, {"labels": row.labels, "withLabel": true}).then(() => {
        row.k++
      })
    },
    onAddLabelSubmit(row) {
      this.$refs[row.name].validate((valid) => {
        if (valid) {
          let key = row.form.key
          if (!row.labels) {
            row.labels = []
          }
          if (row.labels.indexOf(key) === -1) {
            row.labels.push(key)
            updateCluster(row.name, {"labels": row.labels, "withLabel": true}).then(() => {
              row.showAddLabelVisible = false
            })
          } else {
            this.$message({
              type: 'warning',
              message: this.$t("commons.msg.duplicate_failed"),
            });
          }
          row.showAddLabelVisible = false
        }
      })
    },
    onDetail(name) {
      this.$router.push({name: "ClusterMembers", params: {name: name}})
    },
    onEdit(name) {
      this.$router.push({name: "ClusterEdit", params: {name: name}})
    },
    onDelete(name) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        this.ps = []
        if (name) {
          this.ps.push(deleteCluster(name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteCluster(select.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
            .then(() => {
              this.search()
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.search()
            })
        }
      });
    },
    onGotoDashboard(row) {
      if (row.accessable) {
        sessionStorage.removeItem("namespace")
        const url = `${process.env.VUE_APP_DASHBOARD_URL_PREFIX}/dashboard?cluster=${row.name}`
        window.open(url, "_blank")
      } else {
        this.$message.error(this.$t('business.cluster.user_not_in_cluster'))
      }
    },
    checkPrem(obj) {
      return checkPermissions(obj)
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
    //导出选中集群的所有原始helm releases
    async onExportAllHelmReleases(){
       this.loading=true
       await downloadHelmReleases(this.selects)
       this.loading=false
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