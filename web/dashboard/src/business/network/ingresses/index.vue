<template>
  <layout-content header="Ingresses">
    <complex-table :pagination-config="page" :data="data" :selects.sync="selects" @search="search" v-loading="loading">
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
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" prop="metadata.name">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.namespace.namespace')" prop="metadata.namespace">
        <template v-slot:default="{row}">
          {{ row.metadata.namespace }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('business.network.target')" prop="metadata.rules">
        <template v-slot:default="{row}">
          <div v-for="(rule,index) in row.spec.rules" :key="index">
            <div v-for="(path,index) in rule.http.paths" :key="index">
              <el-link :href="'http://' + rule.host + (path.path ? path.path : '')" target="_blank">
                {{ "http://" + rule.host + (path.path ? path.path : "") }}
              </el-link>
              >
              <el-link @click="toResource('Service',row.metadata.namespace,path.backend.service.name)">
                {{
                  path.backend.service ? path.backend.service.name : ""
                }}
              </el-link>
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
import {listNamespace} from "@/api/namespaces"
import ComplexTable from "@/components/complex-table"
import KoTableOperations from "@/components/ko-table-operations"
import {deleteIngress, listIngresses} from "@/api/ingress"
import {mixin} from "@/utils/resourceRoutes"

export default {
  name: "Ingresses",
  components: { LayoutContent, ComplexTable, KoTableOperations },
  mixins: [mixin],
  data () {
    return {
      data: [],
      page: {
        pageSize: 10,
        nextToken: "",
      },
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
              name: "IngressEdit",
              params: { namespace: row.metadata.namespace, name: row.metadata.name },
              query: { yamlShow: false }
            })
          }
        },
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "IngressEdit",
              params: { name: row.metadata.name, namespace: row.metadata.namespace },
              query: { yamlShow: true }
            })
          }
        },
        {
          label: this.$t("commons.button.download_yaml"),
          icon: "el-icon-download",
          click: (row) => {
            downloadYaml(row.metadata.name + ".yml", row)
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
    }
  },
  methods: {
    search (init) {
      this.loading = true
      if (init) {
        this.page = {
          pageSize: this.page.pageSize,
          nextToken: "",
        }
      }
      listIngresses(this.cluster, this.page.pageSize, this.page.nextToken, this.conditions).then(res => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({
        name: "IngressCreate",
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
          this.ps.push(deleteIngress(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteIngress(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "IngressDetail",
        params: { name: row.metadata.name, namespace: row.metadata.namespace },
        query: { yamlShow: false }
      })
    },
    listAllNameSpaces () {
      listNamespace(this.cluster).then(res => {
        this.namespaces = res.items
      })
    }
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.listAllNameSpaces()
    this.search()
  }
}
</script>

<style scoped>

</style>
