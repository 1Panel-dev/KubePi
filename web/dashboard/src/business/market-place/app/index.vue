<template>
  <layout-content :header="$t('business.chart.app_installed')">
    <div style="float: left">
      <el-button v-has-permissions="{scope:'cluster',apiGroup:'kubepi.org',resource:'appmarkets',verb:'create'}" type="primary" size="small"
                  @click="onDelete()" :disabled="selects.length === 0">
        {{ $t("commons.button.delete") }}
      </el-button>
      <el-button  type="primary" size="small"
                  @click="onExportSelected()" :disabled="selects.length === 0">
        {{ $t("commons.button.export") }}
      </el-button>
    </div>
    <complex-table v-loading="loading" :data="data" :search-config="searchConfig" :selects.sync="selects"
                   :pagination-config="paginationConfig"
                   @search="search">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" fix>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="namespace" fix>
      </el-table-column>
      <el-table-column :label="'Chart'" prop="chart" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.name }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="status" fix>
        <template v-slot:default="{row}">
          {{ row.info.status }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.chart.version')" prop="version" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.version }}
        </template>
      </el-table-column>
      <el-table-column :label="'App '+$t('business.chart.version')" prop="version" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.appVersion }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="info.first_deployed" fix>
        <template v-slot:default="{row}">
          {{ row.info.first_deployed | age }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {deleteApp, searchInstalled} from "@/api/charts"
import ComplexTable from "@/components/complex-table"
import {checkPermissions} from "@/utils/permission"
import KoTableOperations from "@/components/ko-table-operations"
import {downloadHelmReleases} from "@/utils/helm"
export default {
  name: "Apps",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  props: {},
  data () {
    return {
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      },
      apps: [],
      cluster: "",
      data: [],
      loading: false,
      selects: [],
      buttons: [
        {
          label: this.$t("commons.button.edit")+ "/" + this.$t("business.chart.upgrade"),
          icon: "el-icon-edit",
          disabled: () => {
            return !checkPermissions({scope:'cluster',apiGroup:'kubepi.org',resource:'appmarkets',verb:'update'})
          },
          click: (row) => {
            this.onUpgrade(row)
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          disabled: () => {
            return !checkPermissions({scope:'cluster',apiGroup:'kubepi.org',resource:'appmarkets',verb:'delete'})
          },
          click: (row) => {
            this.onDelete(row)
          }
        },
{
          label: "export",
          icon: "el-icon-download",
          click: (row) => {
            this.onDownloadHelmReleases(row)
          }
        },
      ]
    }
  },
  methods: {
    search () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchInstalled(this.cluster, currentPage, pageSize, this.searchConfig.keywords).then(res => {
        this.data = res.data.items
        this.paginationConfig.total = res.data.total
        this.loading = false
      })
    },
    onDelete (row) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"), {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning",
        }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteApp(this.cluster, row.namespace,row.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteApp(this.cluster,  select.namespace,select.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          this.loading = true
          Promise.all(this.ps)
            .then(() => {
              this.search()
              this.loading = false
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.loading = false
              this.search()
            })
        }
      })
    },
    onUpgrade (row) {
      this.$router.push({ name: "AppUpgrade", params: {name: row.name } })
},
    //导出helm releases原始chart
    async onDownloadHelmReleases(row) {
       this.loading=true
       await downloadHelmReleases(this.cluster,[row])
       this.loading=false
    },
    //导出选中行的helm releases原始chart
    async onExportSelected(){
       this.loading=true
       await downloadHelmReleases(this.cluster,this.selects)
       this.loading=false
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
