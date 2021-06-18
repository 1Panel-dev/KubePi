<template>
  <layout-content header="Namespaces" v-loading="loading">
    <complex-table :selects.sync="selects" :data="data" :pagination-config="page" @search="search()">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.create") }}
          </el-button>
          <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()">
            {{ $t("commons.button.delete") }}
          </el-button>
        </el-button-group>
      </template>
      <template #toolbar>
        <el-input :placeholder="$t('commons.button.search')" suffix-icon="el-icon-search" clearable v-model="searchName" @change="search(true)" @clear="search(true)"></el-input>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" fix>
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.cluster.label')" prop="metadata.labels" min-width="200px">
        <template v-slot:default="{row}">
          <el-tag v-for="(value,key,index) in row.metadata.labels" v-bind:key="index" type="info" size="mini">
            {{ key }}={{ value }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.status')" prop="metadata.status" fix>
        <template v-slot:default="{row}">
          <el-button v-if="row.status.phase ==='Active'" type="success" size="mini" plain round>
            {{
              row.status.phase
            }}
          </el-button>
          <el-button v-if="row.status.phase ==='Terminating'" type="warning" size="mini" plain round>
            {{
              row.status.phase
            }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" prop="metadata.creationTimestamp" fix>
        <template v-slot:default="{row}">
          {{ row.metadata.creationTimestamp | datetimeFormat }}
        </template>
      </el-table-column>
      <ko-table-operations :buttons="buttons" :label="$t('commons.table.action')"></ko-table-operations>
    </complex-table>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import ComplexTable from "@/components/complex-table"
import {listNamespace, deleteNamespace} from "@/api/namespaces"
import KoTableOperations from "@/components/ko-table-operations"

export default {
  name: "NamespaceList",
  components: { KoTableOperations, ComplexTable, LayoutContent },
  data () {
    return {
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "NamespaceEdit",
              params: { name: row.metadata.name, cluster: this.clusterName, yamlShow: false }
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "NamespaceEdit",
              params: { name: row.metadata.name, cluster: this.clusterName, yamlShow: true }
            })
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          }
        },
      ],
      data: [],
      selects: [],
      loading: false,
      page: {
        pageSize: 10,
        nextToken: "",
      },
      clusterName: "test1",
      searchName: ""
    }
  },
  methods: {
    onCreate () {
      this.$router.push({ name: "NamespaceCreate" })
    },
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: "",
        }
      }
      listNamespace(this.clusterName, this.page.pageSize, this.page.nextToken,this.searchName).then((res) => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
      }).catch(error => {
        console.log(error)
      }).finally(() => {
        this.loading = false
      })
    },
    openDetail (row) {
      this.$router.push({ name: "NamespaceDetail", params: { name: row.metadata.name, cluster: "test1" } })
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
          this.ps.push(deleteNamespace(this.clusterName, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteNamespace(this.clusterName, select.metadata.name))
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
    }
  },
  created () {
    this.search()
  }
}
</script>

<style scoped>

</style>