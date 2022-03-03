<template>
  <layout-content>
    <div style="float: left; margin-bottom: 20px">
      <el-button type="primary" size="small" @click="onCreate">{{ $t("commons.button.add") }}</el-button>
      <el-button :disabled="selects.length===0" type="primary" size="small" @click="onDelete()">{{ $t("commons.button.delete") }}</el-button>
    </div>

    <complex-table :data="data" :selects.sync="selects">
      <el-table-column type="selection" fix></el-table-column>
      <el-table-column :label="$t('commons.table.name')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.repo }}
        </template>
      </el-table-column>
      <el-table-column :label="$t('commons.table.age')" min-width="100" fix>
        <template v-slot:default="{row}">
          {{ row.createAt | ageFormat }}
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
          <el-select v-model="repoForm.repos" style="width: 85%" filterable multiple>
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
      selects: [],
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
        this.ps = []
        if (raw) {
          this.ps.push(deleteClusterRepo(this.name, raw.repo))
        } else {
          if (this.selects.length > 0) {
            for (const select of this.selects) {
              this.ps.push(deleteClusterRepo(this.name, select.repo))
            }
          }
        }
        if (this.ps.length !== 0) { 
          Promise.all(this.ps)
            .then(() => {
              this.list()
              this.$message({
                type: "success",
                message: this.$t("commons.msg.delete_success"),
              })
            })
            .catch(() => {
              this.list()
            })
        }
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
