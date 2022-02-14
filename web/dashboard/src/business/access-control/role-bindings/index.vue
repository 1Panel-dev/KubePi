<template>
  <layout-content header="Role Bindings">
    <div style="float: left">
      <el-button
          v-has-permissions="{scope:'namespace',apiGroup:'rbac.authorization.k8s.io',resource:'rolebindings',verb:'create'}"
          type="primary" size="small"
          @click="yamlCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'rbac.authorization.k8s.io',resource:'rolebindings',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'rbac.authorization.k8s.io',resource:'rolebindings',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :selects.sync="selects" :data="data" @search="search" v-loading="loading" :pagination-config="paginationConfig"
                   :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace" show-overflow-tooltip>
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column label="Role" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span>{{ row.roleRef.kind }}/{{ row.roleRef.name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Users" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-for="(subject,index) in row.subjects" v-bind:key="index">
            <span v-if="subject.kind === 'User'">{{ subject.name }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="Groups" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-for="(subject,index) in row.subjects" v-bind:key="index">
            <span v-if="subject.kind === 'Group'">{{ subject.name }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="ServiceAccounts" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span v-for="(subject,index) in row.subjects" v-bind:key="index">
            <span v-if="subject.kind === 'ServiceAccount'">{{ subject.name }}</span>
          </span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | age }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteRoleBinding, getRoleBinding, listRoleBindings} from "@/api/rolebings"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "RoleBindings",
  components: {ComplexTable, LayoutContent, KoTableOperations},
  data() {
    return {
      data: [],
      selects: [],
      loading: false,
      cluster: "",
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "RoleBindingEdit",
              params: {namespace: row.metadata.namespace, name: row.metadata.name},
              query: {yamlShow: false}
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: 'namespace',
              apiGroup: "rbac.authorization.k8s.io",
              resource: "clusterroles",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "RoleBindingEdit",
              params: {namespace: row.metadata.namespace, name: row.metadata.name},
              query: {yamlShow: true}
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: 'cluster',
              apiGroup: "rbac.authorization.k8s.io",
              resource: "clusterroles",
              verb: "update"
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getRoleBinding(this.cluster, row.metadata.namespace, row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: 'cluster',
              apiGroup: "rbac.authorization.k8s.io",
              resource: "clusterroles",
              verb: "delete"
            })
          }
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: ""
      }
    }
  },
  methods: {
    search(resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      listRoleBindings(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate() {
      this.$router.push({
        name: "RoleBindingCreate", query: {yamlShow: false}
      })
    },
    yamlCreate() {
      this.$router.push({name: "RoleBindingCreateYaml", query: {type: "rolebindings"}})
    },
    onDelete(row) {
      this.$confirm(
          this.$t("commons.confirm_message.delete"),
          this.$t("commons.message_box.prompt"), {
            confirmButtonText: this.$t("commons.button.confirm"),
            cancelButtonText: this.$t("commons.button.cancel"),
            type: "warning",
          }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteRoleBinding(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteRoleBinding(this.cluster, select.metadata.namespace, select.metadata.name))
            }
          }
        }
        if (this.ps.length !== 0) {
          Promise.all(this.ps)
              .then(() => {
                this.search(true)
                this.$message({
                  type: "success",
                  message: this.$t("commons.msg.delete_success"),
                })
              })
              .catch(() => {
                this.search(true)
              })
        }
      })
    },
    openDetail(row) {
      this.$router.push({
        name: "RoleBindingDetail",
        params: {name: row.metadata.name, namespace: row.metadata.namespace}
      })
    }
  },
  created() {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
