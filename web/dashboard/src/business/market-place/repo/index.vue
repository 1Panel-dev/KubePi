<template>
  <layout-content :header="$t('business.chart.chart')">
    <div style="float: left; margin-bottom: 10px">
      <el-button v-has-permissions="{scope:'cluster',apiGroup:'kubepi.org',resource:'appmarkets',verb:'create'}"
                  type="primary" size="small"
                  @click="onCreate">
        {{ $t("commons.button.add") }}
      </el-button>
    </div>
    <complex-table v-loading="loading" :selects.sync="selects" :data="data"
                   @search="search">
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" fix>
        <template v-slot:default="{row}">
          <span :style="{color: row.color}">{{ row.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="'URL'" prop="url" min-width="80" fix>
        <template v-slot:default="{row}">
          {{ row.url }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {deleteRepo, listRepos, syncRepo} from "@/api/charts"
import ComplexTable from "@/components/complex-table"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "RepoManagement",
  components: { ComplexTable, LayoutContent },
  props: {},
  data () {
    return {
      loading: false,
      items: [],
      data: [],
      selects: [],
      cluster: "",
      isSubmitGoing: false,
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "kubepi.org",
              resource: "appmarkets",
              verb: "update"
            })
          },
          click: (row) => {
            this.onEdit(row)
          }
        },
        {
          label: this.$t("commons.button.sync"),
          icon: "el-icon-refresh",
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "kubepi.org",
              resource: "appmarkets",
              verb: "update"
            })
          },
          click: (row) => {
            this.sync(row.name)
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row.name)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "kubepi.org",
              resource: "appmarkets",
              verb: "delete"
            })
          },
        },
      ]
    }
  },
  methods: {
    search () {
      this.loading = true
      listRepos(this.cluster).then(data => {
        this.loading = false
        this.data = data.data
      })
    },
    onCreate () {
      this.$router.push({ name: "RepoCreate" })
    },
    onEdit (row) {
      this.$router.push({ name: "RepoEdit", params: { name: row.name } })
    },
    // onDetail (name) {
    //   this.$router.push({ name: "RepoEdit", params: { name: name } })
    // },
    onDelete (name) {
      if (this.isSubmitGoing) {
        return
      }
      this.isSubmitGoing = true
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning"
      }).then(() => {
        deleteRepo(this.cluster, name).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success"),
          })
          this.search()
        }).finally(() => {
          this.isSubmitGoing = false
        })
      })
    },
    sync (name) {
      this.loading = true
      syncRepo(this.cluster, name).then(() => {
        this.$message({
          type: "success",
          message: this.$t("commons.msg.sync_success"),
        })
      }).finally(() => {
        this.loading = false
      })
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
