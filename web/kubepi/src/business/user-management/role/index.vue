<template>
  <layout-content :header="$t('business.user.role_list')">
    <complex-table :search-config="searchConfig" :selects.sync="selects" :data="data"
                   :pagination-config="paginationConfig" @search="search">
      <template #header>
        <el-button-group>
          <el-button v-has-permissions="{resource:'roles',verb:'create'}" type="primary"
                     size="small"
                     @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
        </el-button-group>
      </template>
      <el-table-column :label="$t('commons.table.name')" min-width="100" show-overflow-tooltip fix>
        <template v-slot:default="{row}">
          <span class="span-link" @click="onDetail(row.name)">
            {{ row.name }}
          </span>
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.description')" min-width="150" fix>
        <template v-slot:default="{row}">
          {{ translate(row.description) }}
        </template>
      </el-table-column>

      <el-table-column :label="$t('commons.table.built_in')" min-width="60" fix>
        <template v-slot:default="{row}">
          {{ $t("commons.bool." + row.builtIn) }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.age')" min-width="60" fix>
        <template v-slot:default="{row}">
          {{ row.createAt |ageFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {searchRoles, deleteRole} from "@/api/roles"
import {checkPermissions} from "@/utils/permission";

export default {
  name: "Role",
  components: {LayoutContent, ComplexTable},
  data() {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.onEdit(row.name)
          },
          disabled: (row) => {
            return row.builtIn || !checkPermissions({resource: 'roles', verb: 'update'})
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row.name)
          },
          disabled: (row) => {
            return row.builtIn || !checkPermissions({resource: 'roles', verb: 'delete'})
          }
        },
      ],
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
            field: "builtIn", label: this.$t("commons.table.built_in"), component: "FuComplexSelect",
            options: [{label: this.$t("commons.bool.true"), value: "true"}, {
              label: this.$t("commons.bool.false"),
              value: "false"
            }],
          },
        ]
      },
      data: [],
      selects: [],
      isSubmitGoing: false,
    }
  },
  methods: {
    search(conditions) {
      console.log(conditions)
      this.loading = true
      const {currentPage, pageSize} = this.paginationConfig
      searchRoles(currentPage, pageSize, conditions).then(data => {
        this.loading = false
        this.data = data.data.items
        this.paginationConfig.total = data.data.total
      })
    },
    onCreate() {
      this.$router.push({name: "RoleCreate"})

    },

    translate(a) {
      return a.startsWith("i18n_") ? this.$t(a) : a
    },

    onDetail(name) {
      this.$router.push({name: "RoleDetail", params: {name: name}})
    },
    onEdit(name) {
      this.$router.push({name: "RoleEdit", params: {name: name}})
    },
    onDelete(name) {
      if (this.isSubmitGoing) {
        return
      }
      this.isSubmitGoing = false
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteRole(name).then(() => {
          this.$message({
            type: 'success',
            message: this.$t("commons.msg.delete_success"),
          });
          this.search()
        })
      }).finally(() => {
        this.isSubmitGoing = false
      });
    }
  },
  created() {
    this.search()
  }
}
</script>

<style scoped>

</style>