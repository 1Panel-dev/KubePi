<template>
  <layout-content :header="$t('business.chart.app_installed')">
    <complex-table v-loading="loading" :data="data" :search-config="searchConfig" :selects.sync="selects"
                   :pagination-config="paginationConfig"
                   @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'charts',verb:'delete'}" type="primary" size="small"
                     @click="onDelete()" :disabled="selects.length === 0">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
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
      <el-table-column :label="$t('business.chart.version')" prop="version" fix>
        <template v-slot:default="{row}">
          {{ row.chart.metadata.version }}
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
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          disabled: () => {
            return !checkPermissions({ resource: "charts", verb: "delete" })
          },
          click: (row) => {
            this.onDelete(row)
          }
        },
        {
          label: this.$t("business.chart.upgrade"),
          icon: "el-icon-edit",
          disabled: () => {
            return !checkPermissions({ resource: "charts", verb: "upgrade" })
          },
          click: (row) => {
            this.onUpgrade(row)
          }
        },
      ]
    }
  },
  methods: {
    search (change) {
      if (change) {
        this.paginationConfig = {
          currentPage: 1,
          pageSize: 10,
          total: 0,
        }
      }
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
          this.ps.push(deleteApp(this.cluster, row.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteApp(this.cluster, select.name))
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
