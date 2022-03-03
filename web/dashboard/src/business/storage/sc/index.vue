<template>
  <layout-content header="Storage Class">
    <div style="float: left">
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'storage.k8s.io',resource:'storageclasses',verb:'create'}"
                  @click="onCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'storage.k8s.io',resource:'storageclasses',verb:'delete'}"
                  :disabled="selects.length===0" @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.storage.provisioner')" prop="provisioner">
        <template v-slot:default="{row}">
          {{ row.provisioner }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.storage.reclaimPolicy')" prop="reclaimPolicy">
        <template v-slot:default="{row}">
          {{ row.reclaimPolicy }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.storage.default')">
        <template v-slot:default="{row}">
          <i v-if="checkDefault(row)"
             class="el-icon-check"></i>
          <i v-else class="el-icon-minus"></i>
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
import ComplexTable from "@/components/complex-table/index"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {changeStorageClass, deleteStorageClasses, getStorageClass, listStorageClasses} from "@/api/storageclass"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "StorageClasses",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
      selects: [],
      cluster: "",
      loading: false,
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "StorageClassEdit",
              params: { name: row.metadata.name },
              query: { yamlShow: false },
            })
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "storage.k8s.io",
              resource: "storageclasses",
              verb: "update"
            })
          },
        },
        {
          label: this.$t("business.storage.change_default"),
          icon: "el-icon-check",
          click: (row) => {
            this.changeDefault(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "storage.k8s.io",
              resource: "storageclasses",
              verb: "update"
            })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", getStorageClass(this.cluster, row.metadata.name))
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({
              scope: "cluster",
              apiGroup: "storage.k8s.io",
              resource: "storageclasses",
              verb: "delete"
            })
          },
        },
      ],
      paginationConfig: {
        currentPage: 1,
        pageSize: 10,
        total: 0,
      },
      searchConfig: {
        keywords: "",
      },
    }
  },
  methods: {
    search () {
      this.loading = true
      const { currentPage, pageSize } = this.paginationConfig
      listStorageClasses(this.cluster, true, this.searchConfig.keywords, currentPage, pageSize).then((res) => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "StorageClassCreateYaml",
        query: { type: "storageclasses" },
      })
    },
    onDelete (row) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.prompt"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: "warning",
      }).then(() => {
        this.ps = []
        if (row) {
          this.ps.push(deleteStorageClasses(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteStorageClasses(this.cluster, select.metadata.name))
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
    openDetail (row) {
      this.$router.push({
        name: "StorageClassDetail",
        params: { name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
    changeDefault (row) {
      if (this.checkDefault(row)) {
        changeStorageClass(this.cluster, row.metadata.name, { metadata: { annotations: { "storageclass.kubernetes.io/is-default-class": "false" } } }).then(() => {
          this.$message({
            type: "success",
            message: this.$t("commons.msg.update_success"),
          })
          this.search(true)
        })
      } else {
        listStorageClasses(this.cluster).then(res => {
          this.ps = []
          for (const item of res.items) {
            if (this.checkDefault(item)) {
              this.ps.push(changeStorageClass(this.cluster, item.metadata.name, { metadata: { annotations: { "storageclass.kubernetes.io/is-default-class": "false" } } }))
            }
          }
          this.ps.push(changeStorageClass(this.cluster, row.metadata.name, { metadata: { annotations: { "storageclass.kubernetes.io/is-default-class": "true" } } }))
          if (this.ps.length !== 0) {
            Promise.all(this.ps)
              .then(() => {
                this.$message({
                  type: "success",
                  message: this.$t("commons.msg.update_success"),
                })
                this.search(true)
              })
              .catch(() => {
                this.search(true)
              })
          }
        })
      }
    },
    checkDefault (row) {
      return row.metadata.annotations && row.metadata.annotations["storageclass.kubernetes.io/is-default-class"] && row.metadata.annotations["storageclass.kubernetes.io/is-default-class"] === "true"
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  },
}
</script>

<style scoped>
</style>
