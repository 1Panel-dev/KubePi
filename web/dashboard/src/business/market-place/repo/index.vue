<template>
  <layout-content :header="$t('business.chart.chart')">
    <complex-table v-loading="loading"  :selects.sync="selects" :data="data"
                   @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'charts',verb:'create'}" type="primary" size="small"
                     @click="onCreate">
            {{ $t("commons.button.add") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="80" fix>
        <template v-slot:default="{row}">
          <span :style="{color: row.color}">{{ row.name }}</span>
        </template>
      </el-table-column>
<!--      <el-table-column :label="$t('business.chart.type')" prop="type" min-width="80" fix>-->
<!--        <template v-slot:default="{row}">-->
<!--          {{ row.type }}-->
<!--        </template>-->
<!--      </el-table-column>-->
      <el-table-column :label="'URL'" prop="url" min-width="80" fix>
        <template v-slot:default="{row}">
          {{ row.url }}
        </template>
      </el-table-column>
      <!--      <el-table-column :label="$t('business.chart.branch')" prop="branch" min-width="80" fix>-->
      <!--        <template v-slot:default="{row}">-->
      <!--          {{ row.spec.gitBranch }}-->
      <!--        </template>-->
      <!--      </el-table-column>-->
<!--      <el-table-column :label="$t('commons.table.created_time')" min-width="120" fix>-->
<!--        <template v-slot:default="{row}">-->
<!--          {{ row.createAt | datetimeFormat }}-->
<!--        </template>-->
<!--      </el-table-column>-->
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {deleteChart, listRepos} from "@/api/charts"
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
          click: (row) => {
            this.onDetail(row.name)
          },
          disabled: () => {
            return !checkPermissions({ resource: "charts", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row.name)
          },
          disabled: () => {
            return !checkPermissions({ resource: "charts", verb: "delete" })
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
    onDetail (name) {
      this.$router.push({ name: "RepoEdit", params: { name: name } })
    },
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
        deleteChart(name).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success"),
          })
          this.search()
        }).finally(() => {
          this.isSubmitGoing = false
        })
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
