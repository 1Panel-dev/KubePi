<template>
  <layout-content :header="$t('business.image_repos.list')">
    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig" @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'imagerepos',verb:'create'}" type="primary" size="small"
                     @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column :label="$t('commons.table.name')" prop="name" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.name }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.type')" prop="type" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.endpoint')" prop="endPoint" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('business.image_repos.repo')" prop="repoName" min-width="100" fix>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
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
import {searchRepos} from "@/api/imagerepos"
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
            this.onEdit(row.name)
          },
          disabled: () => {
            return !checkPermissions({resource: "imagerepos", verb: "create"})
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
      isSubmitGoing: false
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
        path: '/imagerepos/create',
        query: { mode: "create" }
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
