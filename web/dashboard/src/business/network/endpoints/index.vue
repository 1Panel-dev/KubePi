<template>
  <layout-content header="Endpoints">
    <div style="float: left">
      <el-button type="primary" size="small" @click="yamlCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'endpoints',verb:'create'}">
        YAML
      </el-button>
      <el-button type="primary" size="small" @click="onCreate"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'endpoints',verb:'create'}">
        {{ $t("commons.button.create") }}
      </el-button>
      <el-button type="primary" size="small" :disabled="selects.length===0" @click="onDelete()"
                  v-has-permissions="{scope:'namespace',apiGroup:'',resource:'endpoints',verb:'delete'}">
        {{ $t("commons.button.delete") }}
      </el-button>
    </div>
    <complex-table :data="data" :selects.sync="selects" @search="search" v-loading="loading"
                   :pagination-config="paginationConfig" :search-config="searchConfig"
                   :showFullTextSwitch="true" @update:isFullTextSearch="OnIsFullTextSearchChange">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name" min-width="120px" show-overflow-tooltip>
        <template v-slot:default="{row}">
          <span class="span-link" @click="openDetail(row)">{{ row.metadata.name }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column label="Endpoints" prop="subsets" min-width="200px">
        <template v-slot:default="{row}">
          <div v-for="(subset,index) in row.subsets" v-bind:key="index">
            <div v-for="(address,index) in subset.addresses" v-bind:key="index" style="display:inline-block">
              <el-tag v-for="(port,index) in subset.ports" v-bind:key="index" type="info" size="mini">{{
                  address.ip
                }}:{{ port.port }}
              </el-tag>
            </div>
          </div>
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
import {downloadYaml} from "@/utils/actions"
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteEndPoint, getEndPoint, listEndPoints} from "@/api/endpoints"
import {checkPermissions} from "@/utils/permission"
import { searchFullTextItems } from "@/api/fulltextsearch/fulltextsearch"

export default {
  name: "Endpoints",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  data () {
    return {
      data: [],
      selects: [],
      cluster: "",
      loading: false,
      buttons: [
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "EndpointEdit",
              params: { name: row.metadata.name, namespace: row.metadata.namespace },
              query: { yamlShow: true }
            })
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "endpoints", verb: "update" })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml",getEndPoint(this.cluster,row.metadata.namespace,row.metadata.name))
          }
        },
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          },
          disabled: () => {
            return !checkPermissions({ scope: "namespace", apiGroup: "", resource: "endpoints", verb: "delete" })
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
      },
      isFullTextSearch: false
    }
  },
  methods: {
    search (resetPage) {
      this.loading = true
      if (resetPage) {
        this.paginationConfig.currentPage = 1
      }
      if(!this.isFullTextSearch){
        listEndPoints(this.cluster, true, this.searchConfig.keywords, this.paginationConfig.currentPage, this.paginationConfig.pageSize).then(res => {
        this.data = res.items
        this.loading = false
        this.paginationConfig.total = res.total
       })
      } else {
        listEndPoints(this.cluster, false)
        .then((res) => {
          const results = searchFullTextItems(res.items,this.searchConfig.keywords);
          this.data =results.slice(this.paginationConfig.currentPage*this.paginationConfig.pageSize-this.paginationConfig.pageSize,this.paginationConfig.currentPage*this.paginationConfig.pageSize)
          this.paginationConfig.total = results.length
        }).finally(() => {
          this.loading = false
        }) 
      }
    },
    onCreate () {
      this.$router.push({
        name: "EndpointCreate", query: { yamlShow: false }
      })
    },
    yamlCreate () {
      this.$router.push({
        name: "EndpointCreateYaml", query: { type: "endpoints" }
      })
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
          this.ps.push(deleteEndPoint(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteEndPoint(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "EndpointDetail",
        params: { name: row.metadata.name, namespace: row.metadata.namespace }
      })
    },
    //改变选项"是否全文搜索"
    OnIsFullTextSearchChange(val){
      this.isFullTextSearch=val
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
