<template>
  <layout-content :header="$t('business.image_repos.list')">
    <div style="float: left; margin-bottom: 20px">
      <el-button v-has-permissions="{resource:'imagerepos',verb:'create'}" type="primary" size="small" @click="onCreate">{{ $t("commons.button.add") }}</el-button>
      <el-button v-has-permissions="{resource:'imagerepos',verb:'delete'}" :disabled="selects.length===0" type="primary" size="small" @click="onDelete()">{{ $t("commons.button.delete") }}</el-button>
    </div>

    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig" @search="search">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="100" fix>
        <template v-slot:default="{row}">
          <el-link @click="onDetail(row.name)">{{ row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.type')" prop="type" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.endpoint')" prop="endPoint" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.downloadUrl')" prop="downloadUrl" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.repo')" prop="repoName" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('commons.table.age')" min-width="80" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | ageFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')"/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {deleteRepo, searchRepos} from "@/api/imagerepos"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "ImageRepo",
  components: { ComplexTable, LayoutContent },
  props: {},
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              path: `/imagerepos/${row.name}/edit`,
              query: {mode: "edit" }
            })
          },
          disabled: () => {
            return !checkPermissions({ resource: "imagerepos", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row.name)
          },
          disabled: () => {
            return !checkPermissions({ resource: "imagerepos", verb: "delete" })
          }
        },
      ],
      searchConfig: {
        quickPlaceholder: this.$t("commons.search.quickSearch"),
        components: [
          { field: "name", label: this.$t("commons.table.name"), component: "FuComplexInput", defaultOperator: "eq" },
        ]
      },
      selects: [],
      data: [],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
    }
  },
  methods: {
    search (conditions) {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      searchRepos(currentPage, pageSize, conditions).then(data => {
        this.loading = false
        this.data = data.data.items
        this.paginationConfig.total = data.data.total
      })
    },
    onCreate () {
      this.$router.push({
        path: "/imagerepos/create",
        query: { mode: "create" }
      })
    },
    onDelete (name) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning"
      }).then(() => {
        this.ps = []
        if (name) {
          this.ps.push(deleteRepo(name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteRepo(select.name))
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
      })
    },
    onDetail(repo) {
      this.$router.push({
        name: "ImageRepoDetail",
        params: { repo: repo }
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
