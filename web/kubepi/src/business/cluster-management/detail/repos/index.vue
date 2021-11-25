<template>
  <layout-content>
    <complex-table :data="data">
      <template #header>
        <el-button-group>
          <el-button type="primary" size="small" @click="onCreate">
            {{ $t("commons.button.add") }}
          </el-button>
        </el-button-group>
        <br/>
      </template>
      <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.repo }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.created_time')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | datetimeFormat }}
        </template>
      </el-table-column>
      <fu-table-operations :buttons="buttons" :label="$t('commons.table.action')" fix/>
    </complex-table>
    <el-dialog :visible.sync="formDialogOpened" :close-on-click-modal="false"
               :title="$t('business.cluster.repo_auth')"  v-loading="isSubmitGoing"   z-index="10"
               width="70%"
               center>
      <el-form element-loading-spinner="el-icon-loading"
               element-loading-background="rgba(0, 0, 0, 0.8)" :model="repoForm" label-position="left"
               label-width="144px">
        <el-form-item :label="$t('business.cluster.repo')">
          <el-select v-model="repoForm.repo" style="width: 85%" filterable>
            <el-option v-for="(item, index) in repos" :key="index" :value="item.name">
              {{ item.name }}
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
                <el-button @click="formDialogOpened = false">{{ $t("commons.button.cancel") }}</el-button>
                <el-button type="primary" @click="onConfirm">{{ $t("commons.button.confirm") }}</el-button>
      </span>
    </el-dialog>
  </layout-content>
</template>

<script>
import LayoutContent from "@/components/layout/LayoutContent"
import {addClusterRepo, deleteClusterRepo, listClusterRepos} from "@/api/clusters"
import ComplexTable from "@/components/complex-table"
import {listRepoByCluster} from "@/api/imagerepos"

export default {
  name: "ClusterRepo",
  components: { LayoutContent, ComplexTable },
  props: {
    name: String
  },
  data () {
    return {
      data: [],
      loading: false,
      isSubmitGoing: false,
      formDialogOpened: false,
      repoForm: {},
      repos: [],
      buttons: [
        {
          label: this.$t("commons.button.delete"),
          icon: "el-icon-delete",
          click: (row) => {
            this.onDelete(row)
          }
        },
      ]
    }
  },
  methods: {
    list () {
      this.loading = true
      listClusterRepos(this.name).then(res => {
        this.data = res.data
      }).finally(() => {
        this.loading = false
      })
    },
    onCreate () {
      this.formDialogOpened = true
      this.listRepos()
    },
    onDelete(raw) {
      this.$confirm(this.$t("commons.confirm_message.delete"), this.$t("commons.message_box.alert"), {
        confirmButtonText: this.$t("commons.button.confirm"),
        cancelButtonText: this.$t("commons.button.cancel"),
        type: 'warning'
      }).then(() => {
        deleteClusterRepo(this.name, raw.repo).then(() => {
          this.$message.success(this.$t('commons.msg.delete_success'))
          this.list()
        })
      });
    },
    listRepos () {
      listRepoByCluster(this.name).then(res => {
        this.repos = res.data
      })
    },
    onConfirm() {
      this.isSubmitGoing = true
      this.repoForm.cluster = this.name
      addClusterRepo(this.name,this.repoForm).then(() =>{
        this.$message({
          type: "success",
          message: this.$t("commons.msg.create_success")
        })
        this.formDialogOpened = false
        this.list()
      }).finally(() =>{
        this.isSubmitGoing = false
      })
    }
  },
  created () {
    this.list()
  }
}
</script>

<style scoped>

</style>
