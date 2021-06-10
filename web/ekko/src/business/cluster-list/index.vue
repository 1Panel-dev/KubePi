<template>
  <layout-content :header="$t('business.cluster.list')">
    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onImport">
            {{ $t("commons.button.import") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="del()">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" min-width="100" prop="name" fix>
        <template v-slot:default="{row}">
          <el-link href="/dashboard">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" min-width="100" prop="name" fix>
        <template v-slot:default="{row}">
          {{ row.status }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {deleteBy, searchCluster} from "@/api/clusters"

export default {
  name: "ClusterList",
  components: { ComplexTable, LayoutContent },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.del(row.name)
          }
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          { field: "name", label: this.$t("commons.table.name"), component: "FuComplexInput", defaultOperator: "eq" },
        ],
      },
      data: [],
      selects: []
    }
  },
  methods: {
    onImport () {
      this.$router.push({ name: "ClusterImport" })
    },
    search () {
      searchCluster(this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.data.items
        this.paginationConfig.total = res.data.total
      })
    },
    del (name) {
      this.$confirm(
        this.$t("commons.confirm_message.delete"),
        this.$t("commons.message_box.prompt"),
        {
          confirmButtonText: this.$t("commons.button.confirm"),
          cancelButtonText: this.$t("commons.button.cancel"),
          type: "warning"
        }
      ).then(() => {
        const ps = []
        if (name) {
          ps.push(deleteBy(name))
        } else {
          for (const item of this.selects) {
            ps.push(deleteBy(item.name))
          }
        }
        Promise.all(ps).then(() => {
          this.search()
          this.$message({
            type: "success",
            message: this.$t("commons.msg.delete_success")
          })
        }).catch(() => {
          this.search()
        })
      })
    }
  },
  created () {
    this.search()
  }
}
</script>

<style scoped>

</style>