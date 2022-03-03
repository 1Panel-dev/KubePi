<template>
  <layout-content header="Persistent Volumes">
    <div style="float: left">
      <el-button v-has-permissions="{scope:'cluster',apiGroup:'',resource:'persistentvolumes',verb:'create'}"
                  type="primary" size="small"
                  @click="yamlCreate">
        YAML
      </el-button>
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'',resource:'persistentvolumes',verb:'create'}"
                  @click="onCreate">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small"
                  v-has-permissions="{scope:'cluster',apiGroup:'',resource:'persistentvolumes',verb:'delete'}"
                  :disabled="selects.length===0" @click="onDelete()">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" min-width="200px" prop="metadata.name"
                       show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="metadata.namespace">
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase ==='Bound' || row.status.phase ==='Available'" type="success" size="mini"
                     plain round>
            {{ row.status.phase }}
          </el-button>
          <el-button v-else type="" size="mini" plain round>
            {{ row.status.phase }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.storage.accessModes')" prop="metadata.labels">
        <template v-slot:default="{row}">
          <el-tag v-for="(name,index) in row.spec.accessModes " :key="index" type="info" size="mini">{{ name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.storage.volumeMode')" prop="spec.volumeMode">

      </el-table-column>
      <el-table-column :label="$t('business.storage.persistentVolumeReclaimPolicy')" prop="spec.persistentVolumeReclaimPolicy">

      </el-table-column>
      <el-table-column :label="$t('business.storage.capacity')" prop="metadata.labels">
        <template v-slot:default="{row}">
          {{ row.spec.capacity["storage"] }}
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
import {deletePvs, getPv, listPvs} from "@/api/pv"
import {checkPermissions} from "@/utils/permission"

export default {
  name: "PersistentVolumes",
  components: { ComplexTable, LayoutContent, KoTableOperations },
  data () {
    return {
      data: [],
      selects: [],
      cluster: "",
      loading: false,
      namespaces: [],
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "PersistentVolumeEdit",
              params: { name: row.metadata.name },
              query: { yamlShow: false },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "cluster", apiGroup: "", resource: "persistentvolumes", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "PersistentVolumeEdit",
              params: { name: row.metadata.name },
              query: { yamlShow: true },
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "cluster", apiGroup: "", resource: "persistentvolumes", verb: "update" })
          },
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getPv(this.cluster,row.metadata.name))
          },
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "cluster", apiGroup: "", resource: "persistentvolumes", verb: "delete" })
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
      listPvs(this.cluster, true, this.searchConfig.keywords, currentPage, pageSize).then((res) => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
      })
    },
    onCreate () {
      this.$router.push({
        name: "PersistentVolumeCreate",
        query: { yamlShow: false },
      })
    },
    yamlCreate () {
      this.$router.push({
        name: "PersistentVolumeCreateYaml",
        query: { type: "persistentvolumes" },
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
          this.ps.push(deletePvs(this.cluster, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deletePvs(this.cluster, select.metadata.name))
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
    openDetail (row) {
      this.$router.push({
        name: "PersistentVolumeDetail",
        params: { name: row.metadata.name },
        query: { yamlShow: false },
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  },
}
</script>

<style scoped>
</style>
