<template>
  <layout-content header="PersistentVolumes">
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
      <template #toolbar>
        <el-select v-model="conditions" @change="search(true)">
          <el-option label="All Namespaces" value=""></el-option>
          <el-option v-for="namespace in namespaces"
                     :key="namespace.metadata.name"
                     :label="namespace.metadata.name"
                     :value="namespace.metadata.name">
          </el-option>
        </el-select>
      </template>
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" min-width="200px" prop="metadata.name">
        <template v-slot:default="{row}">
          <el-link @click="openDetail(row)">{{ row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="metadata.namespace">
        <template v-slot:default="{row}">
          {{ row.status.phase }}
        </template>
      </el-table-column>
      <el-table-column label="ACCESS MODES" prop="metadata.labels" min-width="200px">
        <template v-slot:default="{row}">
          <div v-for="(name,index) in row.spec.accessModes " :key="index" style="display:inline-block">
            <el-tag>{{ name }}</el-tag>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="大小" prop="metadata.labels">
        <template v-slot:default="{row}">
          {{ row.spec.capacity.storage }}
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
import ComplexTable from "@/components/complex-table/index"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"
import {deletePvs, listPvs} from "../../../api/pv";



export default {
  name: "Pvs",
  components: { ComplexTable, LayoutContent, KoTableOperations },
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
      conditions: "",
      namespaces: [],
      buttons: [
        {
          label: this.$t("commons.button.edit"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "PvsEdit",
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
              name: "PvsEdit",
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
      listPvs(this.cluster, this.page.pageSize, this.page.nextToken, this.conditions).then(res => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.loading = false
        console.log("打印数据",this.data)
      })
    },
    onCreate () {
      this.$router.push({
        name: "PvsCreate",
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
          this.ps.push(deletePvs(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deletePvs(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "PvsDetail",
        params: { name: row.metadata.name, namespace: row.metadata.namespace },
        query: { yamlShow: false }
      })
    },
  },
  created () {
    this.cluster = this.$route.query.cluster
    this.search()
  }
}
</script>

<style scoped>

</style>
