<template>
  <layout-content header="ResourceQuotas">
    <complex-table :pagination-config="page" :data="data" :selects.sync="selects" v-loading="loading" @search="search">
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
      <el-table-column label="Request" min-width="100px">
        <template v-slot:default="{row}">
         <el-tag type="info" size="mini" v-if="row.status.used && row.status.hard"> pods: {{row.status.used.pods}} / {{row.status.hard.pods}}</el-tag>
         <el-tag type="info" size="mini" v-if="row.status.used && row.status.hard"> cpu: {{row.status.used['requests.cpu']}} / {{row.status.hard['requests.cpu']}}</el-tag>
         <el-tag type="info" size="mini" v-if="row.status.used && row.status.hard"> memory: {{row.status.used['requests.memory']}} / {{row.status.hard['requests.memory']}}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="Limit" min-width="100px">
        <template v-slot:default="{row}">
          <el-tag type="info" size="mini" v-if="row.status.used && row.status.hard"> cpu: {{row.status.used['limits.cpu']}} / {{row.status.hard['limits.cpu']}}</el-tag>
          <el-tag type="info" size="mini" v-if="row.status.used && row.status.hard"> memory: {{row.status.used['limits.memory']}} / {{row.status.hard['limits.memory']}}</el-tag>
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
import {deleteResourceQuota, listResourceQuotas} from "@/api/resourcequota"
import {downloadYaml} from "@/utils/actions"
import KoTableOperations from "@/components/ko-table-operations"

export default {
  name: "ResourceQuotas",
  components: { ComplexTable, LayoutContent,KoTableOperations },
  data () {
    return {
      data: [],
      page: {
        pageSize: 10,
        nextToken: ""
      },
      selects: [],
      loading: false,
      cluster: "",
      buttons: [
        {
          label: this.$t("commons.button.edit_yaml"),
          icon: "el-icon-edit",
          click: (row) => {
            this.$router.push({
              name: "ResourceQuotaEdit",
              params: {namespace:row.metadata.namespace,name:row.metadata.name}
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
          nextToken: ""
        }
      }
      listResourceQuotas(this.cluster, this.page.pageSize, this.page.nextToken).then(res => {
        this.data = res.items
        this.page.nextToken = res.metadata["continue"] ? res.metadata["continue"] : ""
        this.loading = false
      })
    },
    onCreate () {
      this.$router.push({
        name: "ResourceQuotaCreate",
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
          this.ps.push(deleteResourceQuota(this.cluster, row.metadata.namespace, row.metadata.name))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteResourceQuota(this.cluster, select.metadata.namespace, select.metadata.name))
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
        name: "ResourceQuotaDetail",
        params: {namespace:row.metadata.namespace,name:row.metadata.name}
      })
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
